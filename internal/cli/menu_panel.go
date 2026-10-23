package cli

import (
	"strings"

	"github.com/Cimmerial/tps/internal/ui"
)

type menuPanelSection struct {
	Label string
	Rows  []numberedMenuRow
}

func stripANSI(s string) string {
	var out strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '\033' {
			i++
			if i < len(s) && s[i] == '[' {
				i++
				for i < len(s) && (s[i] < 0x40 || s[i] > 0x7e) {
					i++
				}
				if i < len(s) {
					i++
				}
				continue
			}
			if i < len(s) {
				i++
			}
			continue
		}
		out.WriteByte(s[i])
		i++
	}
	return out.String()
}

func choiceMenuOption(desc string) string {
	if strings.Contains(desc, "\033") {
		return desc
	}
	return stripANSI(desc)
}

func choiceMenuMaxNumWidth(sections []menuPanelSection, footer []numberedMenuRow) int {
	max := 0
	for _, sec := range sections {
		for _, r := range sec.Rows {
			if n := len(strings.TrimSpace(r.Num)); n > max {
				max = n
			}
		}
	}
	for _, r := range footer {
		if n := len(strings.TrimSpace(r.Num)); n > max {
			max = n
		}
	}
	return max
}

func printChoiceMenuPanel(title string, hints []string, sections []menuPanelSection, footer []numberedMenuRow) {
	box := ui.NewMenuBox(ui.ChromeChoice)
	if strings.TrimSpace(title) != "" {
		box.Title(ui.HeaderChoice(title))
	}
	maxNum := choiceMenuMaxNumWidth(sections, footer)
	for _, h := range hints {
		box.Plain(ui.SectionTrail{LastSection: true, LastRow: false}, ui.Gray(stripANSI(h)))
	}
	secCount := len(sections)
	for si, sec := range sections {
		lastSec := si == secCount-1
		if strings.TrimSpace(sec.Label) != "" {
			box.Section(ui.SubHeaderChoice(stripANSI(sec.Label)), lastSec && len(sec.Rows) == 0 && len(footer) == 0)
		}
		rows := sec.Rows
		for ri, r := range rows {
			trail := ui.SectionTrail{
				LastSection: lastSec,
				LastRow:     ri == len(rows)-1 && len(footer) == 0,
			}
			opt := choiceMenuOption(r.Desc)
			mid := r.Status
			if strings.TrimSpace(stripANSI(mid)) == "" {
				mid = stripANSI(r.Abbrev)
			}
			num := padPlain(strings.TrimSpace(r.Num), maxNum)
			box.ChoiceRow(trail, num, opt, mid, "")
			if strings.TrimSpace(stripANSI(r.HelpLine)) != "" {
				box.ChoiceHelpRow(trail, choiceMenuHelpLine(r.HelpLine))
			}
		}
	}
	for fi, r := range footer {
		trail := ui.SectionTrail{LastSection: true, LastRow: fi == len(footer)-1}
		opt := choiceMenuOption(r.Desc)
		mid := r.Status
		if strings.TrimSpace(stripANSI(mid)) == "" {
			mid = stripANSI(r.Abbrev)
		}
		num := padPlain(strings.TrimSpace(r.Num), maxNum)
		box.ChoiceRow(trail, num, opt, mid, "")
	}
	box.Flush()
}

// printChoiceEditMenuPanel is a numbered choice list in a red edit box (title uses HeaderChoice).
func printChoiceEditMenuPanel(title string, sections []menuPanelSection, footer []numberedMenuRow) {
	printNumberedChoicePanel(ui.ChromeChoice, title, ui.HeaderChoice, ui.SubHeaderChoice, choiceEditHelpLine, sections, footer)
}

