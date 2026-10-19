package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Cimmerial/tps/internal/ui"
)

func Run(dir string, args []string) (int, error) {
	if len(args) == 0 {
		return 1, fmt.Errorf("no git command given")
	}
	if _, err := exec.LookPath("git"); err != nil {
		return 1, fmt.Errorf("git not found in PATH")
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			return exit.ExitCode(), nil
		}
		return 1, err
	}
	return 0, nil
}

// ShipOptions configures git ship; nil fields use defaults.
type ShipOptions struct {
	OnSuspicious func([]string) (SuspiciousAction, error)
}

func Ship(dir, message string) (int, error) {
	return ShipWithOptions(dir, message, nil)
}

func ShipWithOptions(dir, message string, opts *ShipOptions) (int, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return 1, fmt.Errorf("commit message required (example: g ship fix the menu)")
	}
	code, err := Run(dir, []string{"add", "-A"})
	if err != nil {
		return code, err
	}
	if code != 0 {
		return code, nil
	}

	if stagedEmpty(dir) {
		fmt.Println(ui.Yellow("nothing to commit"))
		return 0, nil
	}

	onSuspicious := PromptSuspiciousFiles
	if opts != nil && opts.OnSuspicious != nil {
		onSuspicious = opts.OnSuspicious
	}
	if suspicious := SuspiciousPaths(listStagedPaths(dir)); len(suspicious) > 0 {
		action, err := onSuspicious(suspicious)
		if err != nil {
			return 1, err
		}
		switch action {
		case SuspiciousCancel:
			fmt.Println(ui.Dim("ship cancelled"))
			return 0, nil
		case SuspiciousExcludeGitignore:
			if err := unstagePaths(dir, suspicious); err != nil {
				return 1, err
			}
			if err := appendGitignoreLines(dir, suspicious); err != nil {
				return 1, fmt.Errorf(".gitignore: %w", err)
			}
			code, err = Run(dir, []string{"add", ".gitignore"})
			if err != nil {
				return code, err
			}
			if code != 0 {
				return code, nil
			}
			if stagedEmpty(dir) {
				fmt.Println(ui.Yellow("nothing to commit after excluding suspicious files"))
				return 0, nil
			}
		case SuspiciousProceed:
			// continue
		default:
			return 1, fmt.Errorf("invalid suspicious-file choice")
		}
	}

	code, err = Run(dir, []string{"commit", "-m", message})
	if err != nil {
		return code, err
	}
	if code != 0 {
		return code, nil
	}

	code, err = Run(dir, []string{"push"})
	if err != nil {
		return code, err
	}
	if code == 0 {
		printShipOutcome(dir, message)
	}
	return code, nil
}

func printShipOutcome(dir, message string) {
	fmt.Println()
	fmt.Printf("%s %s\n", ui.Green("shipped —"), ui.Bold(message))
	if line := gitOutput(dir, "log", "-1", "--oneline"); line != "" {
		fmt.Println(ui.Cyan(line))
	}
	if line := gitOutput(dir, "status", "-sb"); line != "" {
		fmt.Println(ui.Dim(line))
	}
	fmt.Println()
}

func gitOutput(dir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func stagedEmpty(dir string) bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	cmd.Dir = dir
	return cmd.Run() == nil
}

func listStagedPaths(dir string) []string {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	var paths []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			paths = append(paths, line)
		}
	}
	return paths
}

func unstagePaths(dir string, paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	args := []string{"restore", "--staged", "--"}
	args = append(args, paths...)
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func ParseArgsLine(line string) []string {
	return strings.Fields(strings.TrimSpace(line))
}
