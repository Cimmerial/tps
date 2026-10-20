//go:build windows

package update

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	envStagedFrom   = "TPS_APPLY_STAGED_FROM"
	envStagedTo     = "TPS_APPLY_STAGED_TO"
	FlagApplyStaged = "--tps-apply-staged"
)

// ApplyStagedWindowsBinaryIfNeeded swaps tps.exe.new into place (after u or install).
func ApplyStagedWindowsBinaryIfNeeded() error {
	to, from := stagedWindowsPaths()
	if from == "" || to == "" {
		return nil
	}
	handoff := os.Getenv(envStagedFrom) != ""
	if handoff {
		time.Sleep(350 * time.Millisecond)
	}
	tries := 3
	pause := 50 * time.Millisecond
	if handoff {
		tries = 25
		pause = 150 * time.Millisecond
	}
	var lastErr error
	for i := 0; i < tries; i++ {
		err := os.Rename(from, to)
		if err == nil {
			_ = os.Unsetenv(envStagedFrom)
			_ = os.Unsetenv(envStagedTo)
			return nil
		}
		lastErr = err
		time.Sleep(pause)
	}
	return fmt.Errorf("apply staged binary %s -> %s: %w", from, to, lastErr)
}

func stagedWindowsPaths() (to, from string) {
	if f := os.Getenv(envStagedFrom); f != "" {
		return os.Getenv(envStagedTo), f
	}
	exe, err := os.Executable()
	if err != nil {
		return "", ""
	}
	to, err = filepath.Abs(exe)
	if err != nil {
		to = exe
	}
	if resolved, err := filepath.EvalSymlinks(to); err == nil {
		to = resolved
	}
	from = to + ".new"
	if _, err := os.Stat(from); err != nil {
		return "", ""
	}
	return to, from
}
