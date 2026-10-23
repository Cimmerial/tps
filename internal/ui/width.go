package ui

import (
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

const defaultTerminalColumns = 80

// VisibleWidth is the display width of s, ignoring ANSI escape sequences.
func VisibleWidth(s string) int {
	n := 0
	i := 0
	for i < len(s) {
		if s[i] == '\033' {
			i = skipEscape(s, i)
			continue
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		if size == 0 {
			break
		}
		n++
		i += size
	}
	return n
}

// DisplayWidth is VisibleWidth plus one column when a braille spinner glyph is present
// (many terminals render those at width 2 while the rest of the prompt is width 1).
func DisplayWidth(s string) int {
	w := VisibleWidth(s)
	if brailleSpinnerIn(s) {
		w++
	}
	return w
}

func brailleSpinnerIn(s string) bool {
	for i := 0; i < len(s); {
		if s[i] == '\033' {
			i = skipEscape(s, i)
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if size == 0 {
			break
		}
		if r >= 0x2800 && r <= 0x28FF {
			return true
		}
		i += size
	}
	return false
}

func skipEscape(s string, i int) int {
	if i >= len(s) || s[i] != '\033' {
		return i + 1
	}
	i++
	if i < len(s) && s[i] == '[' {
		i++
		for i < len(s) && (s[i] < 0x40 || s[i] > 0x7e) {
			i++
		}
		if i < len(s) {
			i++
		}
		return i
	}
	if i < len(s) && s[i] == ']' {
		i++
		for i < len(s) && s[i] != '\033' && s[i] != '\007' {
			i++
		}
		if i < len(s) && s[i] == '\033' {
			if i+1 < len(s) && s[i+1] == '\\' {
				return i + 2
			}
		}
		if i < len(s) {
			i++
		}
		return i
	}
	if i < len(s) {
		i++
	}
	return i
}

// TerminalWidth returns stdout width in columns, or 0 if unknown.
func TerminalWidth() int {
	for _, fd := range terminalFDs() {
		w, _, err := term.GetSize(fd)
		if err == nil && w >= 20 {
			return w
		}
	}
	if w := columnsFromEnv(); w >= 20 {
		return w
	}
	return 0
}

// FallbackTerminalWidth is used when TerminalWidth is unknown (padding estimates).
func FallbackTerminalWidth() int {
	if w := columnsFromEnv(); w >= 20 {
		return w
	}
	return defaultTerminalColumns
}

func columnsFromEnv() int {
	raw := strings.TrimSpace(os.Getenv("COLUMNS"))
	if raw == "" {
		return 0
	}
	w, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return w
}

func terminalFDs() []int {
	fds := make([]int, 0, 4)
	for _, f := range []*os.File{os.Stdout, os.Stderr, os.Stdin} {
		if f != nil {
			fds = append(fds, int(f.Fd()))
		}
	}
	if f, err := os.Open("/dev/tty"); err == nil {
		fds = append(fds, int(f.Fd()))
		_ = f.Close()
	}
	return fds
}

// AlignLeftRight prints left and right on one line with spaces so right ends at width columns.
// If width is 0, uses FallbackTerminalWidth. If the line still does not fit, returns two lines.
func AlignLeftRight(left, right string, width int) string {
	if right == "" {
		return left
	}
	if left == "" {
		return right
	}
	knownWidth := width > 0
	if !knownWidth {
		width = FallbackTerminalWidth()
	}
	lw := VisibleWidth(left)
	rw := VisibleWidth(right)
	pad := width - lw - rw
	if pad >= 1 {
		return left + strings.Repeat(" ", pad) + right
	}
	if knownWidth {
		if pad < 1 {
			pad = 1
		}
		return left + strings.Repeat(" ", pad) + right
	}
	pad2 := width - rw
	if pad2 < 0 {
		pad2 = 0
	}
	return left + "\n" + strings.Repeat(" ", pad2) + right
}
