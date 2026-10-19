package ui

import "testing"

func TestVisibleWidth_ignoresANSI(t *testing.T) {
	s := Bold(Blue("TPS CMDS:"))
	if got := VisibleWidth(s); got != len("TPS CMDS:") {
		t.Fatalf("VisibleWidth(header) = %d, want %d", got, len("TPS CMDS:"))
	}
	part := Abbrev("h") + Dim(" help")
	if got := VisibleWidth(part); got != len("h help") {
		t.Fatalf("VisibleWidth(hint) = %d, want %d", got, len("h help"))
	}
}

func TestAlignLeftRight_padsRight(t *testing.T) {
	left := "TPS CMDS: q quit"
	right := " | h help"
	line := AlignLeftRight(left, right, 40)
	if VisibleWidth(line) != 40 {
		t.Fatalf("line visible width = %d, want 40; line=%q", VisibleWidth(line), line)
	}
	if !stringsHasSuffixSpacesBefore(line, right) {
		t.Fatalf("expected right segment at end: %q", line)
	}
}

func stringsHasSuffixSpacesBefore(line, right string) bool {
	return len(line) >= len(right) && line[len(line)-len(right):] == right
}
