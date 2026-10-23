package ui

import (
	"fmt"
	"strings"
)

// UI glyphs used across menus, spinners, and status lines.
const (
	SymbolCheck      = "✓"
	SymbolCross      = "✗"
	SymbolEmDash     = "—"
	SymbolArrow      = "→"
	SymbolEllipsis   = "…"
	SymbolMiddleDot  = "·"
	SymbolCommandKey = "⌘"
	SymbolPrompt     = ">"
)

// BrailleSpinnerFrames animates busy / waiting states (prompt, spinner, terminal count).
var BrailleSpinnerFrames = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

// SymbolEntry documents one glyph for the sym catalog.
type SymbolEntry struct {
	ID    string // lookup id, e.g. check
	Glyph string
	Usage string
}

// SymbolCatalog returns hand-documented UI symbols (stable order).
func SymbolCatalog() []SymbolEntry {
	frames := strings.Join(BrailleSpinnerFrames, " ")
	return []SymbolEntry{
		{ID: "check", Glyph: SymbolCheck, Usage: "set / on / success (edit menus, capture actions, git ship ok)"},
		{ID: "cross", Glyph: SymbolCross, Usage: "unset / off (edit menus)"},
		{ID: "emdash", Glyph: SymbolEmDash, Usage: "empty field; menu prefixes (— shell line)"},
		{ID: "arrow", Glyph: SymbolArrow, Usage: "project switch (a → b); options → capture; INPUT hints"},
		{ID: "ellipsis", Glyph: SymbolEllipsis, Usage: "loading / waiting / truncated text"},
		{ID: "middledot", Glyph: SymbolMiddleDot, Usage: "hint separator (h · e · q); menu joins; capture msgs"},
		{ID: "command", Glyph: SymbolCommandKey, Usage: "install help (⌘V paste)"},
		{ID: "prompt", Glyph: ">:", Usage: "interactive input prompt (magenta info / red choice menus)"},
		{ID: "braille", Glyph: frames, Usage: "all spinner frames — see braille spinner section for each"},
		{ID: "block_light", Glyph: "░", Usage: "logo shading (blurry, clean presets)"},
		{ID: "block_medium", Glyph: "▒", Usage: "logo shading"},
		{ID: "block_dark", Glyph: "▓", Usage: "logo shading (bar-graph style)"},
	}
}

// LookupSymbol resolves an id (case-insensitive). Empty id returns nil.
func LookupSymbol(id string) *SymbolEntry {
	id = strings.TrimSpace(strings.ToLower(id))
	if id == "" {
		return nil
	}
	for _, e := range SymbolCatalog() {
		if strings.EqualFold(e.ID, id) {
			return &e
		}
	}
	for _, e := range brailleFrameEntries() {
		if strings.EqualFold(e.ID, id) {
			return &e
		}
	}
	if strings.HasPrefix(id, "u+") {
		var cp rune
		if _, err := fmt.Sscanf(id, "u+%x", &cp); err == nil {
			g := string(cp)
			e := SymbolEntry{
				ID:    fmt.Sprintf("U+%04X", cp),
				Glyph: g,
				Usage: runeCategory(cp),
			}
			return &e
		}
	}
	return nil
}
