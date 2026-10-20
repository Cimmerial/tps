//go:build linux

package jobs

import "github.com/Cimmerial/tps/internal/config"

// CloseAllTerminalWindows stops every tps-tracked detached job (all linked projects).
func CloseAllTerminalWindows() CloseAllTermsResult {
	cfg, _ := config.Load()
	labels := RegisteredJobLabels(cfg)
	n, _ := KillAllRegistered(cfg)
	if n == 0 {
		return CloseAllTermsResult{}
	}
	if len(labels) > n {
		labels = labels[:n]
	}
	return CloseAllTermsResult{Count: n, Labels: labels}
}
