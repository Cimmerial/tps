package config

import "testing"

func TestNormalizeColorID(t *testing.T) {
	tests := map[string]string{
		"":              "cyan",
		"default":       "cyan",
		"CYAN":          "cyan",
		"blue":          "blue",
		"bright_red":    "bright_red",
		"light_blue":    "bright_blue",
		"grey":          "white",
		"orange":        "orange",
		"purple":        "purple",
		"pink":          "pink",
		"not_a_color":   "cyan",
		"  magenta  ":   "magenta",
	}
	for in, want := range tests {
		if got := NormalizeColorID(in); got != want {
			t.Errorf("NormalizeColorID(%q) = %q want %q", in, got, want)
		}
	}
}

func TestColorANSICodeKnownPresets(t *testing.T) {
	for _, p := range ColorPresets() {
		if ColorANSICode(p.ID) != p.Code {
			t.Fatalf("ColorANSICode(%q) mismatch", p.ID)
		}
	}
}
