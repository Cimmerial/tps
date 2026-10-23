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
	red          = "\033[31m"
	brightWhite  = "\033[97m"
	brightBlue   = "\033[94m"
	gray         = "\033[90m"
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

// Faint styles menu previews of current values (lighter than body text, still readable).
func Faint(s string) string {
	return paint("\033[37m", s)
}
func Cyan(s string) string    { return paint(cyan, s) }
func Green(s string) string   { return paint(green, s) }
func Yellow(s string) string  { return paint(yellow, s) }
func Magenta(s string) string { return paint(magenta, s) }
func Blue(s string) string    { return paint(blue, s) }
func Red(s string) string { return paint(red, s) }
func BrightWhite(s string) string { return paint(brightWhite, s) }
func BrightBlue(s string) string { return paint(brightBlue, s) }
func Gray(s string) string { return paint(gray, s) }

// Abbrev highlights shortcut keys in menus and hints (bright blue).
func Abbrev(s string) string { return BrightBlue(s) }

func Prompt() string {
	return PromptBusy(0, 0)
}

// NewValuePrompt is the prefix when entering a new value in edit field prompts.
func NewValuePrompt() string {
	return Green("new ") + paint(red, ">: ")
}

func promptPrefix(mode ChromeMode) string {
	if mode == ChromeChoice {
		return paint(red, ">: ")
	}
	return paint(magenta, ">: ")
}

func BrailleFrame(tick int) string {
	if len(BrailleSpinnerFrames) == 0 {
		return "⣾"
	}
	return BrailleSpinnerFrames[tick%len(BrailleSpinnerFrames)]
}

func PromptBusy(running int, tick int) string {
	return PromptBusyChrome(promptChrome, running, tick)
}

// PromptBusyChrome builds the prompt prefix for a chrome mode.
func PromptBusyChrome(mode ChromeMode, running int, tick int) string {
	prompt := promptPrefix(mode)
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
