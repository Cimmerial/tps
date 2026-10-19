package config

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCommandMultiShellYAML(t *testing.T) {
	var c Command
	err := yaml.Unmarshal([]byte(`
alias: stack
shell1: echo one
shell2: echo two
`), &c)
	if err != nil {
		t.Fatal(err)
	}
	got := c.ResolvedShells()
	if len(got) != 2 || got[0] != "echo one" || got[1] != "echo two" {
		t.Fatalf("generic shells: %v", got)
	}
}

func TestCommandPlatformShellYAML(t *testing.T) {
	var c Command
	err := yaml.Unmarshal([]byte(`
alias: dev
windows_shell1: win-cmd
windows_shell2: win-cmd2
linux_shell: unix-cmd
`), &c)
	if err != nil {
		t.Fatal(err)
	}
	// Test store directly (GOOS-dependent ResolvedShells).
	st := c.shells
	if st == nil {
		t.Fatal("expected shell store")
	}
	if w := st.ordered(shellWindows); len(w) != 2 || w[0] != "win-cmd" {
		t.Fatalf("windows: %v", w)
	}
	if l := st.ordered(shellLinux); len(l) != 1 || l[0] != "unix-cmd" {
		t.Fatalf("linux: %v", l)
	}
}

func TestCommandMarshalShellKeys(t *testing.T) {
	st := shellStore{generic: map[int]string{1: "a", 2: "b"}}
	c := Command{Alias: "x", shells: &st}
	c.Shell = "a"
	out, err := yaml.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "shell: a") || !strings.Contains(s, "shell2: b") {
		t.Fatalf("marshal shells: %s", s)
	}
}
