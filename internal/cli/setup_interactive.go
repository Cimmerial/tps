package cli

import (
	"fmt"
	"strings"
)

func parseSetupInput(fields []string) (MenuInput, error) {
	if len(fields) == 0 || !strings.EqualFold(fields[0], "setup") {
		return MenuInput{}, fmt.Errorf("setup: expected setup")
	}
	in := MenuInput{Action: ActionSetup}
	for i := 1; i < len(fields); i++ {
		switch fields[i] {
		case "-n", "--name":
			i++
			if i >= len(fields) {
				return MenuInput{}, fmt.Errorf("setup: -n requires a name")
			}
			in.SetupName = fields[i]
		case "--path":
			i++
			if i >= len(fields) {
				return MenuInput{}, fmt.Errorf("setup: --path requires a directory")
			}
			in.SetupPath = fields[i]
		default:
			return MenuInput{}, fmt.Errorf("setup: unknown flag %q", fields[i])
		}
	}
	return in, nil
}
