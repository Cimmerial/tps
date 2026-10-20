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

func isNCFlagToken(s string) bool {
	switch s {
	case "-precursor", "--precursor", "-tags", "--tags", "-t", "--terminal", "-capture", "-m", "-n", "-a", "-c":
		return true
	default:
		return false
	}
}

func parseNewCommandFlags(fields []string) (config.Command, error) {
	var name, abbrev, shell, capture, tagsRaw, precursor string
	detach := false
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "-precursor", "--precursor":
			i++
			if i >= len(fields) {
				return config.Command{}, fmt.Errorf("nc: -precursor requires a command alias or abbrev")
			}
			precursor = fields[i]
		case "-tags", "--tags":
			i++
			if i >= len(fields) {
				return config.Command{}, fmt.Errorf("nc: -tags requires a comma-separated list")
			}
			tagsRaw = fields[i]
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
			i++
			if i >= len(fields) {
				return config.Command{}, fmt.Errorf("nc: -c requires a shell command")
			}
			start := i
			for i < len(fields) && !isNCFlagToken(fields[i]) {
				i++
			}
			shell = strings.Join(fields[start:i], " ")
			i--
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
	cmd := config.Command{
		Alias:       name,
		Abbrev:      abbrev,
		Shell:       shell,
		Description: name,
		Detach:      detach,
		Capture:     capture,
	}
	if tagsRaw != "" {
		cmd.Tags = config.ParseTagList(tagsRaw)
	}
	cmd.Precursor = precursor
	return cmd, nil
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
