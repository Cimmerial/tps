//go:build !darwin && !linux && !windows

package jobs

import "fmt"

func spawnInTerminal(workDir, shellCmd, pidFile, logFile, ttyFile string, focus bool) error {
	return fmt.Errorf("detached terminals are not supported on this OS")
}
