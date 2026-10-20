package jobs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Job struct {
	ID      int       `yaml:"id"`
	Alias   string    `yaml:"alias"`
	Shell   string    `yaml:"shell"`
	PID     int       `yaml:"pid"`
	TTY     string    `yaml:"tty,omitempty"`
	Started time.Time `yaml:"started"`
}

type jobFile struct {
	Jobs      []Job    `yaml:"jobs"`
	History   []string `yaml:"history,omitempty"`
	NextJobID int      `yaml:"next_job_id,omitempty"`
}

func jobsPath(projectID string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tps", "jobs", projectID+".yaml"), nil
}

func readFile(projectID string) (jobFile, error) {
	p, err := jobsPath(projectID)
	if err != nil {
		return jobFile{}, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return jobFile{}, nil
		}
		return jobFile{}, err
	}
	var f jobFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return jobFile{}, err
	}
	return f, nil
}

func writeFile(projectID string, f jobFile) error {
	p, err := jobsPath(projectID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}

func Reap(projectID string) ([]Job, error) {
	f, err := readFile(projectID)
	if err != nil {
		return nil, err
	}
	kept := make([]Job, 0, len(f.Jobs))
	for _, j := range f.Jobs {
		if alive(j.PID) {
			kept = append(kept, j)
		}
	}
	if len(kept) != len(f.Jobs) {
		f.Jobs = kept
		if err := writeFile(projectID, f); err != nil {
			return kept, err
		}
	}
	sort.Slice(kept, func(i, j int) bool { return kept[i].ID < kept[j].ID })
	return kept, nil
}

func nextID(jobs []Job) int {
	max := 0
	for _, j := range jobs {
		if j.ID > max {
			max = j.ID
		}
	}
	return max + 1
}

// allocateJobID returns a monotonic job id (never reused after reap) so Terminal .command titles stay unique.
func allocateJobID(projectID string) (int, error) {
	f, err := readFile(projectID)
	if err != nil {
		return 0, err
	}
	id := f.NextJobID
	if id <= 0 {
		id = nextID(f.Jobs)
	}
	f.NextJobID = id + 1
	if err := writeFile(projectID, f); err != nil {
		return 0, err
	}
	return id, nil
}

func appendHistory(projectID, alias string) error {
	f, err := readFile(projectID)
	if err != nil {
		return err
	}
	f.History = append(f.History, alias)
	return writeFile(projectID, f)
}

func LaunchHistory(projectID string) ([]string, error) {
	f, err := readFile(projectID)
	if err != nil {
		return nil, err
	}
	return f.History, nil
}

func Add(projectID string, job Job) error {
	f, err := readFile(projectID)
	if err != nil {
		return err
	}
	f.Jobs = append(f.Jobs, job)
	return writeFile(projectID, f)
}

// Kill stops a running job by its 1-based index in the `t` list (not the internal job id).
func Kill(projectID string, displayNum int) (int, error) {
	jobs, err := Reap(projectID)
	if err != nil {
		return 0, err
	}
	if displayNum == 0 {
		return killAll(projectID, jobs)
	}
	if displayNum < 1 || displayNum > len(jobs) {
		return 0, fmt.Errorf("no terminal job #%d", displayNum)
	}
	target := jobs[displayNum-1]
	f, err := readFile(projectID)
	if err != nil {
		return 0, err
	}
	var kept []Job
	for _, j := range jobs {
		if j.ID == target.ID {
			stopJob(projectID, j)
			continue
		}
		kept = append(kept, j)
	}
	f.Jobs = kept
	return 1, writeFile(projectID, f)
}

// ListDisplayNumber returns the 1-based `t` list index for an internal job id.
func ListDisplayNumber(list []Job, jobID int) int {
	for i, j := range list {
		if j.ID == jobID {
			return i + 1
		}
	}
	return 0
}

func KillAll(projectID string) (int, error) {
	jobs, err := Reap(projectID)
	if err != nil {
		return 0, err
	}
	return killAll(projectID, jobs)
}

func killAll(projectID string, jobs []Job) (int, error) {
	killed := 0
	for _, j := range jobs {
		if stopJob(projectID, j) {
			killed++
		}
	}
	f, err := readFile(projectID)
	if err != nil {
		return killed, err
	}
	f.Jobs = nil
	if err := writeFile(projectID, f); err != nil {
		return killed, err
	}
	return killed, nil
}

func stopJob(projectID string, j Job) bool {
	killJob(projectID, j)
	return true
}

func killJob(projectID string, j Job) {
	tty := strings.TrimSpace(j.TTY)
	if tty == "" {
		if path, err := TTYFile(projectID, j.ID); err == nil {
			if data, err := os.ReadFile(path); err == nil {
				tty = strings.TrimSpace(string(data))
			}
		}
	}
	_, _ = killProcessTree(j.PID)
	CloseTerminalTab(tty, JobWindowMarker(projectID, j.ID))
	killPortsFromShell(j.Shell)
	cleanupJobFiles(projectID, j.ID)
}

func cleanupJobFiles(projectID string, jobID int) {
	if path, err := TTYFile(projectID, jobID); err == nil {
		_ = os.Remove(path)
	}
	if path, err := PGIDFile(projectID, jobID); err == nil {
		_ = os.Remove(path)
	}
	if path, err := PIDFile(projectID, jobID); err == nil {
		_ = os.Remove(path)
	}
}

func PIDFile(projectID string, id int) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tps", "jobs", "pids", projectID+"-"+fmt.Sprint(id)+".pid"), nil
}

func TTYFile(projectID string, id int) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tps", "jobs", "pids", projectID+"-"+fmt.Sprint(id)+".tty"), nil
}

func PGIDFile(projectID string, id int) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tps", "jobs", "pids", projectID+"-"+fmt.Sprint(id)+".pgid"), nil
}
