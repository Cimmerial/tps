# TPS addons

This document defines how **addon repositories** integrate with [TPS](https://github.com/Cimmerial/tps). Addon behavior in the CLI is implemented against this contract; addon repos should match it.

## Contract iteration

**Current iteration: 1** (must match `addon.CurrentContractIteration` in the tps repo).

Each addon declares the **highest contract iteration it implements** in `tps-addon.yaml`:

```yaml
contract_iteration: 1
```

| Situation | TPS behavior |
|-----------|----------------|
| `contract_iteration` equals current | OK |
| Missing or `0` | Warn: manifest does not declare a contract iteration |
| Less than current | Warn: addon may be missing newer contract features; update the addon repo |
| Greater than current | Warn: addon expects a newer TPS; update tps |

When this document changes in a way addons must follow (new manifest fields, menu schema, install rules, etc.), **bump the iteration**, add a changelog entry below, and ship tps. Addon maintainers bump `contract_iteration` after their repo conforms to that iteration.

### Contract changelog

| Iteration | Summary |
|-----------|---------|
| **1** | Initial contract: `tps-addon.yaml` manifest, `setup` without tps, `menu` trees with `dim_when` / `action` / picks, macOS Finder picks via tps, `iu` install flow, utilities UIC PNG-hub image conversion. |

## What an addon is

An addon is a **separate git repository** that is **linked as a normal TPS project** (`tps project add` / install flow). Commands defined by the addon appear on **that project only**. They do not run while another project is selected in the interactive session.

Install commands (for example `install utilities` / `iu`) may be run from **any** linked project. They clone, link, and register the addon project in config.

## Repository layout

| Path | Required | Purpose |
|------|----------|---------|
| `tps-addon.yaml` | yes | Manifest: id, name, clone URL, minimum TPS version, command definitions |
| `setup` | yes | Entry when someone clones the addon **without** TPS on `PATH` |
| `menus/*.yaml` | if using tree menus | Menu trees referenced by command `menu` fields |
| `bin/` | optional | Executables the menu `action.shell` invokes |

### `tps-addon.yaml`

```yaml
id: tps-utilities          # stable id (used by install and guards)
name: tps-utilities        # default link name
clone_url: https://github.com/Cimmerial/tps-utilities.git
contract_iteration: 1      # highest ADDONS.md iteration this repo implements
min_tps: ""                # optional semver hint
commands:
  - alias: universal image conversion
    abbreviation: uic
    description: convert images via PNG hub
    menu: menus/uic.yaml   # tree menu — no shell required
```

Command fields match TPS project commands in config (`alias`, `abbreviation` / `abbrev`, `description`, `shell`, `tags`, `handler`, `precursor`, …) plus:

| Field | Meaning |
|-------|---------|
| `menu` | Path **relative to the addon repo root** to a menu YAML file. When set, TPS opens the tree instead of running `shell`. |

**Tags:** omit `tags` to show the command under the **general** project folder (same as untagged commands in the project menu).

**Handlers:** prefer `menu` + `bin/` for addons; reserve `handler` for built-ins inside the TPS binary.

## Menu trees (`menus/*.yaml`)

```yaml
title: Universal Image Conversion
items:
  - id: convert
    label: convert
    children:
      - id: to-jpeg
        label: "→ JPEG"
        dim_when:
          magick_write: jpeg
        action:
          picks:
            - kind: files
              prompt: Select images
            - kind: folder
              prompt: Output folder (cancel for default)
              optional: true
          shell: '{root}/bin/tpsu img convert --to jpeg --out [PICK2] [PICK1]'
  - id: back
    label: back
    back: true
```

### Items

| Field | Meaning |
|-------|---------|
| `id` | Stable id (logging / future use) |
| `label` | Shown in the numbered menu |
| `children` | Sub-menu; selecting the item opens another tree level |
| `back` | Return to parent menu (or exit top-level menu) |
| `action` | Leaf: run picks then `shell` |
| `dim_when` | Item stays visible but **dimmed**; choosing it prints the reason and does not run |

### `dim_when`

```yaml
dim_when:
  reason: RAW is not a write target
```

Or probes (evaluated when the menu is drawn):

```yaml
dim_when:
  require_binary: magick
dim_when:
  magick_read: heic
dim_when:
  magick_write: svg
```

Dimmed rows use dim styling. Availability is re-checked each time the menu renders.

### `action`

```yaml
action:
  picks:
    - kind: file
      prompt: Select a file
    - kind: files
      prompt: Select images
    - kind: folder
      prompt: Choose output folder
      optional: true
  shell: '{root}/bin/tpsu img convert --to png --out [PICK2] [PICK1]'
```

| Pick `kind` | Behavior |
|-------------|----------|
| `file` | One file (macOS: Finder; else: type path) |
| `files` | One or more files |
| `folder` | One directory |

`optional: true` on a folder pick: user cancel → empty `[PICKn]` (addon shell must treat empty as default).

**Shell placeholders:**

| Token | Value |
|-------|--------|
| `{root}` | Addon project root (absolute) |
| `[PICK1]` … `[PICKn]` | Picks in order; file lists are space-separated, shell-quoted by TPS |

Addons **must not** call Finder themselves; only TPS performs GUI picks.

### Navigation

- Numbered choices; invalid input shows a short error and redraws.
- `0`, `q`, or `back` item leaves the current level.
- Exiting the top-level menu returns to the normal TPS prompt.

## Install discovery

TPS **help** (`h`) includes a **TPS ADDONS** section at the **bottom** (after project commands). It is **not** listed on the compact `TPS CMDS:` hint line.

| Command | Abbrev | Action |
|---------|--------|--------|
| install utilities | iu | Clone addon, link project, merge `tps-addon.yaml` commands |

Flow:

1. User runs `iu` from any project.
2. macOS: **Finder** choose parent directory; other OS: type parent path.
3. `git clone <clone_url>` into `<parent>/tps-utilities` (skip clone if directory exists and is already the addon).
4. Link with `tps project add` equivalent and merge manifest commands onto that project.
5. User switches to the utilities project (`p` / `pN`) to run addon tools.

## Direct clone without TPS

Cloning only the addon repo does not register it in TPS. The addon **`setup`** script must:

1. If `tps` is not on `PATH`, print how to install TPS (`https://github.com/Cimmerial/tps`, `go install`, `tps setup`).
2. Optionally offer to `git clone` the TPS repo (sibling directory or Finder-chosen parent on macOS).
3. When `tps` is available, tell the user to run `tps setup --path <addon-root>` or use `iu` from an existing session.

Document the same at the top of the addon README.

## Worked example: tps-utilities / UIC

**Project:** linked `tps-utilities` repo.

**Command:** `universal image conversion` / `uic` — no `tags`. Only available when that project is selected.

**Tool:** Universal Image Conversion (UIC), menu `menus/uic.yaml`:

- **convert** — pick files, pick output folder (optional → default), target format chosen from submenu.
- **mass convert** — pick folder, target format from submenu.

### Image formats

| Format | Extensions |
|--------|------------|
| JPEG | `.jpg`, `.jpeg` |
| PNG | `.png` |
| GIF | `.gif` |
| WebP | `.webp` |
| SVG | `.svg` |
| HEIC | `.heic` |
| AVIF | `.avif` |
| TIFF | `.tif`, `.tiff` |
| BMP | `.bmp` |
| RAW | `.raw`, `.cr2`, `.nef`, `.arw` |

### PNG hub

All conversion goes **source → PNG → target**. Skip a hop when source or target is already PNG. Do **not** implement or test every pairwise format; test **each format ↔ PNG** only.

Implementation uses ImageMagick (`magick`). Probe coders at menu time via `dim_when` / `tpsu img probe`.

| Rule | UI |
|------|-----|
| Cannot read format | Dim as **source** paths (mass) or dim “from” where applicable |
| Cannot write format | Dim **→ FORMAT** targets |
| RAW | Never a write target (always dimmed) |
| SVG write | Dim if not a real vector export; SVG remains a valid read (rasterize via PNG hub) |
| GIF via PNG | First frame only; conversion output says so |

### Mass convert

Walk folder for known image extensions, convert each file through the hub, print:

- `ok/attempted` success rate
- Count of skipped non-images
- Up to **5 random** failed basenames (all failures if fewer than five)

## TPS binary version

`min_tps` is an optional **semver hint** for the tps CLI build itself (separate from `contract_iteration`). Bump it when an addon depends on tps features not tied to the addon contract doc.
