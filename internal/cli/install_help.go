package cli

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func devRepoRoot(cfg *config.Config, selected *config.Project) string {
	if selected != nil && canSelfUpdate(selected.Path) {
		return selected.Path
	}
	if cfg != nil {
		for _, p := range cfg.Projects {
			if canSelfUpdate(p.Path) {
				return p.Path
			}
		}
	}
	return ""
}

func printInstallHelp() {
	cfgPath, err := config.Path()
	if err != nil {
		cfgPath = "(os user config dir)/tps/config.yaml"
	}

	fmt.Println(ui.Bold("Install tps"))
	fmt.Println(ui.Dim("  macOS / Linux:"))
	fmt.Printf("    %s\n", ui.Cyan("curl -fsSL https://raw.githubusercontent.com/Cimmerial/tps/main/scripts/install.sh | bash"))
	fmt.Println(ui.Dim("  Windows (PowerShell):"))
	fmt.Printf("    %s\n", ui.Cyan("irm https://raw.githubusercontent.com/Cimmerial/tps/main/scripts/install.ps1 | iex"))
	fmt.Println(ui.Dim("  developers:"))
	fmt.Printf("    %s\n", ui.Cyan("go install github.com/Cimmerial/tps/cmd/tps@latest"))

	fmt.Println()
	fmt.Println(ui.Bold("Link an app repo"))
	fmt.Printf("    %s\n", ui.Cyan("cd your-app && tps setup"))
	fmt.Printf("    %s\n", ui.Cyan("tps project add"))
	fmt.Println(ui.Dim("  writes TPS.md in the repo (how tps works here + limits)"))

	fmt.Println()
	fmt.Println(ui.Bold("Update"))
	fmt.Printf("    %s\n", ui.Cyan("u"))
	fmt.Println(ui.Dim("  dev clone linked: git pull, rebuild, sync every tps on PATH, refresh TPS.md"))
	fmt.Println(ui.Dim("  otherwise: latest GitHub release + same binary sync"))
	fmt.Println(ui.Dim("  opening tps also quietly rebuilds when the dev repo changed (restart if needed)"))

	fmt.Println()
	fmt.Println(ui.Bold("Config"))
	fmt.Printf("    %s\n", ui.Dim(cfgPath))
	fmt.Println()
}
