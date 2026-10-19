package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func runEditCommand(reader *bufio.Reader, cfg *config.Config, projectID string, alias string) error {
	p := cfg.FindByID(projectID)
	if p == nil {
		return fmt.Errorf("project not found")
	}
	if len(p.Commands) == 0 {
		return fmt.Errorf("no commands to edit (tps project cmd add <alias> <shell>)")
	}

	if alias == "" {
		fmt.Println(ui.Dim("pick alias to edit (or: e hello)"))
		for _, c := range p.Commands {
			line := c.Alias
			if c.Abbrev != "" {
				line = fmt.Sprintf("%s (%s)", c.Alias, c.Abbrev)
			}
			fmt.Printf("  %s\n", line)
		}
		fmt.Print(ui.Prompt())
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		alias = strings.TrimSpace(line)
		if alias == "" {
			return nil
		}
	}

	cmd := resolveProjectCommand(p, alias)
	if cmd == nil {
		return fmt.Errorf("no command %q", alias)
	}
	alias = cmd.Alias

	fmt.Printf("shell [%s]: ", cmd.Shell)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	shell := strings.TrimSpace(line)
	if shell == "" {
		fmt.Println(ui.Dim("unchanged"))
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
	fmt.Println(ui.Green(fmt.Sprintf("updated %q", alias)))
	return config.Save(cfg)
}

func resolveProjectCommand(p *config.Project, name string) *config.Command {
	if c := p.FindCommand(name); c != nil {
		return c
	}
	if c := p.FindCommandByAbbrev(name); c != nil {
		return c
	}
	for i := range p.Commands {
		if strings.EqualFold(p.Commands[i].Alias, name) {
			return &p.Commands[i]
		}
	}
	return nil
}
