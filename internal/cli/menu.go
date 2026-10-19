package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

type InputAction int

const (
	ActionRun InputAction = iota
	ActionHelp
	ActionUpdate
	ActionQuit
	ActionClear
	ActionConfig
	ActionEdit
	ActionNoop
	ActionGitShip
	ActionGit
)

type MenuInput struct {
	Action InputAction
	Alias  string
	Text   string
}

func commandLabel(c config.Command) string {
	if c.Description != "" {
		return c.Description
	}
	if c.Abbrev != "" {
		return c.Abbrev
	}
	return c.Alias
}

func printSessionHeader(p *config.Project) {
	fmt.Printf("project: %s\n", p.Name)
	printCompactMenu(p)
	printMenuHints(p)
}

func printCompactMenu(p *config.Project) {
	top := p.TopCommands(config.TopCommandSlots)
	var parts []string
	for i, c := range top {
		parts = append(parts, fmt.Sprintf("%d > %s", i+1, commandLabel(c)))
	}
	if len(parts) > 0 {
		fmt.Println(strings.Join(parts, "  "))
	}
}

func printMenuHints(p *config.Project) {
	rest := len(p.Commands) - len(p.TopCommands(config.TopCommandSlots))
	if rest < 0 {
		rest = 0
	}
	hint := "h help · q quit · clear · edit · c config"
	if rest > 0 {
		hint = fmt.Sprintf("+%d more · %s", rest, hint)
	}
	if canSelfUpdate(p.Path) {
		hint += " · u update"
	}
	hint += " · gs ship · git …"
	fmt.Println(hint)
}

func printFullHelp(p *config.Project) {
	sections := []helpSection{
		{
			title: "session",
			rows: []helpRow{
				{"clear", "tps", "reset screen and menu"},
				{"quit", "q", "exit tps"},
			},
		},
		{
			title: "config",
			rows: []helpRow{
				{"config", "c", "open config.yaml in Cursor/VS Code or $EDITOR"},
				{"edit", "", "change a command (prompts for shell)"},
				{"shell tips", "", "chain with && or use YAML | for multiple lines"},
			},
		},
	}

	if canSelfUpdate(p.Path) {
		sections = append(sections, helpSection{
			title: "update",
			rows: []helpRow{
				{"update", "u", "go install ./cmd/tps and restart tps"},
			},
		})
	}

	sections = append(sections, helpSection{
		title: "git",
		rows: []helpRow{
			{"ship", "gs", "add all, commit, push (gs your message here)"},
			{"status", "", "git status"},
			{"git …", "", "any git command (git diff, git log -1)"},
		},
	})

	sections = append(sections, helpSection{
		title: "help",
		rows: []helpRow{
			{"help", "h", "this list"},
		},
	})

	custom := make([]helpRow, 0, len(p.Commands))
	for _, c := range p.Commands {
		custom = append(custom, projectCommandHelpRow(c))
	}
	if len(custom) > 0 {
		sections = append(sections, helpSection{
			title: "project commands",
			rows:  custom,
		})
	}

	printGroupedHelp(sections)
}

func projectCommandHelpRow(c config.Command) helpRow {
	desc := c.Description
	if desc == "" {
		desc = c.Shell
	}
	return helpRow{command: c.Alias, abbrev: c.Abbrev, description: desc}
}

type helpRow struct {
	command     string
	abbrev      string
	description string
}

type helpSection struct {
	title string
	rows  []helpRow
}

func printGroupedHelp(sections []helpSection) {
	maxCmd, maxAbbrev, maxDesc := 0, 0, 0
	for _, sec := range sections {
		for _, r := range sec.rows {
			if len(r.command) > maxCmd {
				maxCmd = len(r.command)
			}
			if len(r.abbrev) > maxAbbrev {
				maxAbbrev = len(r.abbrev)
			}
			if len(r.description) > maxDesc {
				maxDesc = len(r.description)
			}
		}
	}

	tableWidth := 0
	for _, sec := range sections {
		for _, r := range sec.rows {
			line := fmt.Sprintf("%-*s | %-*s | %s", maxCmd, r.command, maxAbbrev, r.abbrev, r.description)
			if len(line) > tableWidth {
				tableWidth = len(line)
			}
		}
	}

	for i, sec := range sections {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("  %s\n", sec.title)
		barLen := tableWidth
		if len(sec.title) > barLen {
			barLen = len(sec.title)
		}
		fmt.Printf("  %s\n", strings.Repeat("-", barLen))
		for _, r := range sec.rows {
			fmt.Printf("  %-*s | %-*s | %s\n", maxCmd, r.command, maxAbbrev, r.abbrev, r.description)
		}
	}
}

func canSelfUpdate(projectRoot string) bool {
	_, err := os.Stat(filepath.Join(projectRoot, "cmd", "tps"))
	return err == nil
}

func resolveMenuInput(p *config.Project, line string) (MenuInput, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return MenuInput{Action: ActionNoop}, nil
	}
	low := strings.ToLower(line)
	switch low {
	case "q", "quit", "exit":
		return MenuInput{Action: ActionQuit}, nil
	}
	if fields := strings.Fields(line); len(fields) >= 2 && strings.EqualFold(fields[0], "gs") {
		msg := strings.TrimSpace(strings.Join(fields[1:], " "))
		if msg == "" {
			return MenuInput{}, fmt.Errorf("gs needs a message (example: gs fix the menu)")
		}
		return MenuInput{Action: ActionGitShip, Text: msg}, nil
	}
	if low == "gs" {
		return MenuInput{}, fmt.Errorf("gs needs a message (example: gs fix the menu)")
	}
	if strings.HasPrefix(low, "git ") {
		rest := strings.TrimSpace(line[3:])
		if rest == "" {
			return MenuInput{}, fmt.Errorf("git needs a subcommand (example: git status)")
		}
		return MenuInput{Action: ActionGit, Text: rest}, nil
	}
	if c := p.FindCommandByAbbrev(line); c != nil {
		return MenuInput{Action: ActionRun, Alias: c.Alias}, nil
	}
	switch low {
	case "h", "help", "?":
		return MenuInput{Action: ActionHelp}, nil
	case "clear", "tps":
		return MenuInput{Action: ActionClear}, nil
	case "c", "config", "yaml":
		return MenuInput{Action: ActionConfig}, nil
	case "u", "update":
		if canSelfUpdate(p.Path) {
			return MenuInput{Action: ActionUpdate}, nil
		}
		return MenuInput{}, fmt.Errorf("update only works in a repo with cmd/tps")
	}

	if strings.HasPrefix(low, "edit ") {
		alias := strings.TrimSpace(line[5:])
		if alias == "" {
			return MenuInput{Action: ActionEdit}, nil
		}
		return MenuInput{Action: ActionEdit, Alias: alias}, nil
	}
	if low == "edit" {
		return MenuInput{Action: ActionEdit}, nil
	}

	top := p.TopCommands(config.TopCommandSlots)
	if n, err := strconv.Atoi(line); err == nil {
		if n >= 1 && n <= len(top) {
			return MenuInput{Action: ActionRun, Alias: top[n-1].Alias}, nil
		}
		return MenuInput{}, fmt.Errorf("pick 1–%d, a name, or h (h for help)", len(top))
	}

	if c := p.FindCommand(line); c != nil {
		return MenuInput{Action: ActionRun, Alias: c.Alias}, nil
	}
	for _, c := range p.Commands {
		if strings.EqualFold(c.Alias, line) {
			return MenuInput{Action: ActionRun, Alias: c.Alias}, nil
		}
	}
	return MenuInput{}, fmt.Errorf("unknown command %q (h for help)", line)
}
