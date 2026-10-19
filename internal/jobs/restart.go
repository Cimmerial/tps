package jobs

import (
	"fmt"
	"time"

	"github.com/Cimmerial/tps/internal/config"
)

func RestartAll(p *config.Project) ([]Job, error) {
	hist, err := LaunchHistory(p.ID)
	if err != nil {
		return nil, err
	}
	if len(hist) == 0 {
		return nil, fmt.Errorf("no terminal launch history (start a detached command first)")
	}
	_, _ = KillAll(p.ID)
	time.Sleep(200 * time.Millisecond)
	out := make([]Job, 0, len(hist))
	for _, alias := range hist {
		jobs, err := StartAll(p, alias)
		if err != nil {
			return out, fmt.Errorf("restart %q: %w", alias, err)
		}
		out = append(out, jobs...)
		time.Sleep(150 * time.Millisecond)
	}
	return out, nil
}
