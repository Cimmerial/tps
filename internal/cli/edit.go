package cli

import (
	"bufio"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

func runEditCommand(reader *bufio.Reader, cfg *config.Config, projectID string, statusTick int, alias string) error {
	return runEditMenu(reader, cfg, projectID, statusTick, alias)
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
