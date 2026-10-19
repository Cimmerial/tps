package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/term"
)

// ReadEditedLine prints prompt and reads a line with basic editing (arrows, backspace, delete).
func ReadEditedLine(prompt string) (string, error) {
	return ReadEditedLineFunc(func() string { return prompt })
}

// ReadEditedLineFunc is like ReadEditedLine but re-evaluates the prompt on each redraw (e.g. terminal count).
func ReadEditedLineFunc(promptFn func() string) (string, error) {
	return ReadEditedLineFuncHistory(promptFn, nil)
}

// ReadEditedLineFuncHistory is like ReadEditedLineFunc but supports up/down command history when history is non-nil.
func ReadEditedLineFuncHistory(promptFn func() string, history *[]string) (string, error) {
	if !stdinIsTTY() {
		return readLineFallback(promptFn())
	}
	fd := int(os.Stdin.Fd())
	old, err := term.MakeRaw(fd)
	if err != nil {
		return readLineFallback(promptFn())
	}
	defer term.Restore(fd, old)

	fmt.Print("\r\033[K" + promptFn())
	buf := []rune{}
	pos := 0
	histIdx := -1
	var histDraft string
	in := bufio.NewReader(os.Stdin)

	redraw := func() {
		fmt.Print("\r\033[K" + promptFn() + string(buf))
		tail := len(buf) - pos
		if tail > 0 {
			fmt.Printf("\033[%dD", tail)
		}
	}

	setBuf := func(s string) {
		buf = []rune(s)
		pos = len(buf)
		redraw()
	}

	clearHistBrowse := func() {
		histIdx = -1
		histDraft = ""
	}

	endNotices := BeginPromptNotices(func(lines ...string) {
		fmt.Print("\r\033[K")
		for _, line := range lines {
			fmt.Println(line)
		}
		redraw()
	})
	defer endNotices()

	for {
		b, err := in.ReadByte()
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return "", err
			}
			return "", err
		}

		switch b {
		case '\r', '\n':
			fmt.Print("\r\033[K\n")
			line := string(buf)
			appendLineHistory(history, line)
			return line, nil
		case 3: // Ctrl+C
			fmt.Print("\r\033[K\n")
			return "", fmt.Errorf("interrupted")
		case 4: // Ctrl+D
			if len(buf) == 0 {
				fmt.Print("\r\033[K\n")
				return "", io.EOF
			}
		case 23: // Ctrl+W
			clearHistBrowse()
			start := moveWordLeft(buf, pos)
			if start < pos {
				buf = append(buf[:start], buf[pos:]...)
				pos = start
				redraw()
			}
		case 127, 8: // backspace / DEL
			clearHistBrowse()
			if pos > 0 {
				buf = append(buf[:pos-1], buf[pos:]...)
				pos--
				redraw()
			}
		case 27: // ESC
			seq, err := readEscapeSeq(in)
			if err != nil {
				return "", err
			}
			switch seq {
			case "[A", "OA":
				if history != nil && len(*history) > 0 {
					if histIdx == -1 {
						histDraft = string(buf)
						histIdx = len(*history) - 1
					} else if histIdx > 0 {
						histIdx--
					}
					setBuf((*history)[histIdx])
				}
			case "[B", "OB":
				if history != nil && histIdx != -1 {
					if histIdx < len(*history)-1 {
						histIdx++
						setBuf((*history)[histIdx])
					} else {
						clearHistBrowse()
						setBuf(histDraft)
					}
				}
			case "[D", "OD":
				clearHistBrowse()
				if pos > 0 {
					pos--
					redraw()
				}
			case "[C", "OC":
				clearHistBrowse()
				if pos < len(buf) {
					pos++
					redraw()
				}
			case "[1;3D", "[1;5D":
				clearHistBrowse()
				pos = moveWordLeft(buf, pos)
				redraw()
			case "[1;3C", "[1;5C":
				clearHistBrowse()
				pos = moveWordRight(buf, pos)
				redraw()
			case "[3~":
				clearHistBrowse()
				if pos < len(buf) {
					buf = append(buf[:pos], buf[pos+1:]...)
					redraw()
				}
			case "[3;3~":
				clearHistBrowse()
				end := moveWordRight(buf, pos)
				if end > pos {
					buf = append(buf[:pos], buf[end:]...)
					redraw()
				}
			case "[H", "[1~", "OH":
				clearHistBrowse()
				pos = 0
				redraw()
			case "[F", "[4~", "OF":
				clearHistBrowse()
				pos = len(buf)
				redraw()
			case "\x1bb":
				clearHistBrowse()
				pos = moveWordLeft(buf, pos)
				redraw()
			case "\x1bf":
				clearHistBrowse()
				pos = moveWordRight(buf, pos)
				redraw()
			case "\x1bd":
				clearHistBrowse()
				end := moveWordRight(buf, pos)
				if end > pos {
					buf = append(buf[:pos], buf[end:]...)
					redraw()
				}
			case "\x1b\x7f", "\x1b\b":
				clearHistBrowse()
				start := moveWordLeft(buf, pos)
				if start < pos {
					buf = append(buf[:start], buf[pos:]...)
					pos = start
					redraw()
				}
			}
		default:
			if b < 32 {
				continue
			}
			clearHistBrowse()
			r, err := readUTF8Rune(in, b)
			if err != nil {
				return "", err
			}
			if r == utf8.RuneError {
				continue
			}
			buf = append(buf[:pos], append([]rune{r}, buf[pos:]...)...)
			pos++
			redraw()
		}
	}
}

func appendLineHistory(history *[]string, line string) {
	if history == nil || line == "" {
		return
	}
	if len(*history) > 0 && (*history)[len(*history)-1] == line {
		return
	}
	*history = append(*history, line)
}

func isWordRune(r rune) bool {
	return r == '_' || r == '-' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func moveWordLeft(buf []rune, pos int) int {
	if pos == 0 {
		return 0
	}
	i := pos - 1
	for i > 0 && !isWordRune(buf[i]) {
		i--
	}
	for i > 0 && isWordRune(buf[i-1]) {
		i--
	}
	return i
}

func moveWordRight(buf []rune, pos int) int {
	if pos >= len(buf) {
		return len(buf)
	}
	i := pos
	for i < len(buf) && !isWordRune(buf[i]) {
		i++
	}
	for i < len(buf) && isWordRune(buf[i]) {
		i++
	}
	return i
}

func readLineFallback(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

func readUTF8Rune(in *bufio.Reader, first byte) (rune, error) {
	buf := []byte{first}
	for {
		if utf8.FullRune(buf) {
			r, _ := utf8.DecodeRune(buf)
			return r, nil
		}
		more, err := in.ReadByte()
		if err != nil {
			return utf8.RuneError, err
		}
		buf = append(buf, more)
		if len(buf) > utf8.UTFMax {
			return utf8.RuneError, nil
		}
	}
}

func readEscapeSeq(in *bufio.Reader) (string, error) {
	b, err := in.ReadByte()
	if err != nil {
		return "", err
	}
	if b == 'O' {
		c, err := in.ReadByte()
		if err != nil {
			return "", err
		}
		return "O" + string([]byte{c}), nil
	}
	if b != '[' {
		return string([]byte{27, b}), nil
	}
	var seq []byte
	seq = append(seq, '[')
	for {
		c, err := in.ReadByte()
		if err != nil {
			return "", err
		}
		seq = append(seq, c)
		if c == '~' {
			return string(seq), nil
		}
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			return string(seq), nil
		}
		if len(seq) > 16 {
			return string(seq), nil
		}
	}
}
