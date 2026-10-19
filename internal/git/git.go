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

func Ship(dir, message string) (int, error) {
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

func ParseArgsLine(line string) []string {
	return strings.Fields(strings.TrimSpace(line))
}
