package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const defaultCaptureTimeout = 90 * time.Second

func LogFile(projectID string, id int) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, "tps", "jobs", "logs", fmt.Sprintf("%s-%d.log", projectID, id))
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return "", err
	}
	return p, nil
}

// CaptureOptions configures log watching and post-match actions.
type CaptureOptions struct {
	LogPath     string
	Pattern     string
	Actions     []string
	SMSTo       string
	SMSShortcut string
	Timeout     time.Duration
}
