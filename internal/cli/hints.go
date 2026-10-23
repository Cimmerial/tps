package cli

import (
	"os/exec"
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

func allBuiltInHints(p *config.Project) []hintEntry {
	return []hintEntry{
		{key: "help", abbr: "h", label: "help", order: 0, show: func(*config.Project) bool { return true }},
		{key: "edit", abbr: "e", label: "edit", order: 1, show: func(*config.Project) bool { return true }},
		{key: "install_help", abbr: "ih", label: "install", order: 2, show: func(*config.Project) bool {
			_, err := exec.LookPath("tps")
			return err != nil
		}},
		{key: "quit", abbr: "q", label: "quit", order: 3, show: func(*config.Project) bool { return true }},
		{key: "switch", abbr: "p", label: "switch project", order: 4, show: func(*config.Project) bool { return true }},
		{key: "clear", abbr: "tps", label: "refresh", order: 5, show: func(*config.Project) bool { return true }},
		{key: "new_command", abbr: "nc", label: "new cmd", order: 6, show: func(*config.Project) bool { return true }},
		{key: "options", abbr: "o", label: "options", order: 7, show: func(*config.Project) bool { return true }},
		{key: "update", abbr: "u", label: "update", order: 7, show: func(p *config.Project) bool {
			cfg, err := config.Load()
			if err != nil {
				return canSelfUpdate(p.Path)
			}
			return devRepoRoot(cfg, p) != ""
		}},
		{key: "ship", abbr: "g", label: "git ship", order: 8, show: func(*config.Project) bool { return true }},
		{key: "git", abbr: "git", label: "", order: 9, show: func(*config.Project) bool { return true }},
		{key: "terminals", abbr: "t", label: "terminals", order: 10, show: func(*config.Project) bool { return true }},
	}
}

func isPinnedHintKey(key string) bool {
	switch key {
	case "help", "edit", "options", "clear":
		return true
	default:
		return false
	}
}

func builtInMainHints(p *config.Project) []hintEntry {
	out := make([]hintEntry, 0, 12)
	for _, h := range allBuiltInHints(p) {
		if !h.show(p) || isPinnedHintKey(h.key) {
			continue
		}
		out = append(out, h)
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

func builtInPinnedHints(p *config.Project) []hintEntry {
	var out []hintEntry
	for _, key := range []string{"help", "edit", "options", "clear"} {
		for _, h := range allBuiltInHints(p) {
			if h.key == key && h.show(p) {
				out = append(out, h)
				break
			}
		}
	}
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
	dot := ui.Dim(" · ")
	mainParts := make([]string, 0, 12)
	for _, h := range builtInMainHints(p) {
		mainParts = append(mainParts, hintPart(h.abbr, h.label))
	}
	pinnedParts := make([]string, 0, 4)
	for _, h := range builtInPinnedHints(p) {
		pinnedParts = append(pinnedParts, hintPart(h.abbr, h.label))
	}

	mainLine := strings.Join(mainParts, dot)
	pinnedLine := strings.Join(pinnedParts, dot)
	line := mainLine
	if pinnedLine != "" {
		bar := ui.Dim(" | ")
		if line == "" {
			line = pinnedLine
		} else {
			line = line + bar + pinnedLine
		}
	}
	box := ui.NewMenuBox(ui.ChromeSession)
	box.Title(ui.HeaderInfo("TPS CMDS:"))
	box.Plain(ui.SectionTrail{LastSection: true, LastRow: true}, line)
	box.Flush()
}

func metaKeyForAction(action InputAction) string {
	switch action {
	case ActionHelp, ActionClear:
		return ""
	case ActionQuit:
		return "quit"
	case ActionEdit:
		return "edit"
	case ActionConfig:
		return "config"
	case ActionUpdate:
		return "update"
	case ActionGitShip:
		return "ship"
	case ActionExport:
		return "export"
	case ActionImport:
		return "import"
	case ActionSetup:
		return "setup"
	case ActionConfigExamples:
		return "config_examples"
	case ActionGit:
		return "git"
	case ActionSwitchProject:
		return "switch"
	case ActionSwitchProjectBack:
		return "switch_back"
	case ActionInstallHelp:
		return "install_help"
	case ActionNewCommand:
		return "new_command"
	case ActionTerminalCheck:
		return "terminals"
	case ActionOptions:
		return "options"
	case ActionCloseAllTerms:
		return "close_all_terms"
	case ActionOpenInIDE:
		return "open_in_ide"
	case ActionInstallUtilities:
		return "install_utilities"
	default:
		return ""
	}
}
