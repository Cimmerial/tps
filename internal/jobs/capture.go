package jobs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Cimmerial/tps/internal/ui"
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

// WatchCapture polls logPath until pattern matches, then prints and copies the match.
func WatchCapture(logPath, pattern string, timeout time.Duration) {
	if timeout <= 0 {
		timeout = defaultCaptureTimeout
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("capture: invalid pattern: "+err.Error()))
		return
	}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		data, err := os.ReadFile(logPath)
		if err == nil && len(data) > 0 {
			if m := re.FindString(string(data)); m != "" {
				fmt.Print("\n")
				fmt.Println(m)
				if err := ui.CopyToClipboard(m); err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("capture: clipboard: "+err.Error()))
					return
				}
				fmt.Println(ui.Green("copied to clipboard"))
				return
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
	fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("capture: no match yet (still running in terminal)"))
}
