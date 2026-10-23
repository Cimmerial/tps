# TPS execution spec (for commands and scripts)

This file is **copied automatically** from the [tps](https://github.com/Cimmerial/tps) dev repo when you run `tps setup`, link a project, or update tps (`u`). Do not edit it here — change `internal/instance/SPECS.md` in the tps repo instead, then update tps so copies refresh.

For day-to-day keys and limits, see **`TPS.md`** in this repo. For full user docs, see the tps repo `README.md`.

## How commands run

- **Working directory:** every project command runs with **cwd = linked project root** (git root), not the directory where you launched `tps`.
- **Shell invocation:**
  - macOS / Linux: `sh -c '<shell>'`
  - Windows: `powershell -NoProfile -Command '<shell>'`
- **Config location:** OS user config dir (e.g. macOS `~/Library/Application Support/tps/config.yaml`). Per-repo commands live under that project entry in YAML, or add them in-session with **`nc`**.

## Writing shell strings

- Use **paths relative to the project root** (`./scripts/build.sh`, `npm run dev`), not paths that assume cwd is a subfolder.
- Prefer small **launcher scripts** under `scripts/` and thin `shell:` lines that call them; share setup via `source ./scripts/env.sh` in the shell or a **precursor** command.
- **Foreground** commands block the tps hub until they exit. **Detached** commands open in a separate terminal so tps stays interactive.

## YAML / `nc` fields (common)

| Field | Meaning |
|-------|---------|
| `alias` | Full name (`e` / `edit` use this) |
| `abbreviation` / `abbrev` | What you type in tps |
| `description` | Menu text |
| `shell` | Default shell string |
| `detach` | Run in a background terminal (`nc -t`) |
| `close_on_complete` | With `detach`, close the tab when the command exits (macOS) |
| `handoff` | Replace tps in the **same** terminal (`nc -H`) — use for full-screen TUIs and REPLs |
| `precursor` | Run another abbrev first; foreground steps share one shell (`export` / `source` persist) |
| `capture` | Preset (`trycloudflare` / `cf`) or regex; watch log, print URL, clipboard (`nc -m`) |
| `capture_actions` | With `capture`: `print`, `clipboard`, `sms` (`nc -actions`) |
| `foreground` | Stay in this session; skip auto-detach (`nc -fg`) |
| `tags` | Menu folders (first tag = group name) |
| `shell1`, `shell2`, … | Multiple shells (often parallel, with `detach`). `nc`: repeat `-c`, or `-c2` |
| `linux_shell`, `darwin_shell`, `windows_shell` | OS-specific shells (`nc -linux`, `-mac` / `-darwin`, `-win` / `-windows`) |
| `handler` | In-process builtin id (no shell; `nc -handler`) |
| `menu` | Addon tree menu (see `ADDONS.md` in the tps repo; `nc -menu`) |

**Inputs:** use `[INPUT1]`, `[INPUT2]`, … in `shell`. At run time: `abbrev -first -second` fills inputs in order.

**In-session examples:** run **`ce`** inside tps for YAML and `nc` patterns. **`nc -help`** lists every flag.

## `nc` flags

| Flag | Writes |
|------|--------|
| `-n` / `--name` | `alias` (required) |
| `-a` / `--abbrev` | `abbreviation` |
| `-desc` / `--description` | `description` (defaults to the name) |
| `-c` / `-shell` | `shell`. Repeat for `shell2`, `shell3`, … or pass `-c2` |
| `-mac` / `-darwin` | `darwin_shell` (repeat or `-mac2`) |
| `-linux` | `linux_shell` (repeat or `-linux2`) |
| `-win` / `-windows` | `windows_shell` (repeat or `-win2`) |
| `-t` / `--terminal` | `detach` |
| `-close` | `close_on_complete` (also sets detach) |
| `-H` / `--handoff` | `handoff` |
| `-fg` / `--foreground` | `foreground` (stay in this session; skips auto-detach) |
| `-m` / `-capture` | `capture` (also sets detach) |
| `-actions` | `capture_actions`: `print`, `clipboard`, `sms` |
| `-tags` | `tags` (comma-separated; first tag is the menu folder) |
| `-precursor` | `precursor` |
| `-handler` | `handler` (in-process builtin; no shell required) |
| `-menu` | `menu` (addon tree; no shell required) |

A command needs `-c`, `-mac`, `-linux`, `-win`, `-handler`, or `-menu`. On each OS, tps runs that OS’s shells when any are set, and the `-c` shells otherwise. Quote a shell that contains spaces or another `nc` flag (`-t`, `-n`, `-mac`, …).

```
nc -n build -a b -mac ./scripts/build-mac.sh -linux ./scripts/build.sh -win "powershell -File scripts/build.ps1"
nc -n stack -a st -c "npm run dev" -c "npx cloudflared tunnel --url http://localhost:5173" -t
```

## Handoff vs detach

| Pattern | When to use |
|---------|-------------|
| **`handoff: true`** / **`nc -H`** | Interactive app owns the terminal (TUI, `vim`, long REPL). tps does not return until that process exits. |
| **`detach: true`** / **`nc -t`** | Dev server, tunnel, watcher — keeps running while you stay in tps. List/kill with **`t`** / **`tk`**. |

## Repo-specific commands

This file does **not** list this project’s abbrevs. They live in your machine’s `config.yaml`, or in a committed export (`tps export` → `ProjectName.tps.yaml`, merge with `import`). Put **project-only** launcher notes in your own docs or Cursor rules, not in `TPS_SPECS.md`.

## Addons

Separate repos (e.g. tps-utilities) link as projects and declare commands in **`tps-addon.yaml`**. Contract: **`ADDONS.md`** in the tps repo. Install from any linked project with **`iu`** / **`install utilities`**.

## Platform limits (summary)

- Detached tab auto-close, **`vite`** (close all Terminal.app windows), and some capture SMS notify: **macOS-first**; other OSes may leave terminal windows open after `tk`.
- Single user, local machine — not CI or shared hosting.
