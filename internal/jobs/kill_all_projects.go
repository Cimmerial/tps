package jobs

import (
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

// RegisteredJobLabels returns display labels for every active detached job (all projects).
func RegisteredJobLabels(cfg *config.Config) []string {
	if cfg == nil {
		return nil
	}
	var out []string
	for _, p := range cfg.Projects {
		if p.ID == "" {
			continue
		}
		list, err := Reap(p.ID)
		if err != nil {
			continue
		}
		for _, j := range list {
			out = append(out, jobDisplayLabel(j))
		}
	}
	return out
}

func jobDisplayLabel(j Job) string {
	if a := strings.TrimSpace(j.Alias); a != "" {
		return a
	}
	s := strings.TrimSpace(j.Shell)
	if len(s) > 48 {
		return s[:45] + "..."
	}
	if s != "" {
		return s
	}
	return "terminal"
}

// KillAllRegistered kills every detached terminal job tracked by tps across all linked projects.
func KillAllRegistered(cfg *config.Config) (int, error) {
	if cfg == nil {
		return 0, nil
	}
	total := 0
	for _, p := range cfg.Projects {
		if p.ID == "" {
			continue
		}
		n, err := KillAll(p.ID)
		if err != nil {
			return total, err
		}
		total += n
	}
	return total, nil
}
