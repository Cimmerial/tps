package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

type shellMenuPick struct {
	editKey     string
	addPlatform config.ShellPlatform
}

type shellMenuGroup struct {
	title string
	start int
	end   int
}

func editShellSubmenu(reader *bufio.Reader, projectID, alias string) error {
	for {
		ui.ClearScreen()
		_, p, err := loadProject(projectID)
		if err != nil {
			return err
		}
		cmd := resolveProjectCommand(p, alias)
		if cmd == nil {
			return fmt.Errorf("no command %q", alias)
		}

		fmt.Println()
		fmt.Printf("%s\n", ui.HeaderChoice("edit"))
		fmt.Printf("%s\n", ui.SubHeaderChoice(alias+" — shell"))
		fmt.Println()

		rows, picks, groups := editShellMenuBuild(*cmd)
		for i, g := range groups {
			if g.end <= g.start {
				continue
			}
			fmt.Printf("%s\n", ui.SubHeaderChoice(g.title))
			printChoiceEditMenuPanel("", []menuPanelSection{{Rows: rows[g.start:g.end]}}, nil)
			if i < len(groups)-1 {
				fmt.Println()
			}
		}
		fmt.Println()
		printChoiceEditMenuPanel("", nil, []numberedMenuRow{
			{Num: "0", Desc: ui.Dim("back (q)")},
		})
		fmt.Println()

		choice, back, err := readOptionsMenuChoice()
		if err != nil {
			return editPromptErr(err)
		}
		if back {
			return nil
		}
		if choice < 1 || choice > len(picks) {
			fmt.Println(ui.Dim("invalid choice"))
			continue
		}
		pick := picks[choice-1]
		if pick.editKey != "" {
			if err := editPromptShellKey(reader, projectID, alias, pick.editKey); err != nil {
				return err
			}
			continue
		}
		if err := editPromptAddShell(reader, projectID, alias, pick.addPlatform); err != nil {
			return err
		}
	}
}

func editShellMenuBuild(cmd config.Command) ([]numberedMenuRow, []shellMenuPick, []shellMenuGroup) {
	all := cmd.AllShellFields()
	platforms := []struct {
		title    string
		platform config.ShellPlatform
	}{
		{"all platforms", config.ShellPlatformGeneric},
		{"macos", config.ShellPlatformDarwin},
		{"linux", config.ShellPlatformLinux},
		{"windows", config.ShellPlatformWindows},
	}
	var rows []numberedMenuRow
	var picks []shellMenuPick
	var groups []shellMenuGroup
	n := 0
	for _, g := range platforms {
		start := len(rows)
		for _, f := range all {
			if config.ShellFieldPlatform(f.Key) != g.platform {
				continue
			}
			n++
			rows = append(rows, numberedMenuRow{
				Num:      fmt.Sprintf("%d", n),
				Status:   f.Key,
				Desc:     "edit",
				HelpLine: editMenuHelpLine(editTruncateDisplay(f.Value, 72)),
			})
			picks = append(picks, shellMenuPick{editKey: f.Key})
		}
		n++
		rows = append(rows, numberedMenuRow{
			Num:    fmt.Sprintf("%d", n),
			Status: "add",
			Desc:   "new " + editShellPlatformLabel(g.platform) + " line",
		})
		picks = append(picks, shellMenuPick{addPlatform: g.platform})
		groups = append(groups, shellMenuGroup{title: g.title, start: start, end: len(rows)})
	}
	return rows, picks, groups
}

func editShellPlatformLabel(p config.ShellPlatform) string {
	switch p {
	case config.ShellPlatformDarwin:
		return "macos"
	case config.ShellPlatformLinux:
		return "linux"
	case config.ShellPlatformWindows:
		return "windows"
	default:
		return "generic"
	}
}

func editPromptShellKey(reader *bufio.Reader, projectID, alias, key string) error {
	_, p, err := loadProject(projectID)
	if err != nil {
		return err
	}
	cmd := resolveProjectCommand(p, alias)
	if cmd == nil {
		return fmt.Errorf("no command %q", alias)
	}
	cur := ""
	for _, f := range cmd.AllShellFields() {
		if f.Key == key {
			cur = f.Value
			break
		}
	}
	printEditFieldPromptPanel(key, "click Enter to keep current option, input \"-\" to remove this shell line", cur)
	fmt.Println()
	line, err := readEditNewLine(false)
	if err != nil {
		return editPromptErr(err)
	}
	if line == "" {
		fmt.Println(ui.Dim("unchanged"))
		return nil
	}
	line = strings.TrimSpace(line)
	if line == "-" {
		if err := editMutateCommand(projectID, alias, func(c *config.Command) {
			_ = c.RemoveShellByKey(key)
		}); err != nil {
			return err
		}
		fmt.Println(ui.Green("removed"))
		return nil
	}
	if err := editMutateCommand(projectID, alias, func(c *config.Command) {
		_ = c.SetShellByKey(key, line)
	}); err != nil {
		return err
	}
	fmt.Println(ui.Green("shell updated"))
	return nil
}

func editPromptAddShell(reader *bufio.Reader, projectID, alias string, platform config.ShellPlatform) error {
	_, p, err := loadProject(projectID)
	if err != nil {
		return err
	}
	cmd := resolveProjectCommand(p, alias)
	if cmd == nil {
		return fmt.Errorf("no command %q", alias)
	}
	key := cmd.NextShellKey(platform)
	printEditFieldPromptPanel(
		"new "+editShellPlatformLabel(platform)+" shell line",
		"click Enter to cancel",
		"",
	)
	fmt.Println(ui.Dim("key: "+key))
	fmt.Println()
	line, err := readEditNewLine(false)
	if err != nil {
		return editPromptErr(err)
	}
	line = strings.TrimSpace(line)
	if line == "" {
		fmt.Println(ui.Dim("unchanged"))
		return nil
	}
	if err := editMutateCommand(projectID, alias, func(c *config.Command) {
		_ = c.SetShellByKey(key, line)
	}); err != nil {
		return err
	}
	fmt.Println(ui.Green("shell added"))
	return nil
}
