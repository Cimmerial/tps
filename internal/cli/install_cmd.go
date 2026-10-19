package cli

import (
	"github.com/spf13/cobra"
)

var installHelpCmd = &cobra.Command{
	Use:     "install-help",
	Short:   "How to install and update tps on this or another machine",
	Aliases: []string{"ih"},
	RunE: func(cmd *cobra.Command, args []string) error {
		printInstallHelp()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installHelpCmd)
}
