package update

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Cimmerial/tps/internal/git"
)

// gitPullRunner runs git pull; tests may replace it.
var gitPullRunner = defaultGitPullRunner

func defaultGitPullRunner(cmd *exec.Cmd) error {
	return cmd.Run()
}

func gitPullCommand(root string) (*exec.Cmd, error) {
	exe, err := git.Executable()
	if err != nil {
		return nil, err
	}
	return exec.Command(exe, "-C", root, "pull", "--ff-only"), nil
}

// PullFastForward updates the linked tps repo from its upstream (current branch).
func PullFastForward(root string) error {
	cmd, err := gitPullCommand(root)
	if err != nil {
		return fmt.Errorf("git pull: %w", err)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := gitPullRunner(cmd); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return fmt.Errorf("git pull: %w: %s", err, msg)
		}
		return fmt.Errorf("git pull: %w", err)
	}
	return nil
}
