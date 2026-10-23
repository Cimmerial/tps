package cli

import "strings"

const shipCommitMenuSlots = 2

type shipCommitDraft struct {
	Title string
	Body  string
}

func newShipCommitDraft(initial string) *shipCommitDraft {
	return &shipCommitDraft{Title: strings.TrimSpace(initial)}
}

func (d *shipCommitDraft) CommitMessage() string {
	title := strings.TrimSpace(d.Title)
	body := strings.TrimSpace(d.Body)
	if title == "" {
		return ""
	}
	if body == "" {
		return title
	}
	return title + "\n\n" + body
}

func shipCommitSubmenuRows(d *shipCommitDraft) []numberedMenuRow {
	nameSet := strings.TrimSpace(d.Title) != ""
	descSet := strings.TrimSpace(d.Body) != ""
	return []numberedMenuRow{
		{Num: "1", Status: editFieldSetMark(nameSet), Desc: "name", HelpLine: editHelpIf(nameSet, d.Title)},
		{Num: "2", Status: editFieldSetMark(descSet), Desc: "description", HelpLine: editHelpIf(descSet, d.Body)},
	}
}

func shipFileDisplayNumbersOffset(paths []string, offset int) map[string]int {
	base := shipFileDisplayNumbers(paths)
	if offset == 0 {
		return base
	}
	out := make(map[string]int, len(base))
	for p, n := range base {
		out[p] = n + offset
	}
	return out
}

func shipFilePathForDisplayNumberOffset(paths []string, displayNum, offset int) (string, bool) {
	if displayNum <= offset {
		return "", false
	}
	return shipFilePathForDisplayNumber(paths, displayNum-offset)
}
