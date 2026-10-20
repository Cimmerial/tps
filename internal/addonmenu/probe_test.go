package addonmenu

import "testing"

func TestParseMagickFormatLine(t *testing.T) {
	name, mode, ok := parseMagickFormatLine("     PNG* PNG       rw-   Portable Network Graphics")
	if !ok || name != "PNG*" || mode != "rw-" {
		t.Fatalf("got name=%q mode=%q ok=%v", name, mode, ok)
	}
}

func TestMagickFormatNameMatch(t *testing.T) {
	if !formatNameMatches("JPG*", []string{"JPEG", "JPG"}) {
		t.Fatal("jpg")
	}
}
