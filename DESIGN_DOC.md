# TPS — design & current state

## Goal

TPS is a **personal project command hub**: one interactive CLI where linked repos expose short aliases (`r`, `cf`, …) instead of long shell one-liners. It should feel fast to open from any terminal, stay open while you work, and stay usable when projects move or you switch machines.

**Non-goals (for now):** replacing your shell, CI, or IDE; running as a daemon; multi-user hosting.

---

## What exists today (Go CLI)

**Stack:** Go, Cobra, YAML config at `~/Library/Application Support/tps/config.yaml` (macOS `UserConfigDir`).

**Entry:** bare `tps` → interactive session; subcommands for `project`, `setup`, `install-help` (`ih`), etc.

### Interactive session

- ASCII logo, then **`project: <name>`** (and optional **`- tps updated`** after `u`).
- **PRJCT CMDS:** top project commands by usage (abbrev + description).
- **TPS CMDS:** built-ins sorted by meta usage; **`⣾ >`** prompt when detached terminals are running.
- Empty Enter = no-op; session stays open until **`q`**.
- Resolves project from **cwd** (git root + path / remote healing). If cwd is an unlinked app repo but config already has projects, **auto-links** once.
- If cwd doesn’t match, **project picker** or **`p`** to switch.

### Projects

| Mechanism | Purpose |
|-----------|---------|
| `tps project add` / `add --path` | Link a repo (name defaults to folder). |
| `tps setup` / `tps setup -n name` | Install/update `tps`, PATH hint, link cwd (or path), refresh `TPS.md` and `TPS_SPECS.md`. |
| Git remote + path in YAML | **Heal** when folder moves (`ResolveFromCwd`). |

### Project commands

Defined in config per project: `alias`, `shell` / `shell1`… / OS-specific shells, `abbreviation`, optional `description`, `tags`, `precursor`, `detach`, `close_on_complete`, `handoff`, `foreground`, `capture` + `capture_actions`, `handler`, `menu`. In-session **`nc`** sets the same fields (`nc -help`).

| Interactive | Meaning |
|-------------|---------|
| **`nc -n name -a abbrev -c '…'`** | Add command. **`-mac` / `-linux` / `-win`** (and `-c2`, `-mac2`, …) set per-OS and extra shells. **`-t`** detach, **`-close`** close the tab when it exits, **`-H`** handoff, **`-fg`** stay in this session, **`-m`** capture, **`-actions`**, **`-tags`**, **`-desc`**, **`-precursor`**, **`-handler`**, **`-menu`**. **`nc -help`** lists them. |
| **`e` / `edit`** | Change shell for an alias. |
| **`r`**, etc. | Run alias; dev-like shells auto-**detach** on macOS. |

**Detached run (macOS Terminal.app):**

- Opens a **background** Terminal tab/window (`open -g`), does not steal focus from Cursor/VS Code.
- **`t`** list jobs · **`tk`** / **`t -k`** kill this project’s jobs · **`tk -a`** / **`tk -all`** kill every tps-tracked job · **`t -k N`** one job · **`tr`** / **`t -r`** kill all and replay **launch history** in order. (`tc` / `tck` / `tcr` still work.) **`vite`**: from any project on macOS, kill tab processes then close all Terminal.app windows (same strategy as `tk`).
- Kill: process tree + optional **`--port`** cleanup; close tab by TTY or window title (`<project-id>-<n>.command`).
- **`capture`:** tee log in run script, match preset (`trycloudflare` / `cf`) or regex → print URL + **clipboard**.

### Built-in TPS commands (representative)

| Keys | Role |
|------|------|
| `h` | Help |
| `c` | Open config in Cursor/VS Code/`$EDITOR` (GUI editors non-blocking) |
| `p` | Switch project |
| `u` | Pull dev repo (if linked), rebuild, sync all `tps` on PATH, refresh `TPS.md` and `TPS_SPECS.md`, restart |
| `g ship` / `git ship` | git add/commit/push (`gs` legacy) |
| `g …` / `git …` | Git passthrough |
| `export` / `export -a` | Menu: options, per-project commands, or everything; save dialog (`export -a` skips menu) |
| `import` / `import <file>` | Pick destination (options or project), merge commands by alias |
| `ih` | Install / link / update help |
| `tps` | Clear visible screen and redraw menu at top (keeps scrollback) |
| `clear` | Clear terminal scrollback + refresh (VS Code Cmd+K style) |
| `q` | Quit |

