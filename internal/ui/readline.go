package ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// ReadEditedLine prints prompt and reads a line with basic editing (arrows, backspace, delete).
func ReadEditedLine(prompt string) (string, error) {
	if !stdinIsTTY() {
		return readLineFallback(prompt)
	}
	fd := int(os.Stdin.Fd())
	old, err := term.MakeRaw(fd)
	if err != nil {
		return readLineFallback(prompt)
	}
	defer term.Restore(fd, old)

	fmt.Print(prompt)
	buf := []rune{}
	pos := 0
	in := bufio.NewReader(os.Stdin)

	redraw := func() {
		fmt.Print("\r\033[K" + prompt + string(buf))
		tail := len(buf) - pos
		if tail > 0 {
			fmt.Printf("\033[%dD", tail)
		}
	}

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
			fmt.Println()
			return string(buf), nil
		case 3: // Ctrl+C
			fmt.Println()
			return "", fmt.Errorf("interrupted")
		case 4: // Ctrl+D
			if len(buf) == 0 {
				fmt.Println()
				return "", io.EOF
			}
		case 127, 8: // backspace / DEL
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
			case "[D", "OD":
				if pos > 0 {
					pos--
					redraw()
				}
			case "[C", "OC":
				if pos < len(buf) {
					pos++
					redraw()
				}
			case "[3~":
				if pos < len(buf) {
					buf = append(buf[:pos], buf[pos+1:]...)
					redraw()
				}
			case "[H", "[1~", "OH":
				pos = 0
				redraw()
			case "[F", "[4~", "OF":
				pos = len(buf)
				redraw()
			}
		default:
			if b < 32 {
				continue
			}
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
