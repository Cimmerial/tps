package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	gittool "github.com/Cimmerial/tps/internal/git"
	"github.com/Cimmerial/tps/internal/project"
	"github.com/Cimmerial/tps/internal/runner"
	"github.com/Cimmerial/tps/internal/ui"
	"github.com/Cimmerial/tps/internal/update"
)

func runInteractive() error {
	ui.PrintLogo()

	spin := ui.NewSpinner("loading")
	cfg, err := config.Load()
	if err != nil {
		spin.Stop()
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		spin.Stop()
		return err
	}

	var projectID string
	res, resolveErr := project.ResolveFromCwd(cfg, cwd)
	spin.Stop()
	if resolveErr == nil {
		if res.Healed {
			if err := config.Save(cfg); err != nil {
				return err
			}
		}
		projectID = res.Project.ID
	} else {
		if len(cfg.Projects) == 0 {
			fmt.Println("no projects linked. Run: tps project add")
			return nil
		}
		fmt.Println("select a project:")
		for i, p := range cfg.Projects {
			fmt.Printf("  %d) %s  %s\n", i+1, p.Name, p.Path)
		}
		selected, err := pickProject(cfg)
		if err != nil {
			return err
		}
		if selected == nil {
			return nil
		}
		projectID = selected.ID
	}

	reader := bufio.NewReader(os.Stdin)

	cfg, selected, err := loadProject(projectID)
	if err != nil {
		return err
	}
	if len(selected.Commands) == 0 {
		fmt.Println("no commands for this project. Run: tps project cmd add <alias> <shell>")
		return nil
	}
	printSessionHeader(selected)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		cfg, selected, err = loadProject(projectID)
		if err != nil {
			return err
		}

		input, err := resolveMenuInput(selected, line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}

		switch input.Action {
		case ActionNoop:
			continue
		case ActionQuit:
			return nil
		case ActionHelp:
			printFullHelp(selected)
		case ActionClear:
			if err := refreshInteractiveSession(projectID); err != nil {
				return err
			}
		case ActionConfig:
			if err := config.OpenConfigInEditor(); err != nil {
				fmt.Fprintf(os.Stderr, "config: %v\n", err)
			}
		case ActionUpdate:
			if err := update.InstallAndRestart(selected.Path); err != nil {
				fmt.Fprintf(os.Stderr, "update: %v\n", err)
			}
		case ActionEdit:
			if err := runEditCommand(reader, cfg, projectID, input.Alias); err != nil {
				fmt.Fprintf(os.Stderr, "edit: %v\n", err)
			}
		case ActionGitShip:
			code, err := gittool.Ship(selected.Path, input.Text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "git: %v\n", err)
				continue
			}
			if code != 0 {
				fmt.Fprintf(os.Stderr, "git exited with code %d\n", code)
			}
		case ActionGit:
			args := gittool.ParseArgsLine(input.Text)
			code, err := gittool.Run(selected.Path, args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "git: %v\n", err)
				continue
			}
			if code != 0 {
				fmt.Fprintf(os.Stderr, "git exited with code %d\n", code)
			}
		case ActionRun:
			code, err := runner.RunCommand(selected, input.Alias)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				continue
			}
			if code != 0 {
				fmt.Fprintf(os.Stderr, "command exited with code %d\n", code)
				continue
			}
			selected.RecordCommandUse(input.Alias)
			if err := config.Save(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "could not save usage: %v\n", err)
			}
			if strings.EqualFold(input.Alias, "update") && canSelfUpdate(selected.Path) {
				if err := update.RestartInteractive(); err != nil {
					fmt.Fprintf(os.Stderr, "restart: %v\n", err)
				}
			}
		}
	}
}

func refreshInteractiveSession(projectID string) error {
	ui.ClearScreen()
	ui.PrintLogo()
	_, selected, err := loadProject(projectID)
	if err != nil {
		return err
	}
	printSessionHeader(selected)
	return nil
}

func loadProject(projectID string) (*config.Config, *config.Project, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, err
	}
	selected := cfg.FindByID(projectID)
	if selected == nil {
		return nil, nil, fmt.Errorf("project no longer in config")
	}
	return cfg, selected, nil
}

func pickProject(cfg *config.Config) (*config.Project, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	if line == "" || strings.EqualFold(line, "q") {
		return nil, nil
	}
	n, err := strconv.Atoi(line)
	if err != nil || n < 1 || n > len(cfg.Projects) {
		return nil, fmt.Errorf("invalid selection")
	}
	return &cfg.Projects[n-1], nil
}
