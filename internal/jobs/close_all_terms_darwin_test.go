//go:build darwin

package jobs

import "testing"

func TestClosedTabsSince(t *testing.T) {
	before := []terminalTab{
		{TTY: "/dev/ttys001", Label: "npm run dev"},
		{TTY: "/dev/ttys002", Label: "bash"},
		{TTY: "/dev/ttys003", Label: "vite"},
	}
	after := []terminalTab{{TTY: "ttys002", Label: "bash"}}
	got := closedTabsSince(before, after)
	if got.Count != 2 {
		t.Fatalf("count %d want 2", got.Count)
	}
	if len(got.Labels) != 2 || got.Labels[0] != "npm run dev" || got.Labels[1] != "vite" {
		t.Fatalf("labels %v", got.Labels)
	}
}

func TestTerminalTabDisplayLabel(t *testing.T) {
	if (terminalTab{TTY: "ttys001"}).displayLabel() != "ttys001" {
		t.Fatal("expected tty fallback")
	}
}

func TestCombineCloseResultsFallsBackToTrackedJobs(t *testing.T) {
	got := combineCloseResults([]string{"run", "tunnel"}, CloseAllTermsResult{})
	if got.Count != 2 || len(got.Labels) != 2 {
		t.Fatalf("got %+v", got)
	}
	got = combineCloseResults([]string{"run"}, CloseAllTermsResult{Count: 1, Labels: []string{"npm"}})
	if got.Count != 1 || got.Labels[0] != "npm" {
		t.Fatalf("prefer tabs: %+v", got)
	}
}
