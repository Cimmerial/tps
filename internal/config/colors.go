package config

import "strings"

// ColorPreset is a named terminal foreground color (ANSI 16 or common 256).
type ColorPreset struct {
	ID    string
	Label string
	Code  string // SGR foreground; reset is applied by ui.paint
}

// ColorPresets lists selectable title-card and accent colors (menu order).
func ColorPresets() []ColorPreset {
	return colorPresets
}

// ColorPresetIDs returns normalized ids in menu order.
func ColorPresetIDs() []string {
	out := make([]string, len(colorPresets))
	for i, p := range colorPresets {
		out[i] = p.ID
	}
	return out
}

var colorPresets = []ColorPreset{
	{ID: "cyan", Label: "cyan (default)", Code: "\033[36m"},
	{ID: "blue", Label: "blue", Code: "\033[34m"},
	{ID: "magenta", Label: "magenta", Code: "\033[35m"},
	{ID: "red", Label: "red", Code: "\033[31m"},
	{ID: "green", Label: "green", Code: "\033[32m"},
	{ID: "yellow", Label: "yellow", Code: "\033[33m"},
	{ID: "white", Label: "white", Code: "\033[37m"},
	{ID: "bright_cyan", Label: "bright cyan", Code: "\033[96m"},
	{ID: "bright_blue", Label: "bright blue", Code: "\033[94m"},
	{ID: "bright_magenta", Label: "bright magenta", Code: "\033[95m"},
	{ID: "bright_red", Label: "bright red", Code: "\033[91m"},
	{ID: "bright_green", Label: "bright green", Code: "\033[92m"},
	{ID: "bright_yellow", Label: "bright yellow", Code: "\033[93m"},
	{ID: "bright_white", Label: "bright white", Code: "\033[97m"},
	{ID: "orange", Label: "orange", Code: "\033[38;5;208m"},
	{ID: "purple", Label: "purple", Code: "\033[38;5;135m"},
	{ID: "pink", Label: "pink", Code: "\033[38;5;213m"},
	{ID: "lime", Label: "lime", Code: "\033[92m"},
	{ID: "teal", Label: "teal", Code: "\033[38;5;37m"},
}

var colorAlias = map[string]string{
	"default":     "cyan",
	"light_cyan":  "bright_cyan",
	"light_blue":  "bright_blue",
	"light_green": "bright_green",
	"light_red":   "bright_red",
	"gray":        "white",
	"grey":        "white",
}

// NormalizeColorID returns a canonical preset id (unknown values default to cyan).
func NormalizeColorID(raw string) string {
	s := strings.ToLower(strings.TrimSpace(raw))
	if s == "" {
		return "cyan"
	}
	if canon, ok := colorAlias[s]; ok {
		s = canon
	}
	for _, p := range colorPresets {
		if p.ID == s {
			return s
		}
	}
	return "cyan"
}

// ColorANSICode returns the SGR sequence for a normalized color id.
func ColorANSICode(id string) string {
	id = NormalizeColorID(id)
	for _, p := range colorPresets {
		if p.ID == id {
			return p.Code
		}
	}
	return "\033[36m"
}
