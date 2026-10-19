package config

import "testing"

func TestProjectTags(t *testing.T) {
	p := &Project{
		Commands: []Command{
			{Alias: "a", Shell: "a", Tags: []string{"git"}},
			{Alias: "b", Shell: "b", Tags: []string{"admin", "git"}},
		},
	}
	tags := p.ProjectTags()
	if len(tags) != 2 || tags[0] != "admin" {
		t.Fatalf("tags=%v", tags)
	}
}
