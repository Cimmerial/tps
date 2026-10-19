package cli

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/jobs"
	"github.com/Cimmerial/tps/internal/runner"
	"github.com/Cimmerial/tps/internal/ui"
)

// executeProjectCommand runs precursors (in order) then the requested alias.
// Foreground steps share one shell so export/source in a precursor applies to later steps.
func executeProjectCommand(p *config.Project, alias string, inputs []string) (int, error) {
	chain, err := p.CommandRunChain(alias)
	if err != nil {
		return 1, err
	}
	target := chain[len(chain)-1]
	var fg []string
	for i, a := range chain {
		cmd := p.FindCommand(a)
		if cmd == nil {
			return 1, fmt.Errorf("unknown command %q", a)
		}
		var shellInputs []string
		if a == target {
			shellInputs = inputs
		}
		shells, err := cmd.ShellsWithInputs(shellInputs)
		if err != nil {
			return 1, err
		}
		if i < len(chain)-1 {
			fmt.Println(ui.Dim("precursor: " + a))
		}
		if runner.ShouldDetach(*cmd) {
			code, err := runForegroundChain(p, fg)
			fg = nil
			if err != nil {
				return code, err
			}
			if code != 0 {
				return code, nil
			}
			started, err := jobs.StartAllWithShells(p, a, shells)
			if err != nil {
				return 1, err
			}
			for _, job := range started {
				fmt.Println(ui.Green(fmt.Sprintf("started #%d in new terminal (%s)", job.ID, job.Alias)))
				watchCaptureIfNeeded(p, job)
			}
			continue
		}
		fg = append(fg, shells...)
	}
	return runForegroundChain(p, fg)
}

func runForegroundChain(p *config.Project, parts []string) (int, error) {
	if len(parts) == 0 {
		return 0, nil
	}
	script := strings.Join(parts, " && ")
	var code int
	var err error
	ui.BrailleWaitStderr("loading…", func() {
		code, err = runner.RunScript(p, script)
	})
	return code, err
}
