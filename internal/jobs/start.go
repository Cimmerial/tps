package jobs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Cimmerial/tps/internal/config"
)

func Start(p *config.Project, alias string) (Job, error) {
	jobs, err := StartAll(p, alias)
	if err != nil {
		return Job{}, err
	}
	if len(jobs) == 0 {
		return Job{}, fmt.Errorf("no shell lines for %q", alias)
	}
	return jobs[0], nil
}

// StartAll launches each resolved shell line in its own terminal (one history entry per alias).
func StartAll(p *config.Project, alias string) ([]Job, error) {
	cmd := p.FindCommand(alias)
	if cmd == nil {
		return nil, fmt.Errorf("unknown command %q", alias)
	}
	shells, err := cmd.ShellsWithInputs(nil)
	if err != nil {
		return nil, err
	}
	return StartAllWithShells(p, alias, shells)
}

// StartAllWithShells launches pre-resolved shell lines (e.g. after [INPUTn] substitution).
func StartAllWithShells(p *config.Project, alias string, shells []string) ([]Job, error) {
	if len(shells) == 0 {
		return nil, fmt.Errorf("command %q has no shell", alias)
	}
	out := make([]Job, 0, len(shells))
	for i, shell := range shells {
		job, err := launchShell(p, alias, shell, i == 0)
		if err != nil {
			return out, err
		}
		out = append(out, job)
	}
	return out, nil
}

func launch(p *config.Project, alias string, recordHistory bool) (Job, error) {
	cmd := p.FindCommand(alias)
	if cmd == nil {
		return Job{}, fmt.Errorf("unknown command %q", alias)
	}
	shells := cmd.ResolvedShells()
	if len(shells) == 0 {
		return Job{}, fmt.Errorf("command %q has no shell", alias)
	}
	return launchShell(p, alias, shells[0], recordHistory)
}

func launchShell(p *config.Project, alias, shell string, recordHistory bool) (Job, error) {
	cmd := p.FindCommand(alias)
	if cmd == nil {
		return Job{}, fmt.Errorf("unknown command %q", alias)
	}
	if _, err := Reap(p.ID); err != nil {
		return Job{}, err
	}
	id, err := allocateJobID(p.ID)
	if err != nil {
		return Job{}, err
	}
	pidPath, err := PIDFile(p.ID, id)
	if err != nil {
		return Job{}, err
	}
	ttyPath, err := TTYFile(p.ID, id)
	if err != nil {
		return Job{}, err
	}
	pgidPath, err := PGIDFile(p.ID, id)
	if err != nil {
		return Job{}, err
	}
	_ = os.Remove(pidPath)
	_ = os.Remove(ttyPath)
	_ = os.Remove(pgidPath)

	logPath := ""
	if cmd.Capture != "" {
		logPath, err = LogFile(p.ID, id)
		if err != nil {
			return Job{}, err
		}
		_ = os.Remove(logPath)
	}

	closeOnComplete := cmd != nil && cmd.CloseOnComplete
	if err := spawnInTerminal(p.Path, shell, pidPath, logPath, ttyPath); err != nil {
		return Job{}, err
	}

	pid := 0
	tty := ""
	for i := 0; i < 60; i++ {
		time.Sleep(50 * time.Millisecond)
		if pid == 0 {
			if data, err := os.ReadFile(pidPath); err == nil {
				pid = parsePID(string(data))
			}
		}
		if tty == "" {
			if data, err := os.ReadFile(ttyPath); err == nil {
				tty = strings.TrimSpace(string(data))
			}
		}
		if launchReady(pid, tty) {
			break
		}
	}
	if pid == 0 {
		return Job{}, fmt.Errorf("started terminal but could not read process id (check terminal or logs)")
	}

	job := Job{
		ID:      id,
		Alias:   alias,
		Shell:   shell,
		PID:     pid,
		TTY:     tty,
		Started: time.Now(),
	}
	if err := Add(p.ID, job); err != nil {
		return Job{}, err
	}
	if recordHistory {
		if err := appendHistory(p.ID, alias); err != nil {
			return job, err
		}
	}
	if closeOnComplete {
		_ = spawnCloseWatcherScript(p.ID, job)
		notifyDetachedComplete(p.ID, job)
	}
	return job, nil
}

func parsePID(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
