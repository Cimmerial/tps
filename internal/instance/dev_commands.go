package instance

import (
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

// legacyDevShellBuiltinAliases are mistaken YAML copies of TPS built-ins (u, vite).
var legacyDevShellBuiltinAliases = map[string]bool{
	"update":          true,
	"close all terms": true,
}

// DefaultDevProjectCommands are shell-less commands merged into the linked tps dev repo.
func DefaultDevProjectCommands() []config.Command {
	return []config.Command{
		{
			Alias:       "sym",
			Abbrev:      "sym",
			Description: "ui symbols (sym check, sym braille, …)",
			Tags:        []string{"tps"},
			Handler:     "sym",
		},
		{
			Alias:       "push-release",
			Abbrev:      "rel",
			Description: "commit first, then tag + push (CI release); rel -v0.2.0 -y",
			Tags:        []string{"tps"},
			Handler:     "push-release",
		},
	}
}

// EnsureDefaultDevProjectCommands adds missing default commands on the tps dev repo project.
// Returns true when the project was modified.
func EnsureDefaultDevProjectCommands(p *config.Project) (bool, error) {
	if p == nil || !IsDevRepo(p.Path) {
		return false, nil
	}
	changed := pruneLegacyDevShellBuiltinCommands(p)
	for _, def := range DefaultDevProjectCommands() {
		if existing := p.FindCommand(def.Alias); existing != nil {
			if c := syncDevProjectCommand(existing, def); c {
				changed = true
			}
			if strings.TrimSpace(existing.Handler) == "" && strings.TrimSpace(def.Handler) != "" {
				existing.Handler = def.Handler
				if existing.Description == "" {
					existing.Description = def.Description
				}
				if existing.Abbrev == "" {
					existing.Abbrev = def.Abbrev
				}
				if len(existing.Tags) == 0 && len(def.Tags) > 0 {
					existing.Tags = def.Tags
				}
				changed = true
			}
			continue
		}
		cmd := def
		if err := p.AddCommand(cmd); err != nil {
			return changed, err
		}
		changed = true
	}
	return changed, nil
}

func syncDevProjectCommand(existing *config.Command, def config.Command) bool {
	if existing == nil {
		return false
	}
	changed := false
	if strings.TrimSpace(def.Handler) != "" {
		if existing.Handler != def.Handler {
			existing.Handler = def.Handler
			changed = true
		}
		if existing.Shell != "" {
			existing.Shell = ""
			changed = true
		}
	}
	if def.Abbrev != "" && existing.Abbrev != def.Abbrev {
		existing.Abbrev = def.Abbrev
		changed = true
	}
	if def.Description != "" && existing.Description != def.Description {
		existing.Description = def.Description
		changed = true
	}
	if len(def.Tags) > 0 && !sameStringSlice(existing.Tags, def.Tags) {
		existing.Tags = append([]string(nil), def.Tags...)
		changed = true
	}
	return changed
}

func sameStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func pruneLegacyDevShellBuiltinCommands(p *config.Project) bool {
	if p == nil {
		return false
	}
	changed := false
	kept := make([]config.Command, 0, len(p.Commands))
	for _, c := range p.Commands {
		if isLegacyDevShellBuiltinAlias(c.Alias) {
			changed = true
			continue
		}
		kept = append(kept, c)
	}
	if changed {
		p.Commands = kept
	}
	return changed
}

func isLegacyDevShellBuiltinAlias(alias string) bool {
	return legacyDevShellBuiltinAliases[strings.ToLower(strings.TrimSpace(alias))]
}

// IsNonEditableDevProjectCommand reports commands on the tps dev repo that must not appear in edit.
func IsNonEditableDevProjectCommand(p *config.Project, alias string) bool {
	if p == nil || !IsDevRepo(p.Path) {
		return false
	}
	if isLegacyDevShellBuiltinAlias(alias) {
		return true
	}
	return IsBuiltinDevCommand(p, alias)
}

// IsBuiltinDevCommand reports whether the project command runs in-process (handler set).
func IsBuiltinDevCommand(p *config.Project, alias string) bool {
	if p == nil || !IsDevRepo(p.Path) {
		return false
	}
	cmd := p.FindCommand(alias)
	return cmd != nil && strings.TrimSpace(cmd.Handler) != ""
}
