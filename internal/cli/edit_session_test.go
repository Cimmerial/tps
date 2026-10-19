package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func TestEditSessionChangesRename(t *testing.T) {
	snap := commandSnap{Shell: "npm run x", Abbrev: "ic"}
	before := map[string]commandSnap{"inspect Cimmerial classes": snap}
	after := map[string]commandSnap{"inspect classes": snap}
	changes := editSessionChanges(before, after)
	if len(changes) != 1 {
		t.Fatalf("changes = %v", changes)
	}
	if !strings.Contains(changes[0], "inspect Cimmerial classes") || !strings.Contains(changes[0], "inspect classes") {
		t.Fatalf("expected rename line, got %q", changes[0])
	}
}

func TestEditSessionChanges(t *testing.T) {
	before := map[string]commandSnap{
		"run": {Shell: "npm run dev", Abbrev: "r"},
	}
	after := map[string]commandSnap{
		"run":  {Shell: "npm start", Abbrev: "r"},
		"test": {Shell: "go test ./..."},
	}
	changes := editSessionChanges(before, after)
	if len(changes) != 2 {
		t.Fatalf("changes = %v", changes)
	}
	same := editSessionChanges(before, before)
	if len(same) != 0 {
		t.Fatalf("expected no changes, got %v", same)
	}
}

func TestEditCommandListSkipsNonEditableDevCommands(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "cmd", "tps"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := &config.Project{
		Path: root,
		Commands: []config.Command{
			{Alias: "run", Shell: "go test ./..."},
			{Alias: "sym", Handler: "sym"},
			{Alias: "update", Shell: "go install ./cmd/tps"},
			{Alias: "close all terms", Abbrev: "vite", Shell: "osascript"},
		},
	}
	list := editCommandList(p)
	if len(list) != 1 || list[0].Alias != "run" {
		t.Fatalf("edit list = %+v", list)
	}
}

func TestEditCommandListOrder(t *testing.T) {
	p := &config.Project{
		Commands: []config.Command{
			{Alias: "z", Tags: []string{"git"}},
			{Alias: "a"},
			{Alias: "b", Tags: []string{"admin"}},
		},
	}
	list := editCommandList(p)
	if len(list) != 3 {
		t.Fatalf("len=%d", len(list))
	}
	if list[0].Alias != "a" {
		t.Fatalf("untagged first: %q", list[0].Alias)
	}
}

func TestEditFieldSetMark(t *testing.T) {
	if editFieldSetMark(true) == editFieldSetMark(false) {
		t.Fatal("expected different set/unset marks")
	}
}

func TestEditCommandSubmenuRowsStatus(t *testing.T) {
	cmd := config.Command{
		Alias:       "mtm",
		Shell:       "git merge",
		Abbrev:      "m",
		Description: "mtm",
		Detach:      true,
	}
	rows := editCommandSubmenuRows(cmd)
	if len(rows) != 10 {
		t.Fatalf("len=%d", len(rows))
	}
	if rows[7].Desc != ui.Dim("kill detachment after completion") {
		t.Fatalf("row 8 = %q", rows[7].Desc)
	}
	if rows[0].HelpLine == "" {
		t.Fatal("expected shell help")
	}
	if rows[6].HelpLine == "" {
		t.Fatal("expected detach help")
	}
}

func TestEditCommandSubmenuRowsNoCloseWhenForeground(t *testing.T) {
	cmd := config.Command{Alias: "x", Shell: "echo hi", Detach: false}
	rows := editCommandSubmenuRows(cmd)
	if len(rows) != 9 {
		t.Fatalf("len=%d", len(rows))
	}
	for _, r := range rows {
		if r.Desc == ui.Dim("kill detachment after completion") {
			t.Fatal("close-on-complete row should be hidden when not detached")
		}
	}
}

func TestEditNextDraftAlias(t *testing.T) {
	p := &config.Project{}
	if editNextDraftAlias(p) != "new" {
		t.Fatalf("expected new")
	}
	p.Commands = []config.Command{{Alias: "new", Shell: "x"}}
	if editNextDraftAlias(p) != "new-2" {
		t.Fatalf("expected new-2")
	}
}

func TestProjectRenameTag(t *testing.T) {
	p := &config.Project{
		Commands: []config.Command{
			{Alias: "x", Shell: "x", Tags: []string{"git", "admin"}},
			{Alias: "y", Shell: "y", Tags: []string{"git"}},
		},
	}
	n := p.RenameTag("git", "vcs")
	if n != 2 {
		t.Fatalf("n=%d", n)
	}
	if p.Commands[0].Tags[0] != "vcs" || p.Commands[1].Tags[0] != "vcs" {
		t.Fatalf("tags=%v %v", p.Commands[0].Tags, p.Commands[1].Tags)
	}
}
