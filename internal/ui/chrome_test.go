package ui

import "testing"

func TestFormatMenuColumnsChoice(t *testing.T) {
	s := formatMenuColumns("1", "run", "r", "local", true)
	if s == "" {
		t.Fatal("empty")
	}
	if VisibleWidth(s) < 10 {
		t.Fatalf("short: %q", s)
	}
}

func TestPrintResponseRuleWidth(t *testing.T) {
	w := boxWidth()
	if w < 20 {
		t.Fatalf("width %d", w)
	}
}

func TestWrapBoxContent(t *testing.T) {
	inner := 10
	lines := WrapBoxContent("hello world foo", inner)
	if len(lines) < 2 {
		t.Fatalf("expected wrap, got %v", lines)
	}
}
