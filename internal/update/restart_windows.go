//go:build windows

package update

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Cimmerial/tps/internal/ui"
)

// RestartInteractive starts a fresh tps in this console and exits the current process.
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

	args := argv[1:]
	env := os.Environ()
	if pairs := pendingWindowsReplacements(); len(pairs) > 0 {
		p := pairs[0]
		env = append(env, envStagedFrom+"="+p.Staging, envStagedTo+"="+p.Dest)
		args = append(args, FlagApplyStaged)
		clearWindowsReplaceQueue()
	}

	ui.RestoreStdin()
	fmt.Print("\r\n")

	cmd := exec.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("restart tps: %w", err)
	}
	os.Exit(0)
	return nil
}
