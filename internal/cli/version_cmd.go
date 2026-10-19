package cli

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the tps version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
