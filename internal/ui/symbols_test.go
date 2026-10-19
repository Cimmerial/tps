package ui

import "testing"

func TestSymbolCatalogIncludesCheckAndCross(t *testing.T) {
	var haveCheck, haveCross bool
	for _, e := range SymbolCatalog() {
		if e.ID == "check" && e.Glyph == SymbolCheck {
			haveCheck = true
		}
		if e.ID == "cross" && e.Glyph == SymbolCross {
			haveCross = true
		}
	}
	if !haveCheck || !haveCross {
		t.Fatalf("catalog missing check/cross: check=%v cross=%v", haveCheck, haveCross)
	}
}

func TestLookupSymbol(t *testing.T) {
	if ui := LookupSymbol("check"); ui == nil || ui.Glyph != SymbolCheck {
		t.Fatal("expected check symbol")
	}
	if LookupSymbol("nope") != nil {
		t.Fatal("unknown id should be nil")
	}
}
