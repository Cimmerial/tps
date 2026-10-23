package ui

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
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
	return ReadEditedLineFuncHistoryWithOpts(promptFn, history, DefaultReadLineOptions())
}

// ReadEditedLineFuncHistoryWithOpts is ReadEditedLineFuncHistory with box chrome and menu cancel behavior.
func ReadEditedLineFuncHistoryWithOpts(promptFn func() string, history *[]string, opts ReadLineOptions) (string, error) {
	if !stdinIsTTY() {
		return readLineFallback(promptFn())
	}
	fd := int(os.Stdin.Fd())
	old, err := term.MakeRaw(fd)
	if err != nil {
		return readLineFallback(promptFn())
	}
	defer func() {
		fmt.Print("\033[?7h\033[?2004l")
		term.Restore(fd, old)
	}()

	fmt.Print("\033[?2004h\033[?7l")

	buf := []rune{}
	pos := 0
	histIdx := -1
	var histDraft string
	pasting := false
	boxDrawn := false
	boxMidRows := 0
	inputCursorLine := 0
	in := bufio.NewReader(os.Stdin)

	redrawPlain := func() {
		fmt.Print("\r\033[K" + promptFn() + string(buf))
		tail := len(buf) - pos
		if tail > 0 {
			fmt.Printf("\033[%dD", tail)
		}
	}
	inputCursorFromBefore := func(before string, inner int) (lineIdx, colInLine int) {
		segs := WrapBoxContent(before, inner)
		lineIdx = len(segs) - 1
		if lineIdx < 0 {
			lineIdx = 0
		}
		colInLine = DisplayWidth(segs[lineIdx])
		return lineIdx, colInLine
	}
	placeInputCursor := func(midRows, lineIdx, colInLine int) {
		up := midRows + 1 - lineIdx
		fmt.Printf("\033[%dA\r", up)
		fmt.Printf("\033[%dC", 1+colInLine)
		inputCursorLine = lineIdx
	}
	moveBelowInputBox := func() {
		if !boxDrawn {
			return
		}
		down := boxMidRows - inputCursorLine
		if down > 0 {
			fmt.Printf("\033[%dB", down)
		}
		fmt.Print("\033[1B\r")
	}
	eraseInputBoxLines := func(midRows int) {
		for i := 0; i < midRows+2; i++ {
			fmt.Print("\033[1A\r\033[K")
		}
	}
	printInputBox := func(top string, mids []string, bot string) {
		fmt.Println(top)
		for _, m := range mids {
			fmt.Println(m)
		}
		fmt.Println(bot)
	}
	redrawBoxed := func() {
		w := boxWidth()
		inner := boxInnerWidth(w)
		mode := promptChrome
		content := promptFn() + string(buf)
		mids := boxSideLinesFromContent(mode, w, content)
		top := boxTopLine(mode, w)
		bot := boxBottomLine(mode, w)
		before := promptFn() + string(buf[:pos])
		lineIdx, colInLine := inputCursorFromBefore(before, inner)

		if boxDrawn {
			moveBelowInputBox()
			eraseInputBoxLines(boxMidRows)
		}

		printInputBox(top, mids, bot)
		boxMidRows = len(mids)
		boxDrawn = true
		placeInputCursor(boxMidRows, lineIdx, colInLine)
	}
	redraw := func() {
		if opts.Boxed {
			redrawBoxed()
		} else {
			redrawPlain()
		}
	}
	redrawMaybe := func() {
		if !pasting {
			redraw()
		}
	}
	clearInputBox := func() {
		if opts.Boxed && boxDrawn {
			moveBelowInputBox()
			eraseInputBoxLines(boxMidRows)
			boxDrawn = false
			boxMidRows = 0
			inputCursorLine = 0
		} else {
			fmt.Print("\r\033[K")
		}
	}
	finishLine := func(line string) (string, error) {
		clearInputBox()
		fmt.Println()
		if opts.ResponseRuleOnSubmit && strings.TrimSpace(line) != "" {
			PrintResponseRule()
		}
		appendLineHistory(history, line)
		return line, nil
	}

	setBuf := func(s string) {
		buf = []rune(s)
		pos = len(buf)
		redrawMaybe()
	}

	clearHistBrowse := func() {
		histIdx = -1
		histDraft = ""
	}

	endNotices := BeginPromptNotices(func(lines ...string) {
		if opts.Boxed && boxDrawn {
			clearInputBox()
		} else {
			fmt.Print("\r\033[K")
		}
		for _, line := range lines {
			fmt.Println(line)
		}
		redraw()
	})
	defer endNotices()

	redraw()

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
			if b == '\r' && in.Buffered() > 0 {
				peek, err := in.Peek(1)
				if err == nil && len(peek) > 0 && peek[0] == '\n' {
					_, _ = in.ReadByte()
				}
			}
			if pasting {
				continue
			}
			return finishLine(string(buf))
		case 3: // Ctrl+C
			clearInputBox()
			fmt.Println()
			return "", fmt.Errorf("interrupted")
		case 4: // Ctrl+D
			if len(buf) == 0 {
				clearInputBox()
				fmt.Println()
				return "", io.EOF
			}
		case 23: // Ctrl+W
			clearHistBrowse()
			start := moveWordLeft(buf, pos)
			if start < pos {
				buf = append(buf[:start], buf[pos:]...)
				pos = start
				redrawMaybe()
			}
		case 127, 8: // backspace / DEL
			clearHistBrowse()
			if pos > 0 {
				buf = append(buf[:pos-1], buf[pos:]...)
				pos--
				redrawMaybe()
			}
		case 27: // ESC
			seq, bare, err := readEscapeSeqMaybeBare(in)
			if err != nil {
				return "", err
			}
			if bare {
				if opts.CancelOnEsc {
					return finishLine("")
				}
				continue
			}
			switch seq {
			case "[200~":
				pasting = true
			case "[201~":
				pasting = false
				redraw()
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
					redrawMaybe()
				}
			case "[C", "OC":
				clearHistBrowse()
				if pos < len(buf) {
					pos++
					redrawMaybe()
				}
			case "[1;3D", "[1;5D":
				clearHistBrowse()
				pos = moveWordLeft(buf, pos)
				redrawMaybe()
			case "[1;3C", "[1;5C":
				clearHistBrowse()
				pos = moveWordRight(buf, pos)
				redrawMaybe()
			case "[3~":
				clearHistBrowse()
				if pos < len(buf) {
					buf = append(buf[:pos], buf[pos+1:]...)
					redrawMaybe()
				}
			case "[3;3~":
				clearHistBrowse()
				end := moveWordRight(buf, pos)
				if end > pos {
					buf = append(buf[:pos], buf[end:]...)
					redrawMaybe()
				}
			case "[H", "[1~", "OH":
				clearHistBrowse()
				pos = 0
				redrawMaybe()
			case "[F", "[4~", "OF":
				clearHistBrowse()
				pos = len(buf)
				redrawMaybe()
			case "\x1bb":
				clearHistBrowse()
				pos = moveWordLeft(buf, pos)
				redrawMaybe()
			case "\x1bf":
				clearHistBrowse()
				pos = moveWordRight(buf, pos)
				redrawMaybe()
			case "\x1bd":
				clearHistBrowse()
				end := moveWordRight(buf, pos)
				if end > pos {
					buf = append(buf[:pos], buf[end:]...)
					redrawMaybe()
				}
			case "\x1b\x7f", "\x1b\b":
				clearHistBrowse()
				start := moveWordLeft(buf, pos)
				if start < pos {
					buf = append(buf[:start], buf[pos:]...)
					pos = start
					redrawMaybe()
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
			redrawMaybe()
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

func readEscapeSeqMaybeBare(in *bufio.Reader) (seq string, bare bool, err error) {
	if in.Buffered() > 0 {
		s, e := readEscapeSeq(in)
		return s, false, e
	}
	if err := os.Stdin.SetReadDeadline(time.Now().Add(30 * time.Millisecond)); err == nil {
		b, e := in.ReadByte()
		_ = os.Stdin.SetReadDeadline(time.Time{})
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Timeout() {
				return "", true, nil
			}
			return "", false, e
		}
		s, e := consumeEscapeFromByte(in, b)
		return s, false, e
	}
	s, e := readEscapeSeq(in)
	return s, false, e
}

func consumeEscapeFromByte(in *bufio.Reader, b byte) (string, error) {
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
