package cli

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/addon"
	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/version"
)

func printTPSVersion(p *config.Project) {
	fmt.Println("tps", version.Label())
	if p == nil || !addon.IsAddonProject(p.Path) {
		return
	}
	m, err := addon.LoadManifest(p.Path)
	if err != nil {
		return
	}
	name := strings.TrimSpace(m.Name)
	if name == "" {
		name = strings.TrimSpace(m.ID)
	}
	if name == "" {
		name = "addon"
	}
	ver := addon.ProjectVersion(p.Path)
	if ver == "" {
		fmt.Printf("addon %s (no version in tps-addon.yaml)\n", name)
		return
	}
	fmt.Printf("addon %s %s\n", name, ver)
}
