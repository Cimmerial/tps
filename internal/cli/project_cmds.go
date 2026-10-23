package cli

import (
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

func commandDescription(c config.Command) string {
	if c.Description != "" {
		return c.Description
	}
	return c.Shell
}

func commandColumnWidths(cmds []config.Command) (maxCmd, maxAbbrev, maxDesc int, hasAbbrev bool) {
	for _, c := range cmds {
		if len(c.Alias) > maxCmd {
			maxCmd = len(c.Alias)
		}
		if strings.TrimSpace(c.Abbrev) != "" {
			hasAbbrev = true
			if len(c.Abbrev) > maxAbbrev {
				maxAbbrev = len(c.Abbrev)
			}
		}
		desc := commandDescription(c)
		if len(desc) > maxDesc {
			maxDesc = len(desc)
		}
	}
	return maxCmd, maxAbbrev, maxDesc, hasAbbrev
}

func tagFolderLabel(tag string) string {
	if strings.TrimSpace(tag) == "" {
		return "general"
	}
	return tag
}

func printHelpProjectCommands(p *config.Project) {
	printProjectCommandsGrouped(p, "project commands")
}

func printProjectCommandsGrouped(p *config.Project, title string) {
	groups := p.CommandsByTag()
	if len(groups) == 0 {
		return
	}
	var sections []menuPanelSection
	for _, g := range groups {
		_, _, _, hasAbbrev := commandColumnWidths(g.Commands)
		var rows []numberedMenuRow
		for _, c := range g.Commands {
			row := numberedMenuRow{
				Desc:     c.Alias,
				HelpLine: commandDescription(c),
			}
			if hasAbbrev {
				row.Abbrev = c.Abbrev
			}
			rows = append(rows, row)
		}
		sections = append(sections, menuPanelSection{
			Label: strings.ToUpper(tagFolderLabel(g.Tag)),
			Rows:  rows,
		})
	}
	printSessionMenuPanel(title, sections)
}
