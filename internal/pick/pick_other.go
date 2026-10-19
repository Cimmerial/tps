//go:build !darwin

package pick

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/ui"
)

func ChooseFile(prompt string) (string, error) {
	return readPath(prompt)
}

func ChooseFiles(prompt string) ([]string, error) {
	line, err := readPath(prompt + " (space-separated paths)")
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, ErrCancelled
	}
	return fields, nil
}

func ChooseFolder(prompt string) (string, error) {
	return readPath(prompt)
}

func readPath(prompt string) (string, error) {
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
