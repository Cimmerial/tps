#!/usr/bin/env bash
# Tag a release and push to GitHub. CI (.github/workflows/release.yml) builds assets on v* tags.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

VERSION_FILE="internal/version/version.go"
REMOTE="${TPS_RELEASE_REMOTE:-origin}"
BRANCH="${TPS_RELEASE_BRANCH:-main}"

usage() {
	cat <<'EOF'
Usage: push-release.sh [vX.Y.Z] [-y|--yes] [-n|--dry-run]

  With no version: bump patch from the latest v* tag (v0.1.0 -> v0.1.1).
  Updates internal/version/version.go DevVersion, commits, tags, pushes main + tag.

Environment:
  TPS_RELEASE_REMOTE   git remote (default: origin)
  TPS_RELEASE_BRANCH   branch to push (default: main)
EOF
}

die() {
	echo "push-release: $*" >&2
	exit 1
}

simplified_output() {
	[[ "${TPS_SIMPLIFIED_OUTPUT:-}" == "1" ]]
}

git_quiet() {
	if simplified_output; then
		git -q "$@"
	else
		git "$@"
	fi
}

yes=false
dry_run=false
explicit=""

while [[ $# -gt 0 ]]; do
	case "$1" in
	-y | --yes)
		yes=true
		shift
		;;
	-n | --dry-run)
		dry_run=true
		shift
		;;
	-h | --help)
		usage
		exit 0
		;;
	v*.*)
		explicit="$1"
		shift
		;;
	*)
		die "unknown argument: $1 (try --help)"
		;;
	esac
done

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
	die "not a git repository"
fi

branch="$(git branch --show-current)"
if [[ "$branch" != "$BRANCH" ]]; then
	die "on branch $branch; switch to $BRANCH first"
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
	die "commit or stash all changes before releasing (git status is not clean)"
fi

git fetch "$REMOTE" --tags

latest="$(git tag -l 'v*' --sort=-v:refname | head -n1)"
if [[ -z "$latest" ]]; then
	latest="v0.0.0"
fi

next_tag() {
	local t="${1#v}"
	local major minor patch
	IFS=. read -r major minor patch <<<"$t"
	patch="${patch:-0}"
	echo "v${major}.${minor}.$((patch + 1))"
}

if [[ -n "$explicit" ]]; then
	tag="$explicit"
else
	tag="$(next_tag "$latest")"
fi

if [[ ! "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
	die "version must look like vX.Y.Z (got $tag)"
fi

semver="${tag#v}"
if [[ ! -f "$VERSION_FILE" ]]; then
	die "missing $VERSION_FILE"
fi

if git rev-parse "$tag" >/dev/null 2>&1; then
	die "tag $tag already exists"
fi

echo "Latest tag:  ${latest:-<none>}"
echo "Release tag: $tag"
echo "Remote:      $REMOTE ($BRANCH)"
echo

if [[ "$yes" != true && "$dry_run" != true ]]; then
	read -r -p "Commit DevVersion, tag, and push? [y/N] " ans
	case "$ans" in
	y | Y | yes | YES) ;;
	*) die "aborted" ;;
	esac
fi

bump_dev_version() {
	local ver="$1"
	if [[ "$(uname -s)" == Darwin ]]; then
		sed -i '' "s/const DevVersion = \"[^\"]*\"/const DevVersion = \"$ver\"/" "$VERSION_FILE"
	else
		sed -i "s/const DevVersion = \"[^\"]*\"/const DevVersion = \"$ver\"/" "$VERSION_FILE"
	fi
}

if [[ "$dry_run" == true ]]; then
	echo "[dry-run] would set DevVersion=$semver in $VERSION_FILE"
	echo "[dry-run] would run: go test ./..."
	echo "[dry-run] would commit, tag $tag, push $REMOTE $BRANCH and $tag"
	exit 0
fi

bump_dev_version "$semver"

if simplified_output; then
	echo "running tests..."
	if test_out="$(go test ./... 2>&1)"; then
		ok_count=0
		skip_count=0
		while IFS= read -r line; do
			case "$line" in
			ok\ *) ok_count=$((ok_count + 1)) ;;
			?\ *) skip_count=$((skip_count + 1)) ;;
			esac
		done <<<"$test_out"
		if [[ "$ok_count" -gt 0 && "$skip_count" -gt 0 ]]; then
			echo "tests ok ($ok_count packages, $skip_count without tests)"
		elif [[ "$skip_count" -gt 0 ]]; then
			echo "tests ok ($skip_count packages without tests)"
		elif [[ "$ok_count" -gt 0 ]]; then
			echo "tests ok ($ok_count packages)"
		else
			echo "tests ok"
		fi
	else
		echo "$test_out"
		exit 1
	fi
else
	echo "running tests..."
	go test ./...
fi

if ! git diff --quiet -- "$VERSION_FILE"; then
	git add "$VERSION_FILE"
	if simplified_output; then
		git_quiet commit -m "release $tag"
	else
		git commit -m "release $tag"
	fi
fi

git tag -a "$tag" -m "release $tag"

if simplified_output; then
	git_quiet push "$REMOTE" "$BRANCH"
	git_quiet push "$REMOTE" "$tag"
else
	git push "$REMOTE" "$BRANCH"
	git push "$REMOTE" "$tag"
fi

echo
echo "Pushed $tag. Watch GitHub Actions → release, then users can re-run install.ps1 / install.sh."
