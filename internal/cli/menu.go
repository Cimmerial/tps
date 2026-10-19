package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/jobs"
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
	ActionExport
	ActionImport
	ActionSetup
	ActionConfigExamples
)

type MenuInput struct {
	Action      InputAction
	Alias       string
	ShellInputs []string
	Text        string
	UpdateAll   bool
	ExportAll   bool
	InstallCopy bool
	SetupName   string
	SetupPath   string
	NewCmd      config.Command
	TcKillAll   bool
	TcKillID    int
	TcRestart   bool
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

func printSessionHeader(p *config.Project, projectID string, restartBanner string, statusTick int) {
	name := ui.Bold(ui.Magenta(p.Name))
	switch restartBanner {
	case "updated":
		fmt.Printf("%s %s %s\n", ui.Header("project:"), name, ui.Green("- tps updated"))
	case "uptodate":
		fmt.Printf("%s %s %s\n", ui.Header("project:"), name, ui.Yellow("- already up to date"))
	default:
		fmt.Printf("%s %s\n", ui.Header("project:"), name)
	}
	printCompactMenu(p)
	fmt.Println()
	printMenuHints(p)
	printTerminalStatusLine(projectID, statusTick)
}

func printTerminalStatusLine(projectID string, tick int) {
	active, err := jobs.Reap(projectID)
	if err != nil || len(active) == 0 {
		return
	}
	line := ui.TerminalsRunningLine(len(active), tick)
	if line != "" {
		fmt.Println(line)
	}
}

func printCompactMenu(p *config.Project) {
	printProjectCommandsGrouped(p, "PRJCT CMDS:")
}

func printFullHelp(p *config.Project) {
	sections := []helpSection{
		{
			title: "session",
			rows: []helpRow{
				cmdRow("clear", "tps", "reset screen and menu"),
				cmdRow("quit", "q", "exit tps"),
				cmdRow("project", "p", "switch linked project"),
			},
		},
		{
			title: "config",
			rows: []helpRow{
				cmdRow("config", "c", "open config.yaml in Cursor/VS Code or $EDITOR"),
				cmdRow("edit", "e", "change a command shell (e hello)"),
				cmdRow("new command", "nc", "nc -n name -a abbrev -c shell (-tags admin,git -t terminal, -m capture); run: ic -user"),
				cmdRow("terminals", "t", "t · tk · t -k 1 · tr"),
				cmdRow("export", "export", "export project YAML (export -a for all)"),
				cmdRow("import", "import", "import path/to/file.tps.yaml"),
				cmdRow("config examples", "ce", "YAML & nc examples for new commands"),
				cmdRow("setup", "setup", "install tps + link this repo (setup -n name --path dir)"),
				tipRow("tags", "tags: [admin, git] — first tag is the menu folder"),
				tipRow("precursor", "precursor: fas — run another alias first; foreground steps share one shell"),
				tipRow("multi shell", "shell1/shell2… or linux_shell / darwin_shell / windows_shell (parallel terminals)"),
			},
		},
	}

	cfg, _ := config.Load()
	if devRepoRoot(cfg, p) != "" {
		sections = append(sections, helpSection{
			title: "update",
			rows: []helpRow{
				cmdRow("update", "u", "go install ./cmd/tps and restart tps"),
				cmdRow("update all", "u -a", "install and copy to every tps on PATH and ./tps"),
			},
		})
	}

	sections = append(sections, helpSection{
		title: "git",
		rows: []helpRow{
			cmdRow("git ship", "g", "g your message or g ship … (git ship …; gs legacy); g status, g diff, …"),
			cmdRow("status", "", "git status"),
			cmdRow("git", "", "git status; or git diff, git log -1, …"),
		},
	})

	sections = append(sections, helpSection{
		title: "help",
		rows: []helpRow{
			cmdRow("help", "h", "this list"),
			cmdRow("install", "ih", "install & update tps on this or another machine"),
			cmdRow("install copy", "ih -copy", "copy install commands to clipboard"),
		},
	})

	printTPSCommandHelp(sections)
	fmt.Println()
	printHelpProjectCommands(p)
}

type helpRow struct {
	command     string
	abbrev      string
	description string
	tip         bool // config.yaml modifier, not a runnable tps command
}

func cmdRow(command, abbrev, description string) helpRow {
	return helpRow{command: command, abbrev: abbrev, description: description}
}

func tipRow(name, description string) helpRow {
	return helpRow{command: name, description: description, tip: true}
}

type helpSection struct {
	title string
	rows  []helpRow
}

func printTPSCommandHelp(sections []helpSection) {
	maxCmd, maxAbbrev := 0, 0
	for _, sec := range sections {
		for _, r := range sec.rows {
			if r.tip {
				continue
			}
			if len(r.command) > maxCmd {
				maxCmd = len(r.command)
			}
			if len(r.abbrev) > maxAbbrev {
				maxAbbrev = len(r.abbrev)
			}
		}
	}

	fmt.Println()
	fmt.Printf("%s\n", ui.Header("TPS commands"))
	for i, sec := range sections {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n", ui.SubHeader(sec.title))
		var printedTips bool
		var hadCommands bool
		for _, r := range sec.rows {
			if r.tip {
				if !printedTips {
					if hadCommands {
						fmt.Println()
					}
					fmt.Printf("    %s\n", ui.Dim("tips (config.yaml — not tps commands)"))
					printedTips = true
				}
				fmt.Printf("    %s %s: %s\n", ui.Dim("·"), ui.Dim(r.command), ui.Dim(r.description))
				continue
			}
			cmdCol := ui.Bold(padPlain(r.command, maxCmd))
			abbrPlain := padPlain(r.abbrev, maxAbbrev)
			abbrCol := abbrPlain
			if strings.TrimSpace(r.abbrev) != "" {
				abbrCol = ui.Abbrev(padPlain(r.abbrev, maxAbbrev))
			}
			fmt.Printf("    %s | %s | %s\n", cmdCol, abbrCol, ui.Dim(r.description))
			hadCommands = true
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

func resolveMenuInput(cfg *config.Config, p *config.Project, line string) (MenuInput, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return MenuInput{Action: ActionNoop}, nil
	}
	low := strings.ToLower(line)
	switch low {
	case "q", "quit", "exit":
		return MenuInput{Action: ActionQuit}, nil
	}
	if ship, ok, err := parseGitShipInput(fieldsFromLine(line)); ok {
		if err != nil {
			return MenuInput{}, err
		}
		return MenuInput{Action: ActionGitShip, Text: ship}, nil
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
		switch strings.ToLower(fields[0]) {
		case "tk", "tck":
			return MenuInput{Action: ActionTerminalCheck, TcKillAll: true}, nil
		case "tr", "tcr":
			return MenuInput{Action: ActionTerminalCheck, TcRestart: true}, nil
		}
		if strings.EqualFold(fields[0], "export") {
			return parseExportInput(fields)
		}
		if strings.EqualFold(fields[0], "import") {
			return parseImportInput(fields)
		}
		if strings.EqualFold(fields[0], "setup") {
			return parseSetupInput(fields)
		}
		if strings.EqualFold(fields[0], "ce") || (len(fields) >= 2 && strings.EqualFold(fields[0], "config") && strings.EqualFold(fields[1], "examples")) {
			return MenuInput{Action: ActionConfigExamples}, nil
		}
		if strings.EqualFold(fields[0], "g") {
			if len(fields) == 1 {
				return MenuInput{}, fmt.Errorf("g needs a git command (example: g status)")
			}
			if strings.EqualFold(fields[1], "ship") {
				return MenuInput{}, fmt.Errorf("g ship needs a message (example: g ship fix the menu)")
			}
			if isKnownGitSubcommand(fields[1]) {
				return MenuInput{Action: ActionGit, Text: strings.Join(fields[1:], " ")}, nil
			}
		}
		if strings.EqualFold(fields[0], "t") || strings.EqualFold(fields[0], "tc") || (len(fields) >= 2 && strings.EqualFold(fields[0], "terminal") && strings.EqualFold(fields[1], "check")) {
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
			if devRepoRoot(cfg, p) == "" {
				return MenuInput{}, fmt.Errorf("update needs the tps dev repo linked (tps project add from the tps repo)")
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

	if fields := strings.Fields(line); len(fields) >= 1 {
		if cmd := p.ResolveCommandRef(fields[0]); cmd != nil {
			inputs, err := config.ParseShellInputTokens(fields[1:])
			if err != nil {
				return MenuInput{}, err
			}
			if err := config.ValidateCommandInputs(cmd, inputs); err != nil {
				return MenuInput{}, err
			}
			return MenuInput{Action: ActionRun, Alias: cmd.Alias, ShellInputs: inputs}, nil
		}
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
