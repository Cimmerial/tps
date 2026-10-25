//go:build !windows

package runner

import (
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func handoffScript(p *config.Project, script string) error {
	script = strings.TrimSpace(script)
	if script == "" {
		return nil
	}
	cmd := scriptCommand(script)
	if p != nil && strings.TrimSpace(p.Path) != "" {
		if err := os.Chdir(p.Path); err != nil {
			return err
		}
	}
	env := scriptEnviron(p)
	bin, err := exec.LookPath(cmd.Args[0])
	if err != nil {
		return err
	}
	cmd.Args[0] = bin
	ui.ClearScreen()
	return syscall.Exec(bin, cmd.Args, env)
}
