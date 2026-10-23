# tps

Personal project command hub — short aliases for shell commands, per linked repo. Run `tps` from any terminal, pick the project from your cwd, and stay in an interactive session while dev servers and tunnels run in other terminals.

## Setup

### macOS

```bash
curl -fsSL https://raw.githubusercontent.com/Cimmerial/tps/main/scripts/install.sh | bash
cd your-app-repo && tps setup && tps
```

### Linux

```bash
curl -fsSL https://raw.githubusercontent.com/Cimmerial/tps/main/scripts/install.sh | bash
cd your-app-repo && tps setup && tps
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/Cimmerial/tps/main/scripts/install.ps1 | iex
cd your-app-repo; tps setup; tps
```

Developers: `go install github.com/Cimmerial/tps/cmd/tps@latest`

Config path: OS user config dir + `/tps/config.yaml` (macOS `~/Library/Application Support/tps/config.yaml`, Linux `~/.config/tps/config.yaml`, Windows `%AppData%\tps\config.yaml`). Run `tps ih` in a session for install, link, and update steps on this machine.

---

## What tps is

tps is a **stay-open CLI** for day-to-day work on your own repos. Each linked project stores commands in config: an **alias** (full name), **abbreviation** (what you type), optional **description**, and the **shell** to run from the project root.

- **Interactive first** — bare `tps` opens a session with logo, project commands, and built-in meta commands until you `q`.
- **Cwd-aware** — resolves the active project from git root + stored path; heals when the folder moves but the remote matches.
- **Non-blocking servers** — long-running commands can detach into Terminal.app (macOS) or other terminals so the hub keeps accepting input.
- **Local and personal** — not CI, not multi-user hosting; YAML on your machine is the source of truth.

Each linked app repo also gets **`TPS.md`** and **`TPS_SPECS.md`** at its root (auto-written on `setup`, link, or update; both gitignored). `TPS.md` is a short cheat sheet; `TPS_SPECS.md` is the execution contract for commands and scripts. This README is the full user reference.

---

## Quick start

1. Install tps (see **Setup**).
2. In an app repo: `tps setup` (or `tps project add`) to link it.
3. Run `tps` — type a command **abbrev** to run it, or `h` for help.
4. Add commands with **`nc`** or edit **`config.yaml`** via **`c`**.

---

## Interactive session

When you run `tps` with no subcommand:

| UI | Meaning |
|----|---------|
| **project:** line | Current linked project (and optional update banner after `u`) |
| **PRJCT CMDS** | This project’s commands, grouped by tags / usage |
| **TPS CMDS** | Built-in shortcuts (usage-sorted hints) |
| **`⣾ >` prompt** | Shown while detached terminals for this project are running |

- Empty Enter does nothing; the session stays open until **`q`** / **`quit`**.
- **`tps`** redraws the menu at the top of the screen (keeps scrollback); **`clear`** also clears scrollback.
- If cwd does not match a linked project, tps prompts to pick one or use **`p`** to switch (`p1`, `p2`, … or `p name`). **`pb`** switches back.

---

## Linking projects

| Command | Purpose |
|---------|---------|
| `tps setup` | Install/update tps on PATH, link cwd (or `setup -n name --path dir`), refresh `TPS.md` and `TPS_SPECS.md` |
| `tps project add` | Link a repo (`add --path` optional) |
| `tps project list` | List linked projects |

Projects are stored in `config.yaml` with `id`, `name`, `path`, `git_remote`, and `commands`. If you move a clone, tps can update `path` when the git remote still matches.

---

## Project commands

Commands live under each project in config (or add them in-session).

### Fields (common)

| Field | Meaning |
|-------|---------|
| `alias` | Full name (used by `e` / `edit`) |
| `abbreviation` / `abbrev` | What you type in tps |
| `description` | Shown in the menu |
| `shell` | Run via `sh -c` from project root |
| `detach` | Open in a background terminal instead of blocking the hub (`nc -t`) |
| `close_on_complete` | With `detach`, close the tab when the command exits (macOS; `nc -close`) |
| `handoff` | Replace tps in the same terminal for TUIs/REPLs (`nc -H`) |
| `foreground` | Stay in this session; skip auto-detach heuristics (`nc -fg`) |
| `capture` | Watch log for a URL preset (`trycloudflare` / `cf`) or regex; print + clipboard (+ optional SMS on macOS; `nc -m`) |
| `capture_actions` | With `capture`: `print`, `clipboard`, `sms` (`nc -actions`) |
| `tags` | Menu folders — **first tag** is the group name |
| `precursor` | Run another command’s abbrev first; foreground steps share one shell |
| `shell1`, `shell2`, … | Multiple shells in parallel (often with `detach`; `nc`: repeat `-c` or `-c2`) |
| `handler` | In-process builtin (dev-repo commands; `nc -handler`) |
| `linux_shell`, `darwin_shell`, `windows_shell` | OS-specific shell strings (`nc -linux`, `-mac`, `-win`) |
| `menu` | Addon tree menu (see [ADDONS.md](ADDONS.md)) |

### Shell placeholders

Use `[INPUT1]`, `[INPUT2]`, … in `shell`. At run time, pass values after the abbrev: `ic -Cimmerial` → first `-value` fills `INPUT1`, and so on.

