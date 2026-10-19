//go:build !darwin

package jobs

import "fmt"

func spawnInTerminal(workDir, shellCmd, pidFile, logFile string) error {
	return fmt.Errorf("detached terminals are only supported on macOS")
}
