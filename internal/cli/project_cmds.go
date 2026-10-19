package cli

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
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

func formatCommandTableRow(c config.Command, maxCmd, maxAbbrev int, hasAbbrev bool) string {
	cmdCol := ui.Bold(padPlain(c.Alias, maxCmd))
	descCol := ui.Dim(commandDescription(c))
	if !hasAbbrev {
		return fmt.Sprintf("%s | %s", cmdCol, descCol)
	}
	abbrPlain := padPlain(c.Abbrev, maxAbbrev)
	abbrCol := abbrPlain
	if strings.TrimSpace(c.Abbrev) != "" {
		abbrCol = ui.Abbrev(padPlain(c.Abbrev, maxAbbrev))
	}
	return fmt.Sprintf("%s | %s | %s", cmdCol, abbrCol, descCol)
}

func tagFolderLabel(tag string) string {
	if strings.TrimSpace(tag) == "" {
		return "general"
	}
	return tag
}

func printHelpProjectCommands(p *config.Project) {
	groups := p.CommandsByTag()
	if len(groups) == 0 {
		return
	}
	fmt.Printf("%s\n", ui.Header("project commands"))
	for i, g := range groups {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n", ui.SubHeader("tag: "+tagFolderLabel(g.Tag)))
		maxCmd, maxAbbrev, _, hasAbbrev := commandColumnWidths(g.Commands)
		for _, c := range g.Commands {
			fmt.Printf("    %s\n", formatCommandTableRow(c, maxCmd, maxAbbrev, hasAbbrev))
		}
	}
}

func printProjectCommandsGrouped(p *config.Project, title string) {
	groups := p.CommandsByTag()
	if len(groups) == 0 {
		return
	}
	fmt.Printf("%s\n", ui.Header(title))
	for i, g := range groups {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n", ui.SubHeader(tagFolderLabel(g.Tag)))
		maxCmd, maxAbbrev, _, hasAbbrev := commandColumnWidths(g.Commands)
		for _, c := range g.Commands {
			fmt.Printf("    %s\n", formatCommandTableRow(c, maxCmd, maxAbbrev, hasAbbrev))
		}
	}
}
