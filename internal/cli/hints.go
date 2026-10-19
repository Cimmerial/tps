package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

type hintEntry struct {
	key   string
	abbr  string
	label string
	order int
	show  func(*config.Project) bool
}

func builtInHints(p *config.Project) []hintEntry {
	all := []hintEntry{
		{key: "help", abbr: "h", label: "help", order: 0, show: func(*config.Project) bool { return true }},
		{key: "install_help", abbr: "ih", label: "install", order: 1, show: func(*config.Project) bool { return true }},
		{key: "quit", abbr: "q", label: "quit", order: 2, show: func(*config.Project) bool { return true }},
		{key: "switch", abbr: "p", label: "project", order: 3, show: func(*config.Project) bool { return true }},
		{key: "clear", abbr: "tps", label: "clear", order: 4, show: func(*config.Project) bool { return true }},
		{key: "edit", abbr: "e", label: "edit", order: 5, show: func(*config.Project) bool { return true }},
		{key: "new_command", abbr: "nc", label: "new cmd", order: 5, show: func(*config.Project) bool { return true }},
		{key: "config", abbr: "c", label: "config", order: 6, show: func(*config.Project) bool { return true }},
		{key: "update", abbr: "u", label: "update", order: 7, show: func(p *config.Project) bool {
			cfg, err := config.Load()
			if err != nil {
				return canSelfUpdate(p.Path)
			}
			return devRepoRoot(cfg, p) != ""
		}},
		{key: "ship", abbr: "gs", label: "ship", order: 8, show: func(*config.Project) bool { return true }},
		{key: "git", abbr: "git", label: "", order: 9, show: func(*config.Project) bool { return true }},
		{key: "terminals", abbr: "tc", label: "terminals", order: 10, show: func(*config.Project) bool { return true }},
	}
	out := make([]hintEntry, 0, len(all))
	for _, h := range all {
		if h.show(p) {
			out = append(out, h)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		ci := p.MetaUseCount(out[i].key)
		cj := p.MetaUseCount(out[j].key)
		if ci != cj {
			return ci > cj
		}
		return out[i].order < out[j].order
	})
	return out
}

func hintPart(abbr, label string) string {
	abbr = strings.TrimSpace(abbr)
	label = strings.TrimSpace(label)
	if abbr == "" {
		return ui.Dim(label)
	}
	if label == "" || strings.EqualFold(abbr, label) {
		return ui.Abbrev(abbr)
	}
	return ui.Abbrev(abbr) + ui.Dim(" "+label)
}

func printMenuHints(p *config.Project) {
	rest := len(p.Commands) - len(p.TopCommands(config.TopCommandSlots))
	if rest < 0 {
		rest = 0
	}

	sep := ui.Dim(" · ")
	parts := make([]string, 0, 12)
	for _, h := range builtInHints(p) {
		parts = append(parts, hintPart(h.abbr, h.label))
	}
	line := strings.Join(parts, sep)
	if rest > 0 {
		line += sep + ui.Dim(fmt.Sprintf("+%d more", rest))
	}
	fmt.Printf("%s %s\n", ui.Bold("TPS CMDS:"), line)
}

func metaKeyForAction(action InputAction) string {
	switch action {
	case ActionHelp:
		return "help"
	case ActionQuit:
		return "quit"
	case ActionClear:
		return "clear"
	case ActionEdit:
		return "edit"
	case ActionConfig:
		return "config"
	case ActionUpdate:
		return "update"
	case ActionGitShip:
		return "ship"
	case ActionGit:
		return "git"
	case ActionSwitchProject:
		return "switch"
	case ActionInstallHelp:
		return "install_help"
	case ActionNewCommand:
		return "new_command"
	case ActionTerminalCheck:
		return "terminals"
	default:
		return ""
	}
}
