package jobs

import (
	"os"
	"strings"
	"sync"
	"time"
)

var (
	detachedCompleteMu sync.Mutex
	onDetachedComplete func(Job, time.Duration)
)

// OnDetachedJobComplete registers a callback when a close_on_complete job's tracked process exits.
func OnDetachedJobComplete(fn func(Job, time.Duration)) {
	detachedCompleteMu.Lock()
	onDetachedComplete = fn
	detachedCompleteMu.Unlock()
}

func notifyDetachedComplete(projectID string, job Job) {
	go waitDetachedComplete(projectID, job)
}

func waitDetachedComplete(projectID string, job Job) {
	start := job.Started
	if start.IsZero() {
		start = time.Now()
	}
	deadline := time.Now().Add(48 * time.Hour)
	for time.Now().Before(deadline) {
		if !alive(job.PID) {
			elapsed := time.Since(start)
			marker := JobWindowMarker(projectID, job.ID)
			tryCloseDetachedTab(jobTTY(projectID, job), marker)
			cleanupJobFiles(projectID, job.ID)
			_, _ = Reap(projectID)
			detachedCompleteMu.Lock()
			fn := onDetachedComplete
			detachedCompleteMu.Unlock()
			if fn != nil {
				fn(job, elapsed)
			}
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func jobTTY(projectID string, job Job) string {
	tty := strings.TrimSpace(job.TTY)
	if tty != "" {
		return tty
	}
	path, err := TTYFile(projectID, job.ID)
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
