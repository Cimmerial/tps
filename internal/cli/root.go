package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command: `tps`
var rootCmd = &cobra.Command{
	Use:   "tps",
	Short: "Project command hub",
	Long:  "Link projects and run short command aliases.",
	// No args → interactive mode (add this in session 4)
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInteractive() // stub: return fmt.Errorf("not implemented")
	},
}

// Execute is called from main.
func Execute() error {
	sessionRestartBanner = peelRestartFlagsFromArgv()
	return rootCmd.Execute()
}

func init() {
	// Global flags later, e.g. --config
	// rootCmd.PersistentFlags().StringVar(&cfgPath, "config", "", "config file")

	rootCmd.AddCommand(projectCmd)
}

// Optional: print errors nicely (cobra already prints Usage in many cases)
func exitWithError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}