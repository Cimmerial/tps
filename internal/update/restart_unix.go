//go:build !windows

package update

import (
	"fmt"
	"os"
	"syscall"
)

// RestartInteractive replaces this process with a fresh tps invocation (same args).
func RestartInteractive(updated bool) error {
	bin, argv, err := buildRestartArgv()
	if err != nil {
		return fmt.Errorf("locate tps: %w", err)
	}
	if updated {
		argv = append(argv, flagUpdated)
	} else {
		argv = append(argv, flagUpToDate)
	}
	return syscall.Exec(bin, argv, os.Environ())
}
