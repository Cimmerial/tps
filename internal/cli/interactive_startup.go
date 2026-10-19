package cli

import (
	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/project"
)

type interactiveProjectChoice struct {
	ProjectID  string
	SaveConfig bool
	Pick       bool
	NoProjects bool
}

// resolveInteractiveProject picks the session project from cwd without auto-linking new repos.
func resolveInteractiveProject(cfg *config.Config, cwd string) (interactiveProjectChoice, error) {
	res, resolveErr := project.ResolveFromCwd(cfg, cwd)
	if resolveErr == nil {
		return interactiveProjectChoice{
			ProjectID:  res.Project.ID,
			SaveConfig: res.Healed,
		}, nil
	}
	if len(cfg.Projects) == 0 {
		return interactiveProjectChoice{NoProjects: true}, nil
	}
	return interactiveProjectChoice{Pick: true}, nil
}
