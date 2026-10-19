package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
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
	ActionSwitchProject
	ActionInstallHelp
	ActionNewCommand
	ActionTerminalCheck
)

type MenuInput struct {
	Action      InputAction
	Alias       string
	Text        string
	UpdateAll   bool
	InstallCopy bool
	NewCmd      config.Command
	TcKillAll   bool
	TcKillID    int
}

func commandDetail(c config.Command) string {
	if c.Description != "" {
		return truncateRunes(c.Description, 28)
	}
	return c.Alias
}

func formatProjectCmdEntry(c config.Command) string {
	if c.Abbrev != "" {
		return ui.Abbrev(c.Abbrev) + ui.Dim(" "+commandDetail(c))
	}
	return ui.Bold(commandDetail(c))
}

func truncateRunes(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

func printSessionHeader(p *config.Project, restartBanner string) {
	name := ui.Bold(ui.Magenta(p.Name))
	switch restartBanner {
	case "updated":
		fmt.Printf("%s %s %s\n", ui.Dim("project:"), name, ui.Green("- tps updated"))
	case "uptodate":
		fmt.Printf("%s %s %s\n", ui.Dim("project:"), name, ui.Yellow("- already up to date"))
	default:
		fmt.Printf("%s %s\n", ui.Dim("project:"), name)
	}
	printCompactMenu(p)
	printMenuHints(p)
}

func printCompactMenu(p *config.Project) {
	top := p.TopCommands(config.TopCommandSlots)
	if len(top) == 0 {
		return
	}
	sep := ui.Dim(" · ")
	parts := make([]string, 0, len(top))
	for _, c := range top {
		parts = append(parts, formatProjectCmdEntry(c))
	}
	fmt.Printf("%s %s\n", ui.Bold("PRJCT CMDS:"), strings.Join(parts, sep))
}

func printFullHelp(p *config.Project) {
	sections := []helpSection{
		{
			title: "session",
			rows: []helpRow{
				{"clear", "tps", "reset screen and menu"},
				{"quit", "q", "exit tps"},
				{"project", "p", "switch linked project"},
			},
		},
		{
			title: "config",
			rows: []helpRow{
				{"config", "c", "open config.yaml in Cursor/VS Code or $EDITOR"},
				{"edit", "e", "change a command shell (e hello)"},
				{"new command", "nc", "nc -n name -a abbrev -c shell (-t terminal, -m capture)"},
				{"terminals", "tc", "tc · tc -k · tc -k 1"},
				{"shell tips", "", "chain with && or use YAML | for multiple lines"},
			},
		},
	}

	if canSelfUpdate(p.Path) {
		sections = append(sections, helpSection{
			title: "update",
			rows: []helpRow{
				{"update", "u", "go install ./cmd/tps and restart tps"},
				{"update all", "u -a", "install and copy to every tps on PATH and ./tps"},
			},
		})
	}

	sections = append(sections, helpSection{
		title: "git",
		rows: []helpRow{
			{"ship", "gs", "add all, commit, push (gs your message here)"},
			{"status", "", "git status"},
			{"git", "", "git status; or git diff, git log -1, …"},
		},
	})

	sections = append(sections, helpSection{
		title: "help",
		rows: []helpRow{
			{"help", "h", "this list"},
			{"install", "ih", "install & update tps on this or another machine"},
			{"install copy", "ih -copy", "copy install commands to clipboard"},
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
		fmt.Printf("  %s\n", ui.Bold(ui.Cyan(sec.title)))
		barLen := tableWidth
		if len(sec.title) > barLen {
			barLen = len(sec.title)
		}
		fmt.Printf("  %s\n", ui.Dim(strings.Repeat("-", barLen)))
		for _, r := range sec.rows {
			cmdCol := ui.Bold(padPlain(r.command, maxCmd))
			abbrPlain := padPlain(r.abbrev, maxAbbrev)
			abbrCol := abbrPlain
			if strings.TrimSpace(r.abbrev) != "" {
				abbrCol = ui.Abbrev(padPlain(r.abbrev, maxAbbrev))
			}
			fmt.Printf("  %s | %s | %s\n", cmdCol, abbrCol, ui.Dim(r.description))
		}
	}
}

func isInstallHelpToken(tok string) bool {
	switch strings.ToLower(tok) {
	case "ih", "install", "install-help":
		return true
	default:
		return false
	}
}

func padPlain(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
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
	if low == "git" {
		return MenuInput{Action: ActionGit, Text: "status"}, nil
	}
	if strings.HasPrefix(low, "git ") {
		rest := strings.TrimSpace(line[3:])
		if rest == "" {
			return MenuInput{Action: ActionGit, Text: "status"}, nil
		}
		return MenuInput{Action: ActionGit, Text: rest}, nil
	}
	switch low {
	case "h", "help", "?":
		return MenuInput{Action: ActionHelp}, nil
	case "clear", "tps":
		return MenuInput{Action: ActionClear}, nil
	case "c", "config", "yaml":
		return MenuInput{Action: ActionConfig}, nil
	case "p", "projects":
		return MenuInput{Action: ActionSwitchProject}, nil
	}

	if strings.HasPrefix(low, "p ") {
		name := strings.TrimSpace(line[2:])
		return MenuInput{Action: ActionSwitchProject, Text: name}, nil
	}

	if fields := strings.Fields(line); len(fields) >= 1 {
		if strings.EqualFold(fields[0], "tc") || (len(fields) >= 2 && strings.EqualFold(fields[0], "terminal") && strings.EqualFold(fields[1], "check")) {
			return parseTerminalCheck(fields)
		}
		if isInstallHelpToken(fields[0]) {
			copyClip := false
			for _, f := range fields[1:] {
				if f == "-copy" || f == "--copy" {
					copyClip = true
				}
			}
			return MenuInput{Action: ActionInstallHelp, InstallCopy: copyClip}, nil
		}
		if strings.EqualFold(fields[0], "u") || strings.EqualFold(fields[0], "update") {
			if !canSelfUpdate(p.Path) {
				return MenuInput{}, fmt.Errorf("update only works in a repo with cmd/tps")
			}
			all := false
			for _, f := range fields[1:] {
				if f == "-a" || f == "--all" {
					all = true
				}
			}
			return MenuInput{Action: ActionUpdate, UpdateAll: all}, nil
		}
		if isNewCommandToken(fields[0]) || (len(fields) >= 2 && isNewCommandPhrase(fields[0], fields[1])) {
			flagFields := fields[1:]
			if len(fields) >= 2 && isNewCommandPhrase(fields[0], fields[1]) {
				flagFields = fields[2:]
			}
			cmd, err := parseNewCommandFlags(flagFields)
			if err != nil {
				return MenuInput{}, err
			}
			return MenuInput{Action: ActionNewCommand, NewCmd: cmd}, nil
		}
	}

	if fields := strings.Fields(line); len(fields) >= 1 && strings.EqualFold(fields[0], "e") {
		if len(fields) == 1 {
			return MenuInput{Action: ActionEdit}, nil
		}
		return MenuInput{Action: ActionEdit, Alias: strings.Join(fields[1:], " ")}, nil
	}
	if strings.HasPrefix(low, "edit ") {
		alias := strings.TrimSpace(line[5:])
		return MenuInput{Action: ActionEdit, Alias: alias}, nil
	}
	if low == "edit" || low == "e" {
		return MenuInput{Action: ActionEdit}, nil
	}

	if c := p.FindCommandByAbbrev(line); c != nil {
		return MenuInput{Action: ActionRun, Alias: c.Alias}, nil
	}

	top := p.TopCommands(config.TopCommandSlots)
	if n, err := strconv.Atoi(line); err == nil {
		if n >= 1 && n <= len(top) {
			return MenuInput{Action: ActionRun, Alias: top[n-1].Alias}, nil
		}
		return MenuInput{}, fmt.Errorf("unknown input (h for help, or use a command abbrev)")
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
