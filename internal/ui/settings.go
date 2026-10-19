package ui

import (
	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui/logos"
)

// UISettings holds runtime theme preferences applied from config.
type UISettings struct {
	LogoID          string
	TitleCardColor  string // config.ColorPresets id (logo only)
	Accent          string // config.ColorPresets id
	ColorEnabled    *bool
}

var current UISettings

// ApplySettings updates global UI theme state (logo, accent, color override).
func ApplySettings(s UISettings) {
	if s.LogoID == "" {
		s.LogoID = "default"
	}
	s.Accent = config.NormalizeColorID(s.Accent)
	s.TitleCardColor = config.NormalizeColorID(s.TitleCardColor)
	current = s
}

func accentCode() string {
	return config.ColorANSICode(current.Accent)
}

func titleCardColorCode() string {
	return config.ColorANSICode(current.TitleCardColor)
}

var baseAccentBright = map[string]string{
	"red":     "\033[91m",
	"green":   "\033[92m",
	"yellow":  "\033[93m",
	"blue":    "\033[94m",
	"magenta": "\033[95m",
	"cyan":    "\033[96m",
	"white":   "\033[97m",
}

func abbrevCode() string {
	if code, ok := baseAccentBright[current.Accent]; ok {
		return code
	}
	return config.ColorANSICode(current.Accent)
}

// PaintColorID renders s in the given preset color when color is enabled.
func PaintColorID(id, s string) string {
	return paint(config.ColorANSICode(config.NormalizeColorID(id)), s)
}

func logoArt() string {
	return logos.Art(current.LogoID)
}
