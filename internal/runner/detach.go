package runner

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

func ShouldDetach(c config.Command) bool {
	detached, _ := DetachInfo(c)
	return detached
}

// DetachInfo reports whether the command opens a new terminal and a short reason for the UI.
func DetachInfo(c config.Command) (detached bool, hint string) {
	if c.ForceForeground {
		return false, "runs in this session (foreground)"
	}
	if len(c.ResolvedShells()) > 1 {
		n := len(c.ResolvedShells())
		return true, fmt.Sprintf("opens in a new terminal (%d parallel shells)", n)
	}
	if c.Detach {
		return true, "opens in a new terminal (detach: true)"
	}
	if c.CaptureActive() {
		return true, "opens in a new terminal (capture)"
	}
	if shellAutoDetaches(primaryShell(c)) {
		return true, "opens in a new terminal (auto: dev / long-running)"
	}
	return false, "runs in this session (foreground)"
}

func primaryShell(c config.Command) string {
	s := strings.ToLower(c.Shell)
	if s == "" {
		shells := c.ResolvedShells()
		if len(shells) > 0 {
			s = strings.ToLower(shells[0])
		}
	}
	return s
}

func shellAutoDetaches(s string) bool {
	switch {
	case strings.Contains(s, "vite"),
		strings.Contains(s, "npm run dev"),
		strings.Contains(s, "pnpm dev"),
		strings.Contains(s, "yarn dev"),
		strings.Contains(s, "dev:mobile"),
		strings.Contains(s, "next dev"),
		strings.Contains(s, "webpack serve"):
		return true
	default:
		return false
	}
}
