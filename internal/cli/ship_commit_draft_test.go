package cli

import "testing"

func TestShipCommitDraftCommitMessage(t *testing.T) {
	d := &shipCommitDraft{Title: "fix menu", Body: "details"}
	if got := d.CommitMessage(); got != "fix menu\n\ndetails" {
		t.Fatalf("got %q", got)
	}
	d.Body = ""
	if got := d.CommitMessage(); got != "fix menu" {
		t.Fatalf("title only: got %q", got)
	}
	d.Title = ""
	if got := d.CommitMessage(); got != "" {
		t.Fatalf("empty title: got %q", got)
	}
}

func TestShipFileDisplayNumbersOffset(t *testing.T) {
	paths := []string{"a.go", "b.go"}
	nums := shipFileDisplayNumbersOffset(paths, shipCommitMenuSlots)
	if nums["a.go"] != 3 || nums["b.go"] != 4 {
		t.Fatalf("nums %v", nums)
	}
	p, ok := shipFilePathForDisplayNumberOffset(paths, 3, shipCommitMenuSlots)
	if !ok || p != "a.go" {
		t.Fatalf("path=%q ok=%v", p, ok)
	}
}
