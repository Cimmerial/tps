package cli

import (
	"os"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/project"
	"github.com/spf13/cobra"
)

var installHelpCopy bool

var installHelpCmd = &cobra.Command{
	Use:     "install-help",
	Short:   "How to install and update tps on this or another machine",
	Aliases: []string{"ih"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if installHelpCopy {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			var selected *config.Project
			cwd, _ := os.Getwd()
			if res, err := project.ResolveFromCwd(cfg, cwd); err == nil {
				selected = res.Project
			}
			return copyInstallCommands(cfg, selected)
		}
		printInstallHelp()
		return nil
	},
}

func init() {
	installHelpCmd.Flags().BoolVar(&installHelpCopy, "copy", false, "copy install commands to the clipboard")
	rootCmd.AddCommand(installHelpCmd)
}
