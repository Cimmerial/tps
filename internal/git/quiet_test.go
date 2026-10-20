package git

import (
	"testing"

	"github.com/Cimmerial/tps/internal/ui"
)

func TestQuietGitArgs(t *testing.T) {
	ui.ApplySettings(ui.UISettings{SimplifiedOutput: true})
	defer ui.ApplySettings(ui.UISettings{})

	got := quietGitArgs([]string{"commit", "-m", "hi"})
	if got[1] != "-q" {
		t.Fatalf("commit args = %v", got)
	}
	got = quietGitArgs([]string{"push", "origin", "main"})
	if got[1] != "-q" {
		t.Fatalf("push args = %v", got)
	}
	got = quietGitArgs([]string{"status"})
	if len(got) != 1 || got[0] != "status" {
		t.Fatalf("status args = %v", got)
	}

	ui.ApplySettings(ui.UISettings{SimplifiedOutput: false})
	got = quietGitArgs([]string{"commit", "-m", "hi"})
	if len(got) != 3 {
		t.Fatalf("verbose commit args = %v", got)
	}
}
