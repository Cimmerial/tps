package jobs

import "runtime"

// launchReady reports whether spawnInTerminal has written enough state to register the job.
func launchReady(pid int, tty string) bool {
	if pid <= 0 {
		return false
	}
	if runtime.GOOS == "darwin" {
		return tty != ""
	}
	return true
}
