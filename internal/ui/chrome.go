package ui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// ChromeMode selects info (magenta) vs choice (red) chrome.
type ChromeMode int

const (
	ChromeInfo    ChromeMode = 0
	ChromeChoice  ChromeMode = 1
	ChromeSession ChromeMode = 2 // home session menus (blue box)
)

var promptChrome = ChromeInfo

// SetPromptChrome sets the color mode for Prompt / PromptBusy until changed.
func SetPromptChrome(mode ChromeMode) {
	promptChrome = mode
}

// PromptChrome returns the active prompt chrome mode.
func PromptChrome() ChromeMode {
	return promptChrome
}

func chromeWire(mode ChromeMode) string {
	switch mode {
	case ChromeChoice:
		return red
	case ChromeSession:
		return blue
	default:
		return magenta
	}
}

func paintChrome(mode ChromeMode, s string) string {
	return paint(chromeWire(mode), s)
}

// PrintResponseRule prints a full-width magenta horizontal rule (response separator).
func PrintResponseRule() {
	w := TerminalWidth()
	if w <= 0 {
		w = FallbackTerminalWidth()
	}
	line := strings.Repeat("─", w)
	if ColorEnabled() {
		fmt.Println(paint(magenta, line))
	} else {
		fmt.Println(line)
	}
}

const (
	boxTopLeft     = '┌'
	boxTopRight    = '┐'
	boxBottomLeft  = '└'
	boxBottomRight = '┘'
	boxHorizontal  = '─'
	boxVertical    = '│'
	treeBranch     = '├'
	treeLast       = '└'
	treeVert       = '│'
	treeHoriz      = '─'
)

func boxWidth() int {
	w := TerminalWidth()
	if w <= 0 {
		w = FallbackTerminalWidth()
	}
	return w
}

func boxInnerWidth(outer int) int {
	if outer < 4 {
		return 1
	}
	return outer - 2
}

func boxTopLine(mode ChromeMode, outer int) string {
	inner := boxInnerWidth(outer)
	h := strings.Repeat(string(boxHorizontal), inner)
	return paintChrome(mode, string(boxTopLeft)+h+string(boxTopRight))
}

func boxBottomLine(mode ChromeMode, outer int) string {
	inner := boxInnerWidth(outer)
	h := strings.Repeat(string(boxHorizontal), inner)
	return paintChrome(mode, string(boxBottomLeft)+h+string(boxBottomRight))
}

func boxSideLine(mode ChromeMode, outer int, content string) string {
	inner := boxInnerWidth(outer)
	pad := inner - VisibleWidth(content)
	if pad < 0 {
		pad = 0
	}
	innerText := content + strings.Repeat(" ", pad)
	return paintChrome(mode, string(boxVertical)) + innerText + paintChrome(mode, string(boxVertical))
}

// wrapForBoxInner splits content into segments that each fit the box inner width.
func wrapForBoxInner(s string, inner int) []string {
	if inner < 1 {
		inner = 1
	}
	var out []string
	rest := s
	for rest != "" {
		piece, tail := splitVisibleMax(rest, inner)
		if piece == "" {
			break
		}
		out = append(out, piece)
		rest = tail
	}
	if len(out) == 0 {
		out = []string{""}
	}
	return out
}

func splitVisibleMax(s string, max int) (head, tail string) {
	if max < 1 || s == "" {
		return "", s
	}
	var b strings.Builder
	used := 0
	i := 0
	for i < len(s) {
		seg, dw, next := nextBoxContentSegment(s, i)
		if dw == 0 {
			b.WriteString(seg)
			i = next
			continue
		}
		if used+dw > max && used > 0 {
			break
		}
		if used+dw > max {
			break
		}
		b.WriteString(seg)
		used += dw
		i = next
	}
	if b.Len() == 0 && i < len(s) {
		seg, _, next := nextBoxContentSegment(s, i)
		b.WriteString(seg)
		i = next
	}
	return b.String(), s[i:]
}

// WrapBoxContent splits s into lines that fit the box inner width (ANSI-aware).
func WrapBoxContent(s string, inner int) []string {
	if inner < 1 {
		inner = 1
	}
	var lines []string
	var b strings.Builder
	lineW := 0
	for i := 0; i < len(s); {
		seg, dw, next := nextBoxContentSegment(s, i)
		i = next
		if dw == 0 {
			b.WriteString(seg)
			continue
		}
		if lineW+dw > inner && lineW > 0 {
			lines = append(lines, b.String())
			b.Reset()
			lineW = 0
		}
		b.WriteString(seg)
		lineW += dw
	}
	if b.Len() > 0 || len(lines) == 0 {
		lines = append(lines, b.String())
	}
	return lines
}

func nextBoxContentSegment(s string, i int) (seg string, displayWidth int, next int) {
	if i < len(s) && s[i] == '\033' {
		j := skipEscape(s, i)
		return s[i:j], 0, j
	}
	r, size := utf8.DecodeRuneInString(s[i:])
	if size == 0 {
		return "", 0, len(s)
	}
	dw := 1
	if r >= 0x2800 && r <= 0x28FF {
		dw = 2
	}
	return s[i : i+size], dw, i + size
}

