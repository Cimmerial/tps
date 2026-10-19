package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCommandTagsYAML(t *testing.T) {
	var c Command
	err := yaml.Unmarshal([]byte(`
alias: merge to main
shell: git checkout main
abbreviation: mtm
tags:
  - git
  - backend
`), &c)
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Tags) != 2 || c.Tags[0] != "git" {
		t.Fatalf("tags: %v", c.Tags)
	}
	if c.PrimaryTag() != "git" {
		t.Fatalf("primary tag: %q", c.PrimaryTag())
	}
}

func TestCommandsByTag(t *testing.T) {
	p := &Project{
		Commands: []Command{
			{Alias: "run", Abbrev: "r"},
			{Alias: "merge", Abbrev: "mtm", Tags: []string{"git"}},
			{Alias: "admin tool", Tags: []string{"admin"}},
		},
	}
	groups := p.CommandsByTag()
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
	if groups[0].Tag != "" || len(groups[0].Commands) != 1 {
		t.Fatalf("untagged group: %+v", groups[0])
	}
	if groups[1].Tag != "admin" || groups[2].Tag != "git" {
		t.Fatalf("tag order: %q, %q", groups[1].Tag, groups[2].Tag)
	}
}
