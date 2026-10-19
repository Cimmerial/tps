package cli

import (
	"os"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/project"
	"github.com/Cimmerial/tps/internal/ui"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <alias>",
	Short: "Run a command alias for the project in the current directory",
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
		code, _, err := executeProjectCommand(res.Project, args[0], nil)
		if err != nil {
			return err
		}
		if code != 0 {
			os.Exit(code)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
