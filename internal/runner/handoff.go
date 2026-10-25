package runner

import (
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

// HandoffScript runs script in this terminal and does not return to tps on success.
// Unix replaces the process (syscall.Exec); Windows runs the script then exits.
func HandoffScript(p *config.Project, script string) error {
	script = strings.TrimSpace(script)
	if script == "" {
		return nil
	}
	return handoffScript(p, script)
}
