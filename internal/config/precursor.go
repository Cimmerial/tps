package config

import (
	"fmt"
	"strings"
)

// CommandRunChain returns precursor commands first, then alias (e.g. tunnel → run).
func (p *Project) CommandRunChain(alias string) ([]string, error) {
	if p == nil {
		return nil, fmt.Errorf("no project")
	}
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return nil, fmt.Errorf("empty command")
	}
	target := p.ResolveCommandRef(alias)
	if target == nil {
		return nil, fmt.Errorf("unknown command %q", alias)
	}
	alias = target.Alias

	var precursors []string
	seen := map[string]bool{alias: true}
	walk := alias
	for {
		cmd := p.FindCommand(walk)
		if cmd == nil {
			return nil, fmt.Errorf("unknown command %q", walk)
		}
		ref := strings.TrimSpace(cmd.Precursor)
		if ref == "" {
			break
		}
		pre := p.ResolveCommandRef(ref)
		if pre == nil {
			return nil, fmt.Errorf("precursor %q not found for %q", ref, walk)
		}
		if seen[pre.Alias] {
			return nil, fmt.Errorf("precursor cycle involving %q", pre.Alias)
		}
		seen[pre.Alias] = true
		precursors = append([]string{pre.Alias}, precursors...)
		walk = pre.Alias
	}
	out := append(precursors, alias)
	return out, nil
}

// ResolveCommandRef finds a command by alias or abbreviation.
func (p *Project) ResolveCommandRef(name string) *Command {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	}
	if c := p.FindCommand(name); c != nil {
		return c
	}
	for i := range p.Commands {
		if strings.EqualFold(p.Commands[i].Alias, name) {
			return &p.Commands[i]
		}
	}
	return p.FindCommandByAbbrev(name)
}
