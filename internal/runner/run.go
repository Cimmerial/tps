package runner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
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
	code, _, err := RunScriptCaptureStderr(p, script)
	return code, err
}

// RunScriptCaptureStderr runs shell like RunScript and returns a copy of stderr (still streamed live).
func RunScriptCaptureStderr(p *config.Project, script string) (int, string, error) {
	script = strings.TrimSpace(script)
	if script == "" {
		return 0, "", nil
	}
	c := scriptCommand(script)
	c.Dir = p.Path
	c.Env = scriptEnviron(p)
	c.Stdin = os.Stdin
	var buf bytes.Buffer
	if ui.SimplifiedOutput() {
		c.Stdout = &buf
		c.Stderr = &buf
	} else {
		c.Stdout = os.Stdout
		c.Stderr = io.MultiWriter(os.Stderr, &buf)
	}
	if err := c.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			if ui.SimplifiedOutput() {
				sum := summarizeFailure(script, buf.String(), exit.ExitCode())
				ui.PrintOutcome(false, sum.label, sum.detail)
				_, _ = os.Stdout.Write(buf.Bytes())
			}
			return exit.ExitCode(), buf.String(), nil
		}
		return 1, buf.String(), err
	}
	if ui.SimplifiedOutput() {
		sum := summarizeSuccess(script, buf.String())
		ui.PrintOutcome(true, sum.label, sum.detail)
	}
	return 0, buf.String(), nil
}
