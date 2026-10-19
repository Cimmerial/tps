package update

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// RestartInteractive replaces this process with a fresh tps invocation (same args).
func RestartInteractive() error {
	bin, err := exec.LookPath("tps")
	if err != nil {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("locate tps: %w", err)
		}
		bin = exe
	}
	argv := make([]string, len(os.Args))
	argv[0] = bin
	for i := 1; i < len(os.Args); i++ {
		argv[i] = os.Args[i]
	}
	if len(argv) == 1 {
		argv = []string{bin}
	}
	return syscall.Exec(bin, argv, os.Environ())
}
