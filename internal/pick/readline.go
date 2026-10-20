package pick

import (
	"strings"

	"github.com/Cimmerial/tps/internal/ui"
)

func readPathLine(prompt string) (string, error) {
	line, err := ui.ReadEditedLine(ui.Dim(prompt) + ": ")
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return "", ErrCancelled
	}
	return line, nil
}

func readPathsLine(prompt string) ([]string, error) {
	line, err := readPathLine(prompt + " (space-separated paths)")
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, ErrCancelled
	}
	return fields, nil
}
