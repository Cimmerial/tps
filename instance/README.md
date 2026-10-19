# TPS in this repo

This file is **copied automatically** from the [tps](https://github.com/Cimmerial/tps) dev repo when you run `tps setup`, link a project, or update tps (`u`). Do not edit it here — change `instance/README.md` in the tps repo instead, then update tps so copies refresh.

## What TPS is

TPS is a **personal project command hub**: one interactive CLI where this repo’s short aliases (`r`, `cf`, …) run your real shell commands from the project root. Open `tps` from any terminal; it picks the project from your current directory (git root + remote healing if the folder moved).

Config lives in macOS `~/Library/Application Support/tps/config.yaml` (commands, usage order, detached jobs). This repo only needs to be **linked** once (`tps setup` or `tps project add`).

## Day to day

| In `tps` | Meaning |
|----------|---------|
| Type an **abbrev** or alias | Run that project command |
| **`nc`** | New command (`nc -n name -a abbrev -c 'shell'`) |
| **`e` / `edit`** | Change a command’s shell |
| **`t`** | List detached terminal jobs (macOS) |
| **`t -k` / `tk`** | Kill detached jobs for this project |
| **`tk -a` / `tk -all`** | Kill every tps-tracked detached job (all projects) |
| **`vite`** | (macOS, any project) Kill processes and close every Terminal.app window |
| **`tr` / `t -r`** | Kill all and replay launch history |
| **`p`** | Switch linked project |
| **`export` / `import`** | Backup or merge project YAML |
| **`ih`** | Install / update help; **`ih -copy`** paste-safe setup script |
| **`u`** | Rebuild tps from your linked tps dev repo and restart |
| **`q`** | Quit (session stays in terminal) |

Detached servers (e.g. `rails s`, Vite) can use **`-t`** on `nc` or auto-detach heuristics on macOS — they open in **Terminal.app** in the background so the hub stays in Cursor/VS Code. In **edit**, toggle **kill detachment after completion** (or set `close_on_complete: true` in YAML) to close the tab when a one-shot detached command exits; long-running servers usually leave this off.

**URL capture:** `nc` with **`-m trycloudflare`** (or a regex) watches the log, prints the URL, and copies it to the clipboard.

## Install / update tps

1. Clone the tps repo and link it: `tps project add` from inside that repo.
2. From any linked app repo: `tps setup` or run `ih -copy` from the app directory and paste the script.
3. Update binary: inside `tps`, run **`u`** (or `go install ./cmd/tps` in the tps repo).

## Current limitations

- **macOS-first for detached terminals** — background tabs use Terminal.app + AppleScript; Linux/other OS run foreground or error.
- **Single user, local machine** — not a shared service or CI runner.
- **No generic multi-step pipelines** — besides log `capture` regex/presets, there is no built-in “wait for step B after A” runner yet.
- **iTerm / other terminals** — not supported for detach; Terminal.app only.
- **Project must be linked** in config; cwd must match git root (or heal via `git remote`).
- **This file** is overwritten on tps update; custom notes belong in your own docs, not `TPS.md`.

For full design notes and roadmap, see `DESIGN_DOC.md` in the tps dev repo.
