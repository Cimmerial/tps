package addonmenu

import (
	"strings"
	"testing"
)

func TestActionSummaryLines(t *testing.T) {
	stderr := strings.Join([]string{
		"Delete original? [y/N]",
		"\x1b[32mwrote /tmp/out.jpg\x1b[0m",
		"",
		"\x1b[32mSUCCESS: done — 1 wrote, 0 unchanged, 0 skipped\x1b[0m",
	}, "\n")
	got := actionSummaryLines(stderr)
	if len(got) != 2 {
		t.Fatalf("got %d lines: %v", len(got), got)
	}
	if !strings.Contains(got[0], "wrote ") {
		t.Fatalf("first line: %q", got[0])
	}
	if !strings.Contains(got[1], "SUCCESS:") {
		t.Fatalf("second line: %q", got[1])
	}
}
