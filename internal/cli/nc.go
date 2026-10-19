package cli

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func isNewCommandToken(s string) bool {
	return strings.EqualFold(s, "nc")
}

func isNewCommandPhrase(a, b string) bool {
	return strings.EqualFold(a, "new") && strings.EqualFold(b, "command")
}

func parseNewCommandFlags(fields []string) (config.Command, error) {
	var name, abbrev, shell, capture string
	detach := false
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "-t", "--terminal":
			detach = true
		case "-capture", "-m":
			i++
			if i >= len(fields) {
				return config.Command{}, fmt.Errorf("nc: -capture requires a pattern or preset (e.g. trycloudflare)")
			}
			capture = fields[i]
		case "-n":
			i++
			if i >= len(fields) {
				return config.Command{}, fmt.Errorf("nc: -n requires a name")
			}
			name = fields[i]
		case "-a":
			i++
			if i >= len(fields) {
				return config.Command{}, fmt.Errorf("nc: -a requires an abbrev")
			}
			abbrev = fields[i]
		case "-c":
			if i+1 >= len(fields) {
				return config.Command{}, fmt.Errorf("nc: -c requires a shell command")
			}
			shell = strings.Join(fields[i+1:], " ")
			i = len(fields)
		default:
			return config.Command{}, fmt.Errorf("nc: unknown flag %q", fields[i])
		}
	}
	if name == "" {
		return config.Command{}, fmt.Errorf("nc: -n name is required")
	}
	if shell == "" {
		return config.Command{}, fmt.Errorf("nc: -c command is required")
	}
	if capture != "" {
		detach = true
	}
	return config.Command{
		Alias:       name,
		Abbrev:      abbrev,
		Shell:       shell,
		Description: name,
		Detach:      detach,
		Capture:     capture,
	}, nil
}

func addProjectCommand(cfg *config.Config, projectID string, cmd config.Command) error {
	p := cfg.FindByID(projectID)
	if p == nil {
		return fmt.Errorf("project not found")
	}
	if err := p.AddCommand(cmd); err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Println(ui.Green(fmt.Sprintf("added command %q", cmd.Alias)))
	return nil
}
