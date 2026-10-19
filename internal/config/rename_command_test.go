package config

import "testing"

func TestRenameCommand(t *testing.T) {
	p := &Project{
		Commands: []Command{
			{Alias: "old name", Shell: "x", Abbrev: "o"},
			{Alias: "child", Shell: "y", Precursor: "old name"},
		},
		CommandUsage: map[string]int{"old name": 3},
	}
	if err := p.RenameCommand("old name", "new name"); err != nil {
		t.Fatal(err)
	}
	if p.FindCommand("new name") == nil || p.FindCommand("old name") != nil {
		t.Fatal("rename failed")
	}
	if p.Commands[1].Precursor != "new name" {
		t.Fatalf("precursor=%q", p.Commands[1].Precursor)
	}
	if p.CommandUsage["new name"] != 3 {
		t.Fatalf("usage=%v", p.CommandUsage)
	}
}
