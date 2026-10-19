package cli

import (
	"path/filepath"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/project"
)

type linkTarget struct {
	Name string
	Path string
}

func linkTargetsForCopy(cfg *config.Config, selected *config.Project, cwd string) []linkTarget {
	seen := map[string]bool{}
	var out []linkTarget

	addFromDir := func(dir string, nameHint string) {
		if dir == "" {
			return
		}
		root, err := project.ProjectRootFromCwd(dir)
		if err != nil || root == "" || canSelfUpdate(root) {
			return
		}
		clean, err := config.CleanPath(root)
		if err != nil || seen[clean] {
			return
		}
		seen[clean] = true
		name := nameHint
		if name == "" {
			name = filepath.Base(clean)
		}
		if cfg != nil {
			if p := cfg.FindByPath(clean); p != nil && p.Name != "" {
				name = p.Name
			}
		}
		out = append(out, linkTarget{Name: name, Path: clean})
	}

	if cwd != "" {
		addFromDir(cwd, "")
	}
	if selected != nil && !canSelfUpdate(selected.Path) {
		addFromDir(selected.Path, selected.Name)
	}
	if cfg != nil {
		for _, p := range cfg.Projects {
			if canSelfUpdate(p.Path) {
				continue
			}
			addFromDir(p.Path, p.Name)
		}
	}
	return out
}
