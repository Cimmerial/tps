package cli

import (
	"bufio"
	"fmt"

	"github.com/Cimmerial/tps/internal/ui"
)

func readLine(_ *bufio.Reader, running int, tick int, padBefore bool) (string, error) {
	if padBefore {
		fmt.Print("\r\033[K\n")
	}
	line, err := ui.ReadEditedLine(ui.PromptBusy(running, tick))
	if err != nil {
		return "", err
	}
	return line, nil
}

func needsGapBeforeNextPrompt(action InputAction) bool {
	switch action {
	case ActionNoop, ActionQuit, ActionClear, ActionSwitchProject, ActionOptions:
		return false
	default:
		return true
	}
}
