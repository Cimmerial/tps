package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/project"
	"github.com/Cimmerial/tps/internal/ui"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
}

var projectAddPath string

var projectAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Link a project (default: current directory)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		spin := ui.NewSpinner("linking project")
		cfg, err := config.Load()
		spin.Stop()
		if err != nil {
			return err
		}

		linkDir, err := os.Getwd()
		if err != nil {
			return err
		}
		if projectAddPath != "" {
			linkDir = projectAddPath
		}
		name := ""
		if len(args) > 0 {
			name = args[0]
		} else {
			root, err := project.ProjectRootFromCwd(linkDir)
			if err != nil {
				return err
			}
			name = filepath.Base(root)
		}

		spin = ui.NewSpinner("linking project")
		linked, err := linkProjectAt(cfg, linkDir, name)
		spin.Stop()
		if err != nil {
			return err
		}
		fmt.Printf("linked %q at %s", linked.Name, linked.Path)
		if linked.GitRemote != "" {
			fmt.Printf(" (%s)", linked.GitRemote)
		}
		fmt.Println()
		return nil
	},
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List linked projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if len(cfg.Projects) == 0 {
			fmt.Println("no projects linked (try: tps project add)")
			return nil
		}
		for _, p := range cfg.Projects {
			line := fmt.Sprintf("%s  %s", p.Name, p.Path)
			if p.GitRemote != "" {
				line += fmt.Sprintf("  [%s]", p.GitRemote)
			}
			fmt.Println(line)
		}
		return nil
	},
}

var projectCmdCmd = &cobra.Command{
	Use:   "cmd",
	Short: "Manage project commands",
}

var projectCmdEditCmd = &cobra.Command{
	Use:   "edit <alias> <shell>",
	Short: "Change a command alias for the project in the current directory",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return mutateProjectCommand(args[0], config.Command{
			Alias: args[0],
			Shell: strings.Join(args[1:], " "),
		}, "updated")
	},
}

var projectCmdRmCmd = &cobra.Command{
	Use:   "rm <alias>",
	Short: "Remove a command alias from the project in the current directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		spin := ui.NewSpinner("resolving project")
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
		res, err := project.ResolveFromCwd(cfg, cwd)
		spin.Stop()
		if err != nil {
			return err
		}
		if res.Healed {
			if err := config.Save(cfg); err != nil {
				return err
			}
		}
		if !res.Project.RemoveCommand(args[0]) {
			return fmt.Errorf("no command %q", args[0])
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Printf("removed command %q from %q\n", args[0], res.Project.Name)
		return nil
	},
}

var projectCmdAddCmd = &cobra.Command{
	Use:   "add <alias> <shell>",
	Short: "Add a command alias to the project in the current directory",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return mutateProjectCommand(args[0], config.Command{
			Alias: args[0],
			Shell: strings.Join(args[1:], " "),
		}, "added")
	},
}

func mutateProjectCommand(alias string, command config.Command, verb string) error {
	spin := ui.NewSpinner("resolving project")
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
	res, err := project.ResolveFromCwd(cfg, cwd)
	spin.Stop()
	if err != nil {
		return err
	}
	if res.Healed {
		if err := config.Save(cfg); err != nil {
			return err
		}
	}
	p := res.Project
	if existing := p.FindCommand(command.Alias); existing != nil {
		if command.Abbrev == "" {
			command.Abbrev = existing.Abbrev
		}
		if command.Description == "" {
			command.Description = existing.Description
		}
	}
	if err := p.AddCommand(command); err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Printf("%s command %q to %q\n", verb, alias, p.Name)
	return nil
}

func init() {
	projectAddCmd.Flags().StringVar(&projectAddPath, "path", "", "project directory to link (default: current directory)")
	projectCmd.AddCommand(projectAddCmd, projectListCmd, projectCmdCmd)
	projectCmdCmd.AddCommand(projectCmdAddCmd, projectCmdEditCmd, projectCmdRmCmd)
}
