package ui

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSymbolDumpIncludesBlockShadeAndBraille(t *testing.T) {
	root, err := repoRootForTest()
	if err != nil {
		t.Fatal(err)
	}
	var glyphs []string
	for _, sec := range SymbolDump(root, "") {
		for _, e := range sec.Entries {
			glyphs = append(glyphs, e.Glyph)
		}
	}
	joined := strings.Join(glyphs, "")
	for _, want := range []string{"░", "▒", "▓", "⣾", "✓", "\\033[2J"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("dump missing %q in combined output", want)
		}
	}
}

func TestSymbolDumpFilterBlock(t *testing.T) {
	root, err := repoRootForTest()
	if err != nil {
		t.Fatal(err)
	}
	sections := SymbolDump(root, "block")
	if len(sections) == 0 {
		t.Fatal("expected sections for block filter")
	}
	found := false
	for _, sec := range sections {
		for _, e := range sec.Entries {
			if strings.Contains(e.Glyph, "▓") || e.ID == "block_dark" {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("block filter should match shade glyphs")
	}
}

func repoRootForTest() (string, error) {
	// module root: internal/ui -> ../..
	return filepath.Clean(filepath.Join("..", "..")), nil
}
