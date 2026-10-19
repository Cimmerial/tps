# tps

Personal project command hub — short aliases for shell commands, per linked repo.

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

Config path: OS user config dir + `/tps/config.yaml` (see `tps ih`).
