package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

func runEditCommand(reader *bufio.Reader, cfg *config.Config, projectID string, alias string) error {
	p := cfg.FindByID(projectID)
	if p == nil {
		return fmt.Errorf("project not found")
	}
	if len(p.Commands) == 0 {
		return fmt.Errorf("no commands to edit")
	}

	if alias == "" {
		printFullHelp(p)
		fmt.Print("edit alias> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		alias = strings.TrimSpace(line)
		if alias == "" {
			return nil
		}
	}

	cmd := p.FindCommand(alias)
	if cmd == nil {
		for _, c := range p.Commands {
			if strings.EqualFold(c.Alias, alias) {
				cmd = &c
				alias = c.Alias
				break
			}
		}
	}
	if cmd == nil {
		return fmt.Errorf("no command %q", alias)
	}

	fmt.Printf("shell [%s]: ", cmd.Shell)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	shell := strings.TrimSpace(line)
	if shell == "" {
		fmt.Println("unchanged")
		return nil
	}

	updated := config.Command{
		Alias:       alias,
		Abbrev:      cmd.Abbrev,
		Shell:       shell,
		Description: cmd.Description,
	}
	if err := p.AddCommand(updated); err != nil {
		return err
	}
	return config.Save(cfg)
}
