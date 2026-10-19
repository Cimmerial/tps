package runner

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

// RunCommand runs the project alias. It returns the child exit code (0 on success).
func RunCommand(p *config.Project, alias string) (int, error) {
	cmd := p.FindCommand(alias)
	if cmd == nil {
		return 1, fmt.Errorf("unknown command %q for project %q", alias, p.Name)
	}
	shells := cmd.ResolvedShells()
	if len(shells) == 0 {
		return 1, fmt.Errorf("command %q has no shell", alias)
	}
	return RunScript(p, strings.Join(shells, " && "))
}

// RunScript runs shell in the project directory (stdin/stdout/stderr attached).
func RunScript(p *config.Project, script string) (int, error) {
	script = strings.TrimSpace(script)
	if script == "" {
		return 0, nil
	}
	c := exec.Command("sh", "-c", script)
	c.Dir = p.Path
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return exit.ExitCode(), nil
		}
		return 1, err
	}
	return 0, nil
}
