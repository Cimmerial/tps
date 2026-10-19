package cli

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/instance"
	"github.com/Cimmerial/tps/internal/ui"
)

func ensureDevProjectCommands(cfg *config.Config, p *config.Project) error {
	if cfg == nil || p == nil {
		return nil
	}
	changed, err := instance.EnsureDefaultDevProjectCommands(p)
	if err != nil {
		return err
	}
	if changed {
		return config.Save(cfg)
	}
	return nil
}

// runDevBuiltin handles shell-less commands on the tps dev repo.
// Returns handled, exit code, error.
func runDevBuiltin(p *config.Project, alias string, inputs []string) (bool, int, error) {
	if p == nil || !instance.IsBuiltinDevCommand(p, alias) {
		return false, 0, nil
	}
	cmd := p.FindCommand(alias)
	handler := strings.ToLower(strings.TrimSpace(cmd.Handler))
	switch handler {
	case "sym":
		printSymbolCatalog(p.Path, inputs)
		return true, 0, nil
	default:
		return true, 1, fmt.Errorf("builtin %q not implemented", alias)
	}
}

func printSymbolCatalog(repoRoot string, inputs []string) {
	filter := ""
	if len(inputs) > 0 {
		filter = strings.TrimSpace(inputs[0])
	}
	fmt.Printf("%s\n", ui.Header("tps symbols"))
	fmt.Println()
	if filter != "" {
		if e := ui.LookupSymbol(filter); e != nil {
			printSymbolLine(*e)
			return
		}
	}
	sections := ui.SymbolDump(repoRoot, filter)
	if len(sections) == 0 {
		fmt.Printf("%s\n", ui.Yellow(fmt.Sprintf("no matches for %q", filter)))
		return
	}
	for i, sec := range sections {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n", ui.SubHeader(sec.Title))
		for _, e := range sec.Entries {
			printSymbolLine(e)
		}
	}
	if filter == "" {
		fmt.Println()
		fmt.Println(ui.Dim("sym <filter>   id, category, or substring (e.g. sym block, sym U+2593, sym braille-0)"))
	}
}

func printSymbolLine(e ui.SymbolEntry) {
	glyph := e.Glyph
	switch e.ID {
	case "check", "block_light", "block_medium", "block_dark":
		if len([]rune(e.Glyph)) == 1 {
			glyph = ui.Green(e.Glyph)
		}
	case "cross":
		glyph = ui.Red(e.Glyph)
	}
	id := e.ID
	if len(id) > 18 {
		id = id[:15] + "…"
	}
	fmt.Printf("  %-18s %s  %s\n", ui.Abbrev(id), glyph, ui.Dim(e.Usage))
}
