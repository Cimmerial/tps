package update

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// gitPullRunner runs git pull; tests may replace it.
var gitPullRunner = defaultGitPullRunner

func defaultGitPullRunner(cmd *exec.Cmd) error {
	return cmd.Run()
}

func gitPullCommand(root string) *exec.Cmd {
	return exec.Command("git", "-C", root, "pull", "--ff-only")
}

// PullFastForward updates the linked tps repo from its upstream (current branch).
func PullFastForward(root string) error {
	cmd := gitPullCommand(root)
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
