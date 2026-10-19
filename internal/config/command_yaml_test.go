package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCommandAbbreviationYAML(t *testing.T) {
	var c Command
	err := yaml.Unmarshal([]byte(`
alias: hello
abbreviation: h
shell: echo hello-tps
`), &c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Abbrev != "h" {
		t.Fatalf("expected abbrev h, got %q", c.Abbrev)
	}
	out, err := yaml.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(string(out), "abbreviation: h") {
		t.Fatalf("marshal should use abbreviation key: %s", out)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
