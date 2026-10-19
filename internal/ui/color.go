package ui

import (
	"fmt"
	"strings"
)

const (
	reset        = "\033[0m"
	bold         = "\033[1m"
	dim          = "\033[2m"
	cyan         = "\033[36m"
	green        = "\033[32m"
	yellow       = "\033[33m"
	magenta      = "\033[35m"
	blue         = "\033[34m"
	red = "\033[31m"
)

func ColorEnabled() bool {
	if current.ColorEnabled != nil {
		return *current.ColorEnabled
	}
	return stdoutIsTTY()
}

func paint(code, s string) string {
	if !ColorEnabled() {
		return s
	}
	return code + s + reset
}

func Bold(s string) string    { return paint(bold, s) }
func Dim(s string) string     { return paint(dim, s) }
func Cyan(s string) string    { return paint(cyan, s) }
func Green(s string) string   { return paint(green, s) }
func Yellow(s string) string  { return paint(yellow, s) }
func Magenta(s string) string { return paint(magenta, s) }
func Blue(s string) string    { return paint(blue, s) }
func Red(s string) string { return paint(red, s) }

// Abbrev highlights shortcut keys in menus and hints (bright accent).
func Abbrev(s string) string { return paint(abbrevCode(), s) }

func Prompt() string {
	return PromptBusy(0, 0)
}

var braillePromptFrames = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

func BrailleFrame(tick int) string {
	if len(braillePromptFrames) == 0 {
		return "⣾"
	}
	return braillePromptFrames[tick%len(braillePromptFrames)]
}

func PromptBusy(running int, tick int) string {
	prompt := paint(accentCode(), "> ")
	if running <= 0 {
		return prompt
	}
	return Dim(BrailleFrame(tick)+" ") + prompt
}

func TerminalsRunningLine(count int, tick int) string {
	if count <= 0 {
		return ""
	}
	return Dim(BrailleFrame(tick) + " " + fmt.Sprintf("%d terminal(s) running", count))
}

func Logo(s string) string {
	if !ColorEnabled() {
		return s
	}
	return paint(titleCardColorCode(), s)
}

func ColorizeLines(s, lineColor string) string {
	if !ColorEnabled() {
		return s
	}
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		lines[i] = paint(lineColor, line)
	}
	return strings.Join(lines, "\n")
}