func printNumberedChoicePanel(
	mode ui.ChromeMode,
	title string,
	titleFn func(string) string,
	sectionFn func(string) string,
	helpLineFn func(string) string,
	sections []menuPanelSection,
	footer []numberedMenuRow,
) {
	box := ui.NewMenuBox(mode)
	if strings.TrimSpace(title) != "" {
		box.Title(titleFn(title))
	}
	maxNum := choiceMenuMaxNumWidth(sections, footer)
	secCount := len(sections)
	for si, sec := range sections {
		lastSec := si == secCount-1
		if strings.TrimSpace(sec.Label) != "" {
			box.Section(sectionFn(stripANSI(sec.Label)), lastSec && len(sec.Rows) == 0 && len(footer) == 0)
		}
		rows := sec.Rows
		for ri, r := range rows {
			trail := ui.SectionTrail{
				LastSection: lastSec,
				LastRow:     ri == len(rows)-1 && len(footer) == 0,
			}
			opt := choiceMenuOption(r.Desc)
			mid := r.Status
			if strings.TrimSpace(stripANSI(mid)) == "" {
				mid = stripANSI(r.Abbrev)
			}
			num := padPlain(strings.TrimSpace(r.Num), maxNum)
			box.ChoiceRow(trail, num, opt, mid, "")
			if strings.TrimSpace(stripANSI(r.HelpLine)) != "" {
				box.ChoiceHelpRow(trail, helpLineFn(r.HelpLine))
			}
		}
	}
	for fi, r := range footer {
		trail := ui.SectionTrail{LastSection: true, LastRow: fi == len(footer)-1}
		opt := choiceMenuOption(r.Desc)
		mid := r.Status
		if strings.TrimSpace(stripANSI(mid)) == "" {
			mid = stripANSI(r.Abbrev)
		}
		num := padPlain(strings.TrimSpace(r.Num), maxNum)
		box.ChoiceRow(trail, num, opt, mid, "")
	}
	box.Flush()
}

func choiceMenuHelpLine(s string) string {
	if strings.Contains(s, "\033") {
		return s
	}
	return ui.Gray(stripANSI(s))
}

func choiceEditHelpLine(s string) string {
	if strings.Contains(s, "\033") {
		return s
	}
	return ui.BrightBlue(stripANSI(s))
}

// printEditFieldPromptPanel shows the field being changed and its current value (red box).
func printEditFieldPromptPanel(field, hint, current string) {
	box := ui.NewMenuBox(ui.ChromeChoice)
	title := ui.HeaderChoice("changing "+strings.TrimSpace(field)) + ui.Dim(" | "+strings.TrimSpace(hint))
	box.Title(title)
	box.Plain(ui.SectionTrail{LastSection: true, LastRow: true}, editFieldCurrentLine(current))
	box.Flush()
}

func editFieldCurrentLine(current string) string {
	cur := strings.TrimSpace(current)
	if cur == "" {
		return "Currently set to: " + ui.Dim(ui.SymbolEmDash)
	}
	return "Currently set to: " + ui.BrightBlue(cur)
}

func printInfoMenuPanel(title string, sections []menuPanelSection) {
	printMenuPanel(ui.ChromeInfo, ui.HeaderInfo(title), sections)
}

func printSessionMenuPanel(title string, sections []menuPanelSection) {
	printMenuPanel(ui.ChromeSession, ui.HeaderInfo(title), sections)
}

func printInfoMenuPanelTitle(titleLine string, sections []menuPanelSection) {
	printMenuPanel(ui.ChromeInfo, titleLine, sections)
}

func printMenuPanel(mode ui.ChromeMode, titleLine string, sections []menuPanelSection) {
	box := ui.NewMenuBox(mode)
	box.Title(titleLine)
	secCount := len(sections)
	for si, sec := range sections {
		lastSec := si == secCount-1
		if strings.TrimSpace(sec.Label) != "" {
			box.Section(ui.SubHeaderInfo(stripANSI(sec.Label)), lastSec && len(sec.Rows) == 0)
		}
		rows := sec.Rows
		for ri, r := range rows {
			trail := ui.SectionTrail{LastSection: lastSec, LastRow: ri == len(rows)-1}
			opt := stripANSI(r.Desc)
			if opt == "" {
				opt = stripANSI(r.Num)
			}
			mid := stripANSI(r.Abbrev)
			if mid == "" {
				mid = stripANSI(r.Status)
			}
			box.InfoRow(trail, opt, mid, stripANSI(r.HelpLine))
		}
	}
	box.Flush()
}
