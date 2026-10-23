# Instance readme source

The canonical text for linked app instance docs lives at:

- **[internal/instance/README.md](../internal/instance/README.md)** → `TPS.md` (cheat sheet)
- **[internal/instance/SPECS.md](../internal/instance/SPECS.md)** → `TPS_SPECS.md` (command/script spec)

Both are embedded in release binaries. Edit the source file when behavior or the execution contract changes, then ship tps so `tps setup` / `u` refresh copies.
