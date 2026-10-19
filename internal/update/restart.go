package update

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const flagUpdated = "--tps-updated"
const flagUpToDate = "--tps-uptodate"

// RestartInteractive replaces this process with a fresh tps invocation (same args).
func RestartInteractive(updated bool) error {
	bin, err := exec.LookPath("tps")
	if err != nil {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("locate tps: %w", err)
		}
		bin = exe
	}

	argv := []string{bin}
	for i := 1; i < len(os.Args); i++ {
		a := os.Args[i]
		if a == flagUpdated || a == flagUpToDate {
			continue
		}
		argv = append(argv, a)
	}
	if updated {
		argv = append(argv, flagUpdated)
	} else {
		argv = append(argv, flagUpToDate)
	}
	return syscall.Exec(bin, argv, os.Environ())
}
