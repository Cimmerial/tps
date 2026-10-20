package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func runPushRelease(p *config.Project, inputs []string) (int, error) {
	if p == nil {
		return 1, fmt.Errorf("push-release: no project")
	}
	script := filepath.Join(p.Path, "scripts", "push-release.sh")
	if _, err := os.Stat(script); err != nil {
		return 1, fmt.Errorf("push-release: missing %s", script)
	}
	args := []string{script}
	for _, raw := range inputs {
		in := strings.TrimSpace(raw)
		if in == "" {
			continue
		}
		switch strings.ToLower(in) {
		case "y", "yes":
			args = append(args, "-y")
		case "n", "no":
			continue
		default:
			v := in
			if !strings.HasPrefix(v, "v") {
				v = "v" + v
			}
			args = append(args, v)
		}
	}
	cmd := exec.Command("bash", args...)
	cmd.Dir = p.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if ui.SimplifiedOutput() {
		cmd.Env = append(os.Environ(), "TPS_SIMPLIFIED_OUTPUT=1")
		if ui.ColorEnabled() {
			cmd.Env = append(cmd.Env, "TPSU_COLOR=1")
		}
	}
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return ee.ExitCode(), nil
		}
		return 1, err
	}
	return 0, nil
}
