package git

import (
	"errors"
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
	cmd := exec.Command("git", quietGitArgs(args)...)
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
	// OnSelect chooses paths to stage; nil means stage everything (git add -A).
	OnSelect func(paths []string) ([]string, error)
	// CommitMessage, when set, is read after OnSelect returns (interactive ship message).
	CommitMessage *string
}

func Ship(dir, message string) (int, error) {
	return ShipWithOptions(dir, message, nil)
}

func ShipWithOptions(dir, message string, opts *ShipOptions) (int, error) {
	interactive := opts != nil && opts.OnSelect != nil
	message = strings.TrimSpace(message)
	if message == "" && !interactive {
		return 1, fmt.Errorf("commit message required (example: g ship fix the menu)")
	}

	var (
		allPaths []string
		code     int
		err      error
	)
	if interactive {
		allPaths, err = ListChangedPaths(dir)
		if err != nil {
			return 1, err
		}
		if len(allPaths) == 0 {
			fmt.Println(ui.Yellow("nothing to commit"))
			return 0, nil
		}
		include, err := opts.OnSelect(allPaths)
		if err != nil {
			if errors.Is(err, ErrShipSelectCancel) {
				fmt.Println(ui.Dim("ship cancelled"))
				return 0, nil
			}
			return 1, err
		}
		if len(include) == 0 {
			fmt.Println(ui.Yellow("nothing to commit"))
			return 0, nil
		}
		if err := StageShipFiles(dir, allPaths, include); err != nil {
			return 1, err
		}
		if opts.CommitMessage != nil {
			message = strings.TrimSpace(*opts.CommitMessage)
		}
	} else {
		code, err = Run(dir, []string{"add", "-A"})
		if err != nil {
			return code, err
		}
		if code != 0 {
			return code, nil
		}
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

	if message == "" {
		return 1, fmt.Errorf("commit message required (set commit name before shipping)")
	}

	code, err = Run(dir, []string{"commit", "-m", message})
	if err != nil {
		return code, err
	}
	if code != 0 {
		printShipStepOutcome(dir, false, "commit", shipStepDetail(dir, message, code))
		return code, nil
	}
	printShipStepOutcome(dir, true, "commit", shipStepDetail(dir, message, 0))

	code, err = Run(dir, []string{"push"})
	if err != nil {
		return code, err
	}
	if code != 0 {
		printShipStepOutcome(dir, false, "push", shipPushDetail(dir, code))
		return code, nil
	}
	printShipStepOutcome(dir, true, "push", shipPushDetail(dir, 0))
	printShipOutcome(dir, message)
	return code, nil
}

func printShipStepOutcome(dir string, ok bool, step, detail string) {
	if !ui.SimplifiedOutput() {
		return
	}
	ui.PrintOutcome(ok, step, detail)
}

func shipStepDetail(dir, message string, exitCode int) string {
	if exitCode != 0 {
		return fmt.Sprintf("exit %d", exitCode)
	}
	if line := gitOutput(dir, "log", "-1", "--oneline"); line != "" {
		return line
	}
	return message
}

func shipPushDetail(dir string, exitCode int) string {
	if exitCode != 0 {
		return fmt.Sprintf("exit %d", exitCode)
	}
	if line := gitOutput(dir, "status", "-sb"); line != "" {
		return line
	}
	if branch := gitOutput(dir, "rev-parse", "--abbrev-ref", "HEAD"); branch != "" {
		return branch
	}
	return "pushed"
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
	if ui.SimplifiedOutput() {
		ui.PrintOutcome(true, "ship", message)
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
