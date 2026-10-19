package cli

import (
	"github.com/Cimmerial/tps/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage tps configuration",
}

var configOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open config.yaml in VS Code or $EDITOR",
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.OpenConfigInEditor()
	},
}

func init() {
	configCmd.AddCommand(configOpenCmd)
	rootCmd.AddCommand(configCmd)
}
