package config

import (
	"fmt"
	"sort"
	"strings"
)

// CommandGroup is a tag folder of project commands (empty tag = untagged).
type CommandGroup struct {
	Tag      string
	Commands []Command
}

// PrimaryTag is the folder name for grouping (first tag, if any).
func (c Command) PrimaryTag() string {
	for _, t := range c.Tags {
		t = strings.TrimSpace(t)
		if t != "" {
			return t
		}
	}
	return ""
}

func parseTagsField(v interface{}) []string {
	if v == nil {
		return nil
	}
	switch x := v.(type) {
	case string:
		return splitTagList(x)
	case []interface{}:
		out := make([]string, 0, len(x))
		for _, item := range x {
			for _, p := range splitTagList(fmt.Sprint(item)) {
				out = append(out, p)
			}
		}
		return out
	case []string:
		out := make([]string, 0, len(x))
		for _, s := range x {
			for _, p := range splitTagList(s) {
				out = append(out, p)
			}
		}
		return out
	default:
		return splitTagList(fmt.Sprint(v))
	}
}

// ParseTagList splits a comma-separated tag string (CLI / import).
func ParseTagList(s string) []string {
	return splitTagList(s)
}

func splitTagList(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// ProjectTags returns distinct non-empty tags used by commands, sorted.
func (p *Project) ProjectTags() []string {
	if p == nil {
		return nil
	}
	seen := map[string]bool{}
	for _, c := range p.Commands {
		for _, t := range c.Tags {
			t = strings.TrimSpace(t)
			if t != "" {
				seen[t] = true
			}
		}
	}
	out := make([]string, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// RenameTag renames old to new in every command's tag list; returns commands updated.
func (p *Project) RenameTag(old, new string) int {
	old = strings.TrimSpace(old)
	new = strings.TrimSpace(new)
	if old == "" || new == "" || old == new {
		return 0
	}
	n := 0
	for i := range p.Commands {
		changed := false
		for j, t := range p.Commands[i].Tags {
			if strings.TrimSpace(t) == old {
				p.Commands[i].Tags[j] = new
				changed = true
			}
		}
		if changed {
			n++
		}
	}
	return n
}

// CommandsByTag returns commands grouped by PrimaryTag; untagged first, then tags A–Z.
func (p *Project) CommandsByTag() []CommandGroup {
	if p == nil || len(p.Commands) == 0 {
		return nil
	}
	byTag := map[string][]Command{}
	var untagged []Command
	for _, rc := range p.ranked() {
		c := rc.Command
		tag := c.PrimaryTag()
		if tag == "" {
			untagged = append(untagged, c)
			continue
		}
		byTag[tag] = append(byTag[tag], c)
	}
	tags := make([]string, 0, len(byTag))
	for t := range byTag {
		tags = append(tags, t)
	}
	sort.Strings(tags)

	groups := make([]CommandGroup, 0, 1+len(tags))
	if len(untagged) > 0 {
		groups = append(groups, CommandGroup{Commands: untagged})
	}
	for _, tag := range tags {
		groups = append(groups, CommandGroup{Tag: tag, Commands: byTag[tag]})
	}
	return groups
}