### In-session editing

| Keys | Action |
|------|--------|
| **`nc`** | New command: `nc -n name -a abbrev -c 'shell'`. Also `-mac` / `-linux` / `-win`, `-c2`, `-t`, `-close`, `-H`, `-fg`, `-m`, `-actions`, `-tags`, `-desc`, `-precursor`, `-handler`, `-menu`. `nc -help` lists every flag |
| **`e`** / **`edit`** | Change shells (`e hello` or pick from menu) |
| **`ce`** | Print YAML / `nc` examples |

Run **`h`** for the full built-in list; **`ce`** mirrors many config patterns.

---

## Built-in tps commands

### Session

| Command | Abbrev | Description |
|---------|--------|-------------|
| clear | | Clear terminal scrollback |
| tps | tps | Refresh logo and menu at top |
| quit | q | Exit tps |
| project | p · p1 | Switch linked project |
| pb | | Switch back to previous project |
| options | o | Logo/color, terminal colors, username, SMS capture, default IDE, preferred browser |
| open in desired ide | ide | Open project in Cursor/VS Code (or default IDE) |
| close all terms | vite | macOS: kill tab processes and close every Terminal.app window |

### Config

| Command | Abbrev | Description |
|---------|--------|-------------|
| config | c | Open `config.yaml` in IDE or `$EDITOR` |
| edit | e | Edit commands, tags, shells |
| new command | nc | Add a command from flags (`nc -help`; same fields as YAML) |
| terminals | t | List jobs · `tk` kill · `tk -a` all projects · `t -k N` one job · `tr` restart history |
| export | export | Export options and/or commands (`export -a` for everything) |
| import | import | Import options or merge commands (`import path/to/file.tps.yaml`) |
| config examples | ce | YAML and `nc` examples |
| setup | setup | Install tps + link repo (`setup -n name --path dir`) |

### Git

| Command | Abbrev | Description |
|---------|--------|-------------|
| git ship | g | `g your message` or `g ship …` (legacy `gs`) |
| status | | `git status` |
| git | | `git status`, or `git diff`, `git log -1`, … |

### Help

| Command | Abbrev | Description |
|---------|--------|-------------|
| help | h | Command list in the terminal |
| version | ver | tps version (addon version on addon projects) |
| install | ih | Install and link on this or another machine |
| doc | documentation | Open this README on GitHub |

### Update (when a tps dev clone is linked)

| Command | Abbrev | Description |
|---------|--------|-------------|
| update | u | `git pull` + rebuild, sync every `tps` on PATH, refresh `TPS.md` and `TPS_SPECS.md` |

Without a dev clone, **`u`** installs the latest [GitHub release](https://github.com/Cimmerial/tps/releases). Opening `tps` may also quietly rebuild when the linked dev repo changed.

---

## Detached terminals and capture

On **macOS**, detached jobs use **Terminal.app** in the background (`open -g`) so focus stays in your editor.

| Keys | Action |
|------|--------|
| **`t`** | List detached jobs for this project |
| **`tk`** / **`t -k`** | Kill this project’s jobs |
| **`tk -a`** / **`tk -all`** | Kill every tps-tracked job (all projects) |
| **`t -k N`** | Kill one job by number |
| **`tr`** / **`t -r`** | Kill all and replay launch history |
| **`vite`** | Kill processes and close all Terminal.app windows |

**Capture:** with `-m trycloudflare` on `nc`, or `capture:` in YAML, tps tails the job log, matches the URL, prints it, and copies to the clipboard (optional SMS notify on macOS via **options**).

Linux and Windows use supported terminal emulators when available; some behaviors (tab auto-close, `vite`, SMS) are **macOS-only** — see **Limitations**.

---

## Addons

Separate git repos (for example [tps-utilities](https://github.com/Cimmerial/tps-utilities)) link as normal projects and declare commands in `tps-addon.yaml`. Install from any project with **`iu`** / **`install utilities`**. Contract and menu trees: [ADDONS.md](ADDONS.md).

---

## CLI subcommands (non-interactive)

| Subcommand | Purpose |
|------------|---------|
| `tps` | Interactive session |
| `tps setup` | Bootstrap PATH + link project |
| `tps project …` | Add, list, remove projects |
| `tps install-help` / `ih` | Same install text as in-session `ih` |

---

## Limitations

- **Single user, local machine** — config and job state are not shared.
- **Project must be linked**; cwd should be under the git root (or heal via remote).
- **SMS capture notify** — macOS only.
- **`vite` / close all Terminal windows** — Terminal.app on macOS only.
- **Detached tab auto-close** — macOS only; other OSes may leave terminal windows open after `tk`.
- **`TPS.md` and `TPS_SPECS.md` in app repos** are overwritten on update — keep custom notes elsewhere.

---

## Developers

- Module: `github.com/Cimmerial/tps`
- Design notes and roadmap: [DESIGN_DOC.md](DESIGN_DOC.md)
- Instance cheat sheet source (embedded into releases): `internal/instance/README.md`
- Run tests: `go test ./...`

Contributions and issues: [github.com/Cimmerial/tps](https://github.com/Cimmerial/tps)
