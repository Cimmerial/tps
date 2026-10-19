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

func commandColumnWidths(cmds []config.Command) (maxCmd, maxAbbrev, maxDesc int) {
	for _, c := range cmds {
		if len(c.Alias) > maxCmd {
			maxCmd = len(c.Alias)
		}
		if len(c.Abbrev) > maxAbbrev {
			maxAbbrev = len(c.Abbrev)
		}
		desc := commandDescription(c)
		if len(desc) > maxDesc {
			maxDesc = len(desc)
		}
	}
	return maxCmd, maxAbbrev, maxDesc
}

func formatCommandTableRow(c config.Command, maxCmd, maxAbbrev int) string {
	cmdCol := ui.Bold(padPlain(c.Alias, maxCmd))
	abbrPlain := padPlain(c.Abbrev, maxAbbrev)
	abbrCol := abbrPlain
	if strings.TrimSpace(c.Abbrev) != "" {
		abbrCol = ui.Abbrev(padPlain(c.Abbrev, maxAbbrev))
	}
	return fmt.Sprintf("%s | %s | %s", cmdCol, abbrCol, ui.Dim(commandDescription(c)))
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
	all := make([]config.Command, 0, len(p.Commands))
	for _, g := range groups {
		all = append(all, g.Commands...)
	}
	maxCmd, maxAbbrev, _ := commandColumnWidths(all)

	fmt.Printf("%s\n", ui.Header("project commands"))
	for i, g := range groups {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n", ui.SubHeader("tag: "+tagFolderLabel(g.Tag)))
		for _, c := range g.Commands {
			fmt.Printf("    %s\n", formatCommandTableRow(c, maxCmd, maxAbbrev))
		}
	}
}

func printProjectCommandsGrouped(p *config.Project, title string) {
	groups := p.CommandsByTag()
	if len(groups) == 0 {
		return
	}
	all := make([]config.Command, 0, len(p.Commands))
	for _, g := range groups {
		all = append(all, g.Commands...)
	}
	maxCmd, maxAbbrev, _ := commandColumnWidths(all)

	fmt.Printf("%s\n", ui.Header(title))
	for i, g := range groups {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n", ui.SubHeader(tagFolderLabel(g.Tag)))
		for _, c := range g.Commands {
			fmt.Printf("    %s\n", formatCommandTableRow(c, maxCmd, maxAbbrev))
		}
	}
}
