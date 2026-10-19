package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSetShellLineClearsMultiShell(t *testing.T) {
	c := Command{Alias: "stack", Shell: "a"}
	c.shells = &shellStore{generic: map[int]string{1: "one", 2: "two"}}
	c.SetShellLine("npm test")
	if c.Shell != "npm test" {
		t.Fatalf("shell=%q", c.Shell)
	}
	if c.shells != nil {
		t.Fatal("expected multi-shell cleared")
	}
}

func TestSetShellByKeyPreservesOthers(t *testing.T) {
	var c Command
	if err := yaml.Unmarshal([]byte(`
alias: mtm
shell: git checkout main
`), &c); err != nil {
		t.Fatal(err)
	}
	if err := c.SetShellByKey("shell2", "echo two"); err != nil {
		t.Fatal(err)
	}
	if err := c.SetShellByKey("darwin_shell", "open -a Terminal"); err != nil {
		t.Fatal(err)
	}
	fields := c.AllShellFields()
	if len(fields) != 3 {
		t.Fatalf("fields=%v", fields)
	}
	if err := c.RemoveShellByKey("shell2"); err != nil {
		t.Fatal(err)
	}
	fields = c.AllShellFields()
	if len(fields) != 2 {
		t.Fatalf("after remove: %v", fields)
	}
}

func TestHasAnyShellLine(t *testing.T) {
	c := Command{Alias: "x"}
	if c.HasAnyShellLine() {
		t.Fatal("expected empty")
	}
	c.Shell = "go test ./..."
	if !c.HasAnyShellLine() {
		t.Fatal("expected shell")
	}
}
