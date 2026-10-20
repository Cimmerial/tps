//go:build windows

package update

import (
	"fmt"
	"os"
	"time"
)

const (
	envStagedFrom = "TPS_APPLY_STAGED_FROM"
	envStagedTo   = "TPS_APPLY_STAGED_TO"
	FlagApplyStaged = "--tps-apply-staged"
)

// ApplyStagedWindowsBinaryIfNeeded swaps tps.exe.new into place after the parent exits.
func ApplyStagedWindowsBinaryIfNeeded() error {
	from := os.Getenv(envStagedFrom)
	to := os.Getenv(envStagedTo)
	if from == "" || to == "" {
		return nil
	}
	time.Sleep(350 * time.Millisecond)
	var lastErr error
	for i := 0; i < 25; i++ {
		if err := os.Rename(from, to); err == nil {
			_ = os.Unsetenv(envStagedFrom)
			_ = os.Unsetenv(envStagedTo)
			return nil
		}
		lastErr = err
		time.Sleep(150 * time.Millisecond)
	}
	return fmt.Errorf("apply staged binary %s -> %s: %w", from, to, lastErr)
}
