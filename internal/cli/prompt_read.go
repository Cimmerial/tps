package cli

import (
	"bufio"
	"fmt"

	"github.com/Cimmerial/tps/internal/jobs"
	"github.com/Cimmerial/tps/internal/ui"
)

var promptLineHistory []string

func readLine(_ *bufio.Reader, projectID string, tick int, padBefore bool) (string, error) {
	drainDetachedCompleteNotices()
	if padBefore {
		fmt.Print("\r\033[K\n")
	}
	line, err := ui.ReadEditedLineFuncHistory(func() string {
		active, _ := jobs.Reap(projectID)
		return ui.PromptBusy(len(active), tick)
	}, &promptLineHistory)
	if err != nil {
		return "", err
	}
	return line, nil
}

func needsGapBeforeNextPrompt(action InputAction) bool {
	switch action {
	case ActionNoop, ActionQuit, ActionClear, ActionNewCommand, ActionSwitchProject, ActionSwitchProjectBack, ActionOptions, ActionVersion:
		return false
	default:
		return true
	}
}
