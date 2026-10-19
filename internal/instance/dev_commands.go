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
