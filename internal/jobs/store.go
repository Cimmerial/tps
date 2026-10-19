package jobs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

type Job struct {
	ID      int       `yaml:"id"`
	Alias   string    `yaml:"alias"`
	Shell   string    `yaml:"shell"`
	PID     int       `yaml:"pid"`
	Started time.Time `yaml:"started"`
}

type jobFile struct {
	Jobs []Job `yaml:"jobs"`
}

func jobsPath(projectID string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tps", "jobs", projectID+".yaml"), nil
}

func load(projectID string) ([]Job, error) {
	p, err := jobsPath(projectID)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var f jobFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return f.Jobs, nil
}

func save(projectID string, jobs []Job) error {
	p, err := jobsPath(projectID)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(jobFile{Jobs: jobs})
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}

func alive(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := syscall.Kill(pid, 0)
	return err == nil
}

func Reap(projectID string) ([]Job, error) {
	jobs, err := load(projectID)
	if err != nil {
		return nil, err
	}
	kept := make([]Job, 0, len(jobs))
	for _, j := range jobs {
		if alive(j.PID) {
			kept = append(kept, j)
		}
	}
	if len(kept) != len(jobs) {
		if err := save(projectID, kept); err != nil {
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

func Add(projectID string, job Job) error {
	jobs, err := load(projectID)
	if err != nil {
		return err
	}
	jobs = append(jobs, job)
	return save(projectID, jobs)
}

func Kill(projectID string, id int) (int, error) {
	jobs, err := Reap(projectID)
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return killAll(projectID, jobs)
	}
	var kept []Job
	killed := 0
	for _, j := range jobs {
		if j.ID == id {
			if terminate(j.PID) {
				killed++
			}
			continue
		}
		kept = append(kept, j)
	}
	if killed == 0 {
		return 0, fmt.Errorf("no terminal job #%d", id)
	}
	return killed, save(projectID, kept)
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
		if terminate(j.PID) {
			killed++
		}
	}
	if err := save(projectID, nil); err != nil {
		return killed, err
	}
	return killed, nil
}

func terminate(pid int) bool {
	if pid <= 0 {
		return false
	}
	if !alive(pid) {
		return false
	}
	_ = syscall.Kill(pid, syscall.SIGTERM)
	_ = syscall.Kill(-pid, syscall.SIGTERM)
	time.Sleep(120 * time.Millisecond)
	if alive(pid) {
		_ = syscall.Kill(pid, syscall.SIGKILL)
		_ = syscall.Kill(-pid, syscall.SIGKILL)
	}
	return true
}

func PIDFile(projectID string, id int) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tps", "jobs", "pids", projectID+"-"+fmt.Sprint(id)+".pid"), nil
}
