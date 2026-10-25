#!/usr/bin/env bash
set -euo pipefail

REPO="Cimmerial/tps"
INSTALL_DIR="${TPS_INSTALL_DIR:-$HOME/.local/bin}"
VERSION="${TPS_VERSION:-latest}"

os="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$os" in
  darwin) goos="darwin" ;;
  linux) goos="linux" ;;
  *) echo "unsupported OS: $os" >&2; exit 1 ;;
esac

arch="$(uname -m)"
case "$arch" in
  x86_64|amd64) goarch="amd64" ;;
  arm64|aarch64) goarch="arm64" ;;
  *) echo "unsupported arch: $arch" >&2; exit 1 ;;
esac

install_go_module() {
  if ! command -v go >/dev/null 2>&1; then
    echo "go not found; install Go or set TPS_VERSION to an existing release" >&2
    return 1
  fi
  if [[ -f go.mod ]] && [[ -d cmd/tps ]]; then
    echo "installing from local clone..."
    go install ./cmd/tps
  else
    echo "installing from module..."
    go install "github.com/${REPO}/cmd/tps@latest"
  fi
  bindir="$(go env GOPATH)/bin"
  mkdir -p "$INSTALL_DIR"
  if [[ -f "$bindir/tps" ]]; then
    cp -f "$bindir/tps" "$INSTALL_DIR/tps"
    chmod +x "$INSTALL_DIR/tps"
  fi
}

install_release() {
  ver="$1"
  if [[ "$ver" == "latest" ]]; then
    url="https://api.github.com/repos/${REPO}/releases/latest"
    ver="$(curl -fsSL "$url" | sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p' | head -1)"
  fi
  ver="${ver#v}"
  asset="tps_${ver}_${goos}_${goarch}.tar.gz"
  base="https://github.com/${REPO}/releases/download/v${ver}"
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' EXIT
  echo "downloading ${asset}..."
  if ! curl -fsSL "${base}/${asset}" -o "$tmp/archive.tar.gz"; then
    return 1
  fi
  if curl -fsSL "${base}/SHA256SUMS" -o "$tmp/SHA256SUMS" 2>/dev/null; then
    (cd "$tmp" && grep " ${asset}\$" SHA256SUMS | shasum -a 256 -c -)
  fi
  tar -xzf "$tmp/archive.tar.gz" -C "$tmp"
  mkdir -p "$INSTALL_DIR"
  install -m 755 "$tmp/tps" "$INSTALL_DIR/tps"
}

ensure_path_line() {
  line='export PATH="$PATH:'"$INSTALL_DIR"'"'
  for rc in "$HOME/.zshrc" "$HOME/.bashrc"; do
    if [[ -f "$rc" ]] || [[ "$rc" == *zshrc ]]; then
      if [[ -f "$rc" ]] && grep -qF "$INSTALL_DIR" "$rc" 2>/dev/null; then
        continue
      fi
      echo "$line" >>"$rc"
      echo "appended PATH to $rc"
    fi
  done
}

ensure_path_now() {
  case ":$PATH:" in
    *":$INSTALL_DIR:"*) ;;
    *) export PATH="$PATH:$INSTALL_DIR" ;;
  esac
}

git_installed() {
  command -v git >/dev/null 2>&1
}

offer_git_install() {
  if git_installed; then
    return 0
  fi
  if [[ "${TPS_SKIP_GIT_PROMPT:-}" == 1 ]]; then
    return 0
  fi
  if [[ ! -t 0 ]] && [[ ! -c /dev/tty ]]; then
    echo "git not found; install git for tps updates (https://git-scm.com)" >&2
    return 0
  fi
  local ans=""
  printf "Git is not installed (needed for tps update and git ship). Install Git now? [Y/n] "
  read -r ans </dev/tty 2>/dev/null || read -r ans
  case "$ans" in
    [Nn]*) return 0 ;;
  esac
  if [[ "$os" == darwin ]]; then
    if command -v brew >/dev/null 2>&1; then
      echo "installing git via Homebrew..."
      brew install git
    else
      echo "install Xcode Command Line Tools (includes git)..."
      xcode-select --install || true
    fi
  elif command -v apt-get >/dev/null 2>&1; then
    echo "installing git via apt..."
    sudo apt-get update -qq && sudo apt-get install -y git
  elif command -v dnf >/dev/null 2>&1; then
    echo "installing git via dnf..."
    sudo dnf install -y git
  elif command -v pacman >/dev/null 2>&1; then
    echo "installing git via pacman..."
    sudo pacman -S --noconfirm git
  else
    echo "install git manually: https://git-scm.com" >&2
    return 0
  fi
  if git_installed; then
    echo "git installed"
  else
    echo "git install finished; open a new shell if git is not found yet" >&2
  fi
}

go_installed() {
  command -v go >/dev/null 2>&1
}

ensure_go_bin_path() {
  if ! go_installed; then
    return 0
  fi
  local bindir
  bindir="$(go env GOPATH)/bin"
  case ":$PATH:" in
    *":$bindir:"*) ;;
    *) export PATH="$PATH:$bindir" ;;
  esac
  line='export PATH="$PATH:'"$bindir"'"'
  for rc in "$HOME/.zshrc" "$HOME/.bashrc"; do
    if [[ -f "$rc" ]] && ! grep -qF "$bindir" "$rc" 2>/dev/null; then
      echo "$line" >>"$rc"
      echo "appended go bin to $rc"
    fi
  done
}

offer_go_install() {
  if go_installed; then
    ensure_go_bin_path
    return 0
  fi
  if [[ "${TPS_SKIP_GO_PROMPT:-}" == 1 ]] || [[ "${TPS_SKIP_TOOLCHAIN_PROMPT:-}" == 1 ]]; then
    return 0
  fi
  if [[ ! -t 0 ]] && [[ ! -c /dev/tty ]]; then
    echo "go not found; install Go for tps dev updates (https://go.dev/dl/)" >&2
    return 0
  fi
  local ans=""
  printf "Go is not installed (needed for tps u / go install from the dev repo). Install Go now? [Y/n] "
  read -r ans </dev/tty 2>/dev/null || read -r ans
  case "$ans" in
    [Nn]*) return 0 ;;
  esac
  if [[ "$os" == darwin ]]; then
    if command -v brew >/dev/null 2>&1; then
      echo "installing go via Homebrew..."
      brew install go
    else
      echo "install Go from https://go.dev/dl/ or: brew install go" >&2
      return 0
    fi
  elif command -v apt-get >/dev/null 2>&1; then
    echo "installing go via apt..."
    sudo apt-get update -qq && sudo apt-get install -y golang-go
  elif command -v dnf >/dev/null 2>&1; then
    echo "installing go via dnf..."
    sudo dnf install -y golang
  elif command -v pacman >/dev/null 2>&1; then
    echo "installing go via pacman..."
    sudo pacman -S --noconfirm go
  else
    echo "install Go manually: https://go.dev/dl/" >&2
    return 0
  fi
  if go_installed; then
    ensure_go_bin_path
    echo "go installed"
  else
    echo "go install finished; open a new shell if go is not found yet" >&2
  fi
}

if install_release "$VERSION"; then
  echo "installed to $INSTALL_DIR/tps"
else
  echo "release install failed; trying go install..." >&2
  install_go_module
fi

ensure_path_line
ensure_path_now
offer_git_install
offer_go_install
echo "run: which tps && tps version"
