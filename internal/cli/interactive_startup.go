package cli

import (
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/project"
)

type interactiveProjectChoice struct {
	ProjectID  string
	SaveConfig bool
	Pick       bool
	NoProjects bool
}

// resolveInteractiveProject picks the session project without auto-linking new repos.
// Restart flag and saved session_project_id override cwd when still valid.
func resolveInteractiveProject(cfg *config.Config, cwd string, restartProjectID string) (interactiveProjectChoice, error) {
	if id := strings.TrimSpace(restartProjectID); id != "" {
		if cfg.FindByID(id) != nil {
			return interactiveProjectChoice{ProjectID: id}, nil
		}
	}
	prefs, err := config.LoadPreferences()
	if err == nil && prefs != nil {
		if id := strings.TrimSpace(prefs.SessionProjectID); id != "" {
			if cfg.FindByID(id) != nil {
				return interactiveProjectChoice{ProjectID: id}, nil
			}
		}
	}

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
