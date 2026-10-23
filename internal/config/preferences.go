package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Preferences are stored separately from project config (same tps dir as config.yaml; see PreferencesPath).
type Preferences struct {
	Logo           string `yaml:"logo,omitempty"`
	TitleCard      string `yaml:"titlecard,omitempty"`
	TitleCardColor string `yaml:"title_card_color,omitempty"` // see config.ColorPresets
	Accent         string `yaml:"accent,omitempty"`             // see config.ColorPresets
	ColorEnabled   *bool  `yaml:"color_enabled,omitempty"`
	SMSTo              string `yaml:"sms_to,omitempty"`                // E.164 for capture SMS
	SMSShortcut        string `yaml:"sms_shortcut,omitempty"`          // optional Shortcuts name (preferred over AppleScript)
	SessionProjectID   string `yaml:"session_project_id,omitempty"`    // interactive session project (survives switch / restart)
	PreferredIDE       string `yaml:"preferred_ide,omitempty"`         // cursor, vscode, auto, …
	PreferredBrowser   string `yaml:"preferred_browser,omitempty"`     // default, chrome, firefox, …
	SimplifiedOutput   bool   `yaml:"simplified_output,omitempty"`     // summarize git/tests/scripts (default verbose)
	Username           string `yaml:"username,omitempty"`              // used in export filenames and bundles
}

func PreferencesPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tps", "preferences.yaml"), nil
}

func LoadPreferences() (*Preferences, error) {
	p, err := PreferencesPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Preferences{}, nil
		}
		return nil, err
	}
	var prefs Preferences
	if err := yaml.Unmarshal(data, &prefs); err != nil {
		return nil, err
	}
	return &prefs, nil
}

// SaveSessionProjectID records the active interactive project (merges into existing preferences).
func SaveSessionProjectID(projectID string) error {
	projectID = strings.TrimSpace(projectID)
	prefs, err := LoadPreferences()
	if err != nil {
		return err
	}
	if prefs.SessionProjectID == projectID {
		return nil
	}
	prefs.SessionProjectID = projectID
	return SavePreferences(prefs)
}

func SavePreferences(prefs *Preferences) error {
	p, err := PreferencesPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(prefs)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}

// LogoID returns the normalized title-card preset id (default when unset).
func (p *Preferences) LogoID() string {
	raw := strings.TrimSpace(p.Logo)
	if raw == "" {
		raw = strings.TrimSpace(p.TitleCard)
	}
	return NormalizeLogoID(raw)
}

func NormalizeLogoID(id string) string {
	id = strings.ToLower(strings.TrimSpace(id))
	switch id {
	case "", "default":
		return "default"
	case "logo", "titlecard", "title-card":
		return "default"
	case "random":
		return "random"
	// legacy preset ids → README names
	case "classic":
		return "baby"
	case "wire":
		return "goopy"
	case "o_bird":
		return "people"
	case "dollar":
		return "dolla_dolla_bill_yall"
	case "tps_blocks":
		return "boring"
	case "tps_gradient":
		return "blurry"
	case "tps_slash", "tps_dots":
		return "poker_face"
	case "pixel":
		return "block"
	case "slashes":
		return "topo"
	case "oo":
		return "lowkey_cool"
	case "sss":
		return "sharp"
	default:
		if id == "" {
			return "default"
		}
		return id
	}
}

// AccentMode returns a canonical color preset id (default cyan).
func (p *Preferences) AccentMode() string {
	return NormalizeColorID(p.Accent)
}

// TitleCardColorMode returns a canonical color preset id for the title card logo (default cyan).
func (p *Preferences) TitleCardColorMode() string {
	return NormalizeColorID(p.TitleCardColor)
}
