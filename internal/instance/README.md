# TPS in this repo

This file is **copied automatically** from the [tps](https://github.com/Cimmerial/tps) dev repo when you run `tps setup`, link a project, or update tps (`u`). Do not edit it here — change `internal/instance/README.md` in the tps repo instead, then update tps so copies refresh.

## What TPS is

TPS is a **personal project command hub**: one interactive CLI where this repo’s short aliases (`r`, `cf`, …) run your real shell commands from the project root. Open `tps` from any terminal; it picks the project from your current directory (git root + remote healing if the folder moved).

Config lives in your OS user config dir (e.g. macOS `~/Library/Application Support/tps/config.yaml`, Linux `~/.config/tps`, Windows `%AppData%\tps`). Link this repo once with `tps setup` or `tps project add`.

## Day to day

| In `tps` | Meaning |
|----------|---------|
| Type an **abbrev** or alias | Run that project command |
| **`nc`** | New command (`nc -n name -a abbrev -c 'shell'`) |
| **`e` / `edit`** | Change a command’s shell |
| **`t`** | List detached terminal jobs |
| **`t -k` / `tk`** | Kill detached jobs for this project |
| **`tk -a` / `tk -all`** | Kill every tps-tracked detached job (all projects) |
| **`vite`** | (macOS) Kill processes and close every Terminal.app window |
| **`tr` / `t -r`** | Kill all and replay launch history |
| **`p`** | Switch linked project |
| **`export` / `import`** | Backup or merge project YAML |
| **`ih`** | Install / link / update help for this machine |
| **`iu`** | Install **tps-utilities** addon (`h` → TPS ADDONS) |
| **`u`** | Update tps (git pull + build if dev repo linked, else GitHub release) |
| **`q`** | Quit (session stays in terminal) |

Detached servers can use **`-t`** on `nc` or auto-detach heuristics — they open in a **background terminal** (Terminal.app on macOS; Windows Terminal or a new console on Windows; common Linux emulators when installed). Tab auto-close on completion is **macOS-only** for now.

**URL capture:** `nc` with **`-m trycloudflare`** (or a regex) watches the log, prints the URL, and copies it to the clipboard.

## Install / update tps

**macOS / Linux:** `curl -fsSL https://raw.githubusercontent.com/Cimmerial/tps/main/scripts/install.sh | bash`

**Windows (PowerShell):** `irm https://raw.githubusercontent.com/Cimmerial/tps/main/scripts/install.ps1 | iex`

Then from this app repo: `tps setup` and `tps`.

Developers: `go install github.com/Cimmerial/tps/cmd/tps@latest`

Opening `tps` quietly rebuilds when a linked dev clone changed; **`u`** also syncs every `tps` on PATH and refreshes this file.

## Current limitations

- **SMS capture notify** — macOS only (Messages / Shortcuts).
- **vite / close all Terminal windows** — macOS Terminal.app only.
- **Detached tab auto-close** — macOS only; Linux/Windows kill jobs via `tk` but may leave the terminal window open.
- **Single user, local machine** — not a shared service or CI runner.
- **Project must be linked** in config; cwd must match git root (or heal via `git remote`).
- **This file** is overwritten on tps update; custom notes belong in your own docs, not `TPS.md`.

For full design notes and roadmap, see `DESIGN_DOC.md` in the tps dev repo.
