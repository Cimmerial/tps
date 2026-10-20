package cli

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"ver"},
	Short:   "Print the tps version",
	Run: func(cmd *cobra.Command, args []string) {
		printTPSVersion(nil)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
