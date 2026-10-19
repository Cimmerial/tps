package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	gittool "github.com/Cimmerial/tps/internal/git"
	"github.com/Cimmerial/tps/internal/jobs"
	"github.com/Cimmerial/tps/internal/ui"
	"github.com/Cimmerial/tps/internal/update"
)

func runInteractive() error {
	restartBanner := takeSessionRestartBanner()
	restartProjectID := takeSessionRestartProjectID()

	spin := ui.NewSpinner("loading")
	if err := applyUIPreferences(); err != nil {
		spin.Stop()
		return err
	}
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
	choice, err := resolveInteractiveProject(cfg, cwd, restartProjectID)
	spin.Stop()
	if err != nil {
		return err
	}
	if choice.NoProjects {
		fmt.Println("no projects linked. Run: tps project add")
		return nil
	}
	if choice.SaveConfig {
		if err := config.Save(cfg); err != nil {
			return err
		}
	}
	projectID = choice.ProjectID
	if choice.Pick {
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
	if err := config.SaveSessionProjectID(projectID); err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)

	cfg, selected, err := loadProject(projectID)
	if err != nil {
		return err
	}
	if err := ensureDevProjectCommands(cfg, selected); err != nil {
		return err
	}
	if restartBanner == "" {
		if root := devRepoRoot(cfg, selected); root != "" {
			spin := ui.NewSpinner("checking for updates")
			if err := update.MaybeAutoUpdate(root); err != nil {
				spin.Stop()
				fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("auto-update: "+err.Error()))
			} else {
				spin.Stop()
			}
		}
	}
	if len(selected.Commands) == 0 {
		fmt.Println(ui.Dim("no project commands — tps project cmd add <alias> <shell>"))
	}
	initDetachedCompleteNotices()

	statusTick := 0
	ui.ClearScreen()
	ui.PrintLogo()
	printSessionHeader(selected, projectID, restartBanner, "", statusTick)

	gapBeforePrompt := false
	for {
		active, _ := jobs.Reap(projectID)
		line, err := readLine(reader, projectID, statusTick, gapBeforePrompt)
		gapBeforePrompt = false
		if len(active) > 0 {
			statusTick++
		}
		if err != nil {
			return err
		}

		cfg, selected, err = loadProject(projectID)
		if err != nil {
			return err
		}

		input, err := resolveMenuInput(cfg, selected, line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			gapBeforePrompt = true
			continue
		}

		switch input.Action {
		case ActionNoop:
			continue
		case ActionQuit:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			return nil
		case ActionHelp:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			printFullHelp(selected)
		case ActionInstallHelp:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			if input.InstallCopy {
				if err := copyInstallCommands(cfg, selected); err != nil {
					fmt.Fprintf(os.Stderr, "ih: %v\n", err)
				}
			} else {
				printInstallHelp()
			}
		case ActionClear:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			if err := refreshInteractiveSession(projectID, statusTick, "", input.ClearScrollback); err != nil {
				return err
			}
		case ActionConfig:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			if err := config.OpenConfigInEditor(); err != nil {
				fmt.Fprintf(os.Stderr, "config: %v\n", err)
			} else {
				fmt.Println(ui.Dim("opened config.yaml (tps stays open)"))
			}
		case ActionNewCommand:
			recordMetaUse(cfg, selected, input.Action)
			if err := addProjectCommand(cfg, projectID, input.NewCmd); err != nil {
				fmt.Fprintf(os.Stderr, "nc: %v\n", err)
				continue
			}
			cfg, selected, err = loadProject(projectID)
			if err != nil {
				return err
			}
		case ActionUpdate:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			root := devRepoRoot(cfg, selected)
			if root == "" {
				fmt.Fprintf(os.Stderr, "update: link the tps dev repo first (tps project add from the tps repo)\n")
				continue
			}
			if err := update.InstallAndRestart(root, input.UpdateAll); err != nil {
				fmt.Fprintf(os.Stderr, "update: %v\n", err)
			}
		case ActionEdit:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			if err := runEditCommand(reader, cfg, projectID, statusTick, input.Alias); err != nil {
				fmt.Fprintf(os.Stderr, "edit: %v\n", err)
			}
		case ActionSwitchProject:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			newID, err := switchProject(reader, cfg, projectID, input.Text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "project: %v\n", err)
				continue
			}
			if newID != "" {
				prevProjectName := selected.Name
				projectID = newID
				if err := config.SaveSessionProjectID(projectID); err != nil {
					return err
				}
				if err := refreshInteractiveSession(projectID, statusTick, prevProjectName, false); err != nil {
					return err
				}
			}
		case ActionGitShip:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			code, err := runGitShip(selected.Path, input.Text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "git: %v\n", err)
				continue
			}
			if code != 0 {
				fmt.Fprintf(os.Stderr, "git exited with code %d\n", code)
			}
		case ActionGit:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			args := gittool.ParseArgsLine(input.Text)
			code, err := gittool.Run(selected.Path, args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "git: %v\n", err)
				continue
			}
			if code != 0 {
				fmt.Fprintf(os.Stderr, "git exited with code %d\n", code)
			}
		case ActionTerminalCheck:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			if err := runTerminalCheck(selected, input); err != nil {
				fmt.Fprintf(os.Stderr, "t: %v\n", err)
			}
		case ActionCloseAllTerms:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			result, err := runCloseAllTerms()
			if err != nil {
				fmt.Fprintf(os.Stderr, "vite: %v\n", err)
			} else {
				printCloseAllTermsDone(result)
			}
		case ActionExport:
			recordMetaUse(cfg, selected, input.Action)
			if err := runExport(cfg, selected, input); err != nil {
				fmt.Fprintf(os.Stderr, "export: %v\n", err)
			}
		case ActionImport:
			recordMetaUse(cfg, selected, input.Action)
			if err := runImport(cfg, input.Text); err != nil {
				fmt.Fprintf(os.Stderr, "import: %v\n", err)
			}
			cfg, selected, err = loadProject(projectID)
			if err != nil {
				return err
			}
		case ActionSetup:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			if err := runSetup(input.SetupName, input.SetupPath); err != nil {
				fmt.Fprintf(os.Stderr, "setup: %v\n", err)
			}
		case ActionConfigExamples:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			printConfigExamples()
		case ActionOptions:
			recordMetaUse(cfg, selected, input.Action)
			config.Save(cfg)
			if err := runOptionsMenu(projectID, statusTick); err != nil {
				fmt.Fprintf(os.Stderr, "options: %v\n", err)
			}
		case ActionRun:
			code, compactPrompt, err := executeProjectCommand(selected, input.Alias, input.ShellInputs)
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
			if strings.EqualFold(input.Alias, "update") && devRepoRoot(cfg, selected) != "" {
				if err := update.RestartInteractive(true); err != nil {
					fmt.Fprintf(os.Stderr, "restart: %v\n", err)
				}
			}
			if compactPrompt {
				gapBeforePrompt = false
				continue
			}
		}
		if needsGapBeforeNextPrompt(input.Action) {
			gapBeforePrompt = true
		}
	}
}

func refreshInteractiveSession(projectID string, statusTick int, prevProjectName string, clearScrollback bool) error {
	if clearScrollback {
		ui.ClearTerminal()
	} else {
		ui.ClearScreen()
	}
	ui.PrintLogo()
	_, selected, err := loadProject(projectID)
	if err != nil {
		return err
	}
	printSessionHeader(selected, projectID, "", prevProjectName, statusTick)
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

func recordMetaUse(cfg *config.Config, p *config.Project, action InputAction) {
	key := metaKeyForAction(action)
	if key == "" {
		return
	}
	p.RecordMetaUse(key)
}

func pickProject(cfg *config.Config) (*config.Project, error) {
	line, err := ui.ReadEditedLine(ui.Prompt())
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
