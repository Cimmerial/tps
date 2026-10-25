//go:build windows

package runner

import (
	"os"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func handoffScript(p *config.Project, script string) error {
	ui.ClearScreen()
	code, err := RunScript(p, script)
	if err != nil {
		return err
	}
	os.Exit(code)
	return nil
}
