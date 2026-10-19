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
	return launch(p, alias, true)
}

func launch(p *config.Project, alias string, recordHistory bool) (Job, error) {
	cmd := p.FindCommand(alias)
	if cmd == nil {
		return Job{}, fmt.Errorf("unknown command %q", alias)
	}
	jobs, err := Reap(p.ID)
	if err != nil {
		return Job{}, err
	}
	id := nextID(jobs)
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

	if err := spawnInTerminal(p.Path, cmd.Shell, pidPath, logPath, ttyPath); err != nil {
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
		if pid > 0 && tty != "" {
			break
		}
	}
	if pid == 0 {
		return Job{}, fmt.Errorf("started terminal but could not read process id (check Terminal.app)")
	}

	job := Job{
		ID:      id,
		Alias:   alias,
		Shell:   cmd.Shell,
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
	return job, nil
}

func parsePID(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
