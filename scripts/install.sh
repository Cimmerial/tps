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

if install_release "$VERSION"; then
  echo "installed to $INSTALL_DIR/tps"
else
  echo "release install failed; trying go install..." >&2
  install_go_module
fi

ensure_path_line
ensure_path_now
echo "run: which tps && tps version"
