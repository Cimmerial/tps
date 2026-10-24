#!/usr/bin/env bash
set -euo pipefail

INSTALL_DIR="${TPS_INSTALL_DIR:-$HOME/.local/bin}"

rm -f "$INSTALL_DIR/tps"
if [[ -d "$INSTALL_DIR" ]] && [[ -z "$(ls -A "$INSTALL_DIR" 2>/dev/null)" ]]; then
  rmdir "$INSTALL_DIR" 2>/dev/null || true
fi

for rc in "$HOME/.zshrc" "$HOME/.bashrc"; do
  [[ -f "$rc" ]] || continue
  if grep -qF "$INSTALL_DIR" "$rc" 2>/dev/null; then
    grep -vF "$INSTALL_DIR" "$rc" >"${rc}.tps-uninstall.$$"
    mv "${rc}.tps-uninstall.$$" "$rc"
    echo "removed PATH line from $rc"
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

hash -r 2>/dev/null || true
echo "tps uninstalled from $INSTALL_DIR (PATH updated in this shell)"