func boxSideLinesFromContent(mode ChromeMode, outer int, content string) []string {
	inner := boxInnerWidth(outer)
	segs := wrapForBoxInner(content, inner)
	out := make([]string, len(segs))
	for i, seg := range segs {
		out[i] = boxSideLine(mode, outer, seg)
	}
	return out
}

// MenuBox renders a terminal-width box with tree-prefixed lines.
type MenuBox struct {
	mode   ChromeMode
	outer  int
	lines  []string
}

// NewMenuBox starts a menu box for the given chrome mode.
func NewMenuBox(mode ChromeMode) *MenuBox {
	return &MenuBox{mode: mode, outer: boxWidth()}
}

// Title adds a top-level title row (no tree prefix).
func (b *MenuBox) Title(title string) {
	b.lines = append(b.lines, title)
}

// Section adds a nested section label (e.g. GENERAL).
func (b *MenuBox) Section(label string, lastInBox bool) {
	prefix := string(treeBranch) + string(treeHoriz) + " "
	if lastInBox {
		prefix = string(treeLast) + string(treeHoriz) + " "
	}
	b.lines = append(b.lines, prefix+label)
}

// SectionTrail controls tree branch characters for nested rows.
type SectionTrail struct {
	LastSection bool
	LastRow     bool
}

// ChoiceRow is number | option | abbrev/status | detail (red | white | blue | gray).
func (b *MenuBox) ChoiceRow(trail SectionTrail, num, option, mid, detail string) {
	b.lines = append(b.lines, b.formatRow(trail, num, option, mid, detail, true))
}

// InfoRow is option | abbrev | detail (white | blue | gray) without a number column.
func (b *MenuBox) InfoRow(trail SectionTrail, option, abbrev, detail string) {
	b.lines = append(b.lines, b.formatRow(trail, "", option, abbrev, detail, false))
}

// Plain adds a line with tree continuation prefix only.
func (b *MenuBox) Plain(trail SectionTrail, text string) {
	b.lines = append(b.lines, b.treePrefix(trail)+text)
}

// ChoiceHelpRow adds a value/detail line under a numbered choice (│ └─ connector).
func (b *MenuBox) ChoiceHelpRow(parentTrail SectionTrail, text string) {
	prefix := b.treePrefix(parentTrail)
	connector := string(treeVert) + "  " + string(treeLast) + string(treeHoriz) + " "
	b.lines = append(b.lines, prefix+connector+text)
}

func (b *MenuBox) treePrefix(trail SectionTrail) string {
	if trail.LastSection {
		return "   "
	}
	return string(treeVert) + "  "
}

func (b *MenuBox) formatRow(trail SectionTrail, num, option, mid, detail string, numbered bool) string {
	branch := string(treeBranch) + string(treeHoriz) + " "
	if trail.LastRow {
		branch = string(treeLast) + string(treeHoriz) + " "
	}
	prefix := b.treePrefix(trail) + branch

	var body string
	if numbered {
		body = formatMenuColumns(num, option, mid, detail, true)
	} else {
		body = formatMenuColumns("", option, mid, detail, false)
	}
	return prefix + body
}

func paintMidColumn(mid string) string {
	if strings.Contains(mid, "\033") {
		return mid
	}
	return BrightBlue(mid)
}

func paintOptionColumn(opt string) string {
	if strings.Contains(opt, "\033") {
		return opt
	}
	return Bold(opt)
}

func formatMenuColumns(num, option, mid, detail string, numbered bool) string {
	sep := Dim(" | ")
	if numbered && strings.TrimSpace(num) != "" {
		numCol := Bold(Red(num))
		if strings.TrimSpace(mid) != "" {
			midCol := paintMidColumn(mid)
			if strings.TrimSpace(detail) != "" {
				return numCol + sep + paintOptionColumn(option) + sep + midCol + sep + Gray(detail)
			}
			return numCol + sep + paintOptionColumn(option) + sep + midCol
		}
		if strings.TrimSpace(detail) != "" {
			return numCol + sep + paintOptionColumn(option) + sep + Gray(detail)
		}
		return numCol + sep + paintOptionColumn(option)
	}
	optCol := paintOptionColumn(option)
	if strings.TrimSpace(mid) != "" {
		midCol := paintMidColumn(mid)
		if strings.TrimSpace(detail) != "" {
			return optCol + sep + midCol + sep + Gray(detail)
		}
		return optCol + sep + midCol
	}
	if strings.TrimSpace(detail) != "" {
		return optCol + sep + Gray(detail)
	}
	return optCol
}

// Flush prints the box to stdout.
func (b *MenuBox) Flush() {
	if len(b.lines) == 0 {
		return
	}
	fmt.Println(boxTopLine(b.mode, b.outer))
	inner := boxInnerWidth(b.outer)
	for _, line := range b.lines {
		for _, seg := range wrapForBoxInner(line, inner) {
			fmt.Println(boxSideLine(b.mode, b.outer, seg))
		}
	}
	fmt.Println(boxBottomLine(b.mode, b.outer))
}
