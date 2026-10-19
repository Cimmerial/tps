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
	_ = os.Remove(pidPath)

	logPath := ""
	if cmd.Capture != "" {
		logPath, err = LogFile(p.ID, id)
		if err != nil {
			return Job{}, err
		}
		_ = os.Remove(logPath)
	}

	if err := spawnInTerminal(p.Path, cmd.Shell, pidPath, logPath); err != nil {
		return Job{}, err
	}

	pid := 0
	for i := 0; i < 40; i++ {
		time.Sleep(50 * time.Millisecond)
		data, err := os.ReadFile(pidPath)
		if err == nil {
			pid = parsePID(string(data))
			if pid > 0 {
				break
			}
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
		Started: time.Now(),
	}
	if err := Add(p.ID, job); err != nil {
		return Job{}, err
	}
	return job, nil
}

func parsePID(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
