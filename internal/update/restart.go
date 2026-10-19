package update

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/Cimmerial/tps/internal/config"
)

const flagUpdated = "--tps-updated"
const flagUpToDate = "--tps-uptodate"
const flagSessionProjectPrefix = "--tps-session-project="

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
	sessionProject := ""
	if prefs, err := config.LoadPreferences(); err == nil && prefs != nil {
		sessionProject = strings.TrimSpace(prefs.SessionProjectID)
	}
	for i := 1; i < len(os.Args); i++ {
		a := os.Args[i]
		if a == flagUpdated || a == flagUpToDate {
			continue
		}
		if strings.HasPrefix(a, flagSessionProjectPrefix) {
			continue
		}
		argv = append(argv, a)
	}
	if sessionProject != "" {
		argv = append(argv, flagSessionProjectPrefix+sessionProject)
	}
	if updated {
		argv = append(argv, flagUpdated)
	} else {
		argv = append(argv, flagUpToDate)
	}
	return syscall.Exec(bin, argv, os.Environ())
}
