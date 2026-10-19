package cli

import (
	"github.com/Cimmerial/tps/internal/ui"
)

func readMenuLine() (string, error) {
	return ui.ReadEditedLine(ui.Prompt())
}
