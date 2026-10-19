package ui

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func ClearScreen() {
	if !stdoutIsTTY() {
		return
	}
	fmt.Print("\033[2J\033[H")
}

// ClearTerminal erases the visible screen and scrollback (VS Code / iTerm Cmd+K style).
func ClearTerminal() {
	if !stdoutIsTTY() {
		return
	}
	fmt.Print("\033[2J\033[3J\033[H")
}

func stdoutIsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func stdinIsTTY() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}
