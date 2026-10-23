package cli

import (
	"os"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	gittool "github.com/Cimmerial/tps/internal/git"
	"github.com/Cimmerial/tps/internal/project"
	"github.com/Cimmerial/tps/internal/ui"
	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Run git in the current project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		return runGitInProject(args)
	},
}

var gitShipCmd = &cobra.Command{
	Use:   "ship [message]",
	Short: "choose files, commit, and push",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := projectRootFromCwd()
		if err != nil {
			return err
		}
		msg := strings.Join(args, " ")
		code, err := runGitShip(root, msg)
		if err != nil {
			return err
		}
		if code != 0 {
			os.Exit(code)
		}
		return nil
	},
}

var gitStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "git status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGitInProject(append([]string{"status"}, args...))
	},
}

func runGitInProject(args []string) error {
	root, err := projectRootFromCwd()
	if err != nil {
		return err
	}
	code, err := gittool.Run(root, args)
	if err != nil {
		return err
	}
	if code != 0 {
		os.Exit(code)
	}
	return nil
}

func projectRootFromCwd() (string, error) {
	spin := ui.NewSpinner("resolving project")
	cfg, err := config.Load()
	if err != nil {
		spin.Stop()
		return "", err
	}
	cwd, err := os.Getwd()
	if err != nil {
		spin.Stop()
		return "", err
	}
	res, err := project.ResolveFromCwd(cfg, cwd)
	spin.Stop()
	if err != nil {
		return "", err
	}
	if res.Healed {
		if err := config.Save(cfg); err != nil {
			return "", err
		}
	}
	return res.Project.Path, nil
}

func init() {
	gitCmd.AddCommand(gitShipCmd, gitStatusCmd)
	rootCmd.AddCommand(gitCmd)
}
