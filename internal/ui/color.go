package ui

import "strings"

const (
	reset        = "\033[0m"
	bold         = "\033[1m"
	dim          = "\033[2m"
	cyan         = "\033[36m"
	green        = "\033[32m"
	yellow       = "\033[33m"
	magenta      = "\033[35m"
	blue         = "\033[34m"
	red          = "\033[31m"
	boldBrightFG = "\033[1;93m"
)

func ColorEnabled() bool {
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

// Abbrev highlights shortcut keys in menus and hints.
func Abbrev(s string) string { return paint(boldBrightFG, s) }

func Prompt() string {
	return PromptBusy(0, 0)
}

func PromptBusy(running int, tick int) string {
	if running <= 0 {
		return Cyan("> ")
	}
	frames := []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}
	frame := frames[tick%len(frames)]
	return Dim(frame+" ") + Cyan("> ")
}

func Logo(s string) string {
	if !ColorEnabled() {
		return s
	}
	return Cyan(s)
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
