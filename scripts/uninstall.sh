#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="${TPS_INSTALL_DIR:-$HOME/.local/bin}"

if [[ -f "$INSTALL_DIR/tps" ]]; then
  rm -f "$INSTALL_DIR/tps"
  echo "removed $INSTALL_DIR/tps"
fi
if [[ -d "$INSTALL_DIR" ]] && [[ -z "$(ls -A "$INSTALL_DIR" 2>/dev/null)" ]]; then
  rmdir "$INSTALL_DIR" 2>/dev/null || true
fi

for rc in "$HOME/.zshrc" "$HOME/.bashrc"; do
  if [[ -f "$rc" ]]; then
    if grep -qF "$INSTALL_DIR" "$rc" 2>/dev/null; then
      grep -vF "$INSTALL_DIR" "$rc" >"$rc.tmp" && mv "$rc.tmp" "$rc"
      echo "removed PATH line from $rc"
    fi
  fi
done

new_path=""
IFS=':'
for p in $PATH; do
  [[ -z "$p" ]] && continue
  [[ "$p" == "$INSTALL_DIR" ]] && continue
  if [[ -z "$new_path" ]]; then
    new_path="$p"
  else
    new_path="$new_path:$p"
  fi
done
unset IFS
export PATH="$new_path"

if command -v hash >/dev/null 2>&1; then
  hash -r 2>/dev/null || true
fi

echo "tps uninstalled"