Startup **auto-update** (quiet `go install`) when a tps dev repo is in config, unless just restarted from `u`.

### Config & persistence

- **Projects:** commands (`shell`, `shell1`…, `tags` for menu folders, OS-specific shells), `command_usage`, `meta_usage`.
- **Jobs:** per-project YAML under `…/tps/jobs/` (active PIDs, **launch history** for `tr`).
- **PIDs/logs:** `…/tps/jobs/pids/`, `…/logs/` for capture.

---

## Architecture (packages)

```
cmd/tps          → main
internal/cli     → interactive loop, menu resolution, t/nc/setup/edit/…
internal/config  → YAML, editor, capture presets
internal/project → git root, remote, link, resolve
internal/runner  → foreground `sh -c` in project root; detach heuristics
internal/jobs    → spawn Terminal, track/kill/restart, capture watch
internal/git     → ship + passthrough
internal/update  → install, fingerprint, restart flags, sync all binaries
internal/ui      → logo, colors, spinner, clipboard, prompt
```

---

## Instance docs (`TPS.md`, `TPS_SPECS.md`)

Each **linked app repo** gets:

- **`TPS.md`** — cheat sheet (limits, common keys). Canonical source: **`internal/instance/README.md`**.
- **`TPS_SPECS.md`** — execution spec for commands, scripts, handoff/detach, YAML fields. Canonical source: **`internal/instance/SPECS.md`**.

Both are embedded in release binaries and written to `<project>/` on sync; both are gitignored.

**When you change behavior, menus, detach/capture, or documented limits:** edit `internal/instance/README.md`, then ship tps.

**When you change how commands/scripts should be written for tps:** edit `internal/instance/SPECS.md`, then ship tps.

Copies refresh when users run `tps setup`, link a project (`tps project add`), or update tps (`u` / auto-update after a rebuilt binary).

---

## Design principles

1. **Interactive first** — optimize for daily use inside `tps`, not one-shot flags.
2. **Minimal scope per change** — small diffs; match existing style.
3. **macOS-first for terminals** — detached jobs and AppleScript; other OSes fall back to foreground or errors.
4. **Config is source of truth** — projects heal by path/remote; usage sorts menus.
5. **Don’t block the hub** — long servers run in other terminals; hub keeps accepting input.

**Addons:** separate repos that link as projects; contract in [ADDONS.md](ADDONS.md). Help lists **TPS ADDONS** (`iu` / `install utilities`).

---

## Open goals / roadmap

| Area | Status |
|------|--------|
| Short aliases per project | **Done** |
| Stay-open REPL + usage-sorted menus | **Done** |
| Open/manage other terminals | **Done** (detach, t, kill, restart, capture) |
| Sequences + wait for events | **Partial** (`capture` regex); no generic step runner yet |
| Export/import projects & commands | **Done** (`export` / `import` / `setup`) |
| Run `tps` from anywhere on PATH | **Done** (`go install`, `ih` PATH line) |
| Stable identity when project moves | **Done** (git remote + path heal) |
| New machine bootstrap | **Done** (`setup`, `ih`, install scripts, optional config copy) |
| Xcode / GUI build from CLI | **Not started** |
| iTerm vs Terminal, Linux terminals | **Not started** (Terminal.app only for detach) |
| `TPS_DEBUG` / structured logging | **Not started** (kill debug removed after validation) |

---

## Resolved questions (from early notes)

- **Run from any terminal:** install `tps` to `GOPATH/bin`, PATH in `~/.zshrc`; bare `tps` from any cwd.
- **Project moved:** match `git_remote` and update stored path on resolve.
- **Computer to computer:** clone tps, `setup`, copy config or re-link projects; same config path on macOS.

---

## Language

**Go** — shipped implementation. Early C++/Python/Rust notes are obsolete; Go fits install-via-`go install`, single binary, and shelling out to git/Terminal/osascript.
