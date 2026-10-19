package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestPreferencesRoundTrip(t *testing.T) {
	dir := t.TempDir()
	var tpsDir string
	switch runtime.GOOS {
	case "darwin":
		t.Setenv("HOME", dir)
		tpsDir = filepath.Join(dir, "Library", "Application Support", "tps")
	case "windows":
		appData := filepath.Join(dir, "AppData", "Roaming")
		t.Setenv("APPDATA", appData)
		tpsDir = filepath.Join(appData, "tps")
	default:
		cfg := filepath.Join(dir, ".config")
		t.Setenv("XDG_CONFIG_HOME", cfg)
		tpsDir = filepath.Join(cfg, "tps")
	}
	if err := os.MkdirAll(tpsDir, 0o700); err != nil {
		t.Fatal(err)
	}

	on := true
	prefs := &Preferences{
		Logo:         "block",
		Accent:       "magenta",
		ColorEnabled: &on,
		SMSTo:        "+15551234567",
	}
	if err := SavePreferences(prefs); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(tpsDir, "preferences.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file: %v", err)
	}
	loaded, err := LoadPreferences()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Logo != "block" || loaded.Accent != "magenta" {
		t.Fatalf("loaded %+v", loaded)
	}
	if loaded.ColorEnabled == nil || !*loaded.ColorEnabled {
		t.Fatal("color_enabled not preserved")
	}
	if loaded.SMSTo != "+15551234567" {
		t.Fatalf("sms_to = %q", loaded.SMSTo)
	}
}

func TestNormalizeLogoID(t *testing.T) {
	tests := map[string]string{
		"":          "default",
		"default":   "default",
		"titlecard": "default",
		"random":    "random",
		"classic":   "baby",
		"wire":      "goopy",
		"tps_dots":  "poker_face",
		"pixel":     "block",
	}
	for in, want := range tests {
		if got := NormalizeLogoID(in); got != want {
			t.Errorf("NormalizeLogoID(%q) = %q want %q", in, got, want)
		}
	}
}

func TestTitleCardColorMode(t *testing.T) {
	p := &Preferences{TitleCardColor: "blue"}
	if p.TitleCardColorMode() != "blue" {
		t.Fatalf("TitleCardColorMode() = %q", p.TitleCardColorMode())
	}
	if (&Preferences{}).TitleCardColorMode() != "cyan" {
		t.Fatal("empty title card color should default to cyan")
	}
}

func TestPreferencesLogoID(t *testing.T) {
	empty := &Preferences{}
	if empty.LogoID() != "default" {
		t.Fatalf("empty LogoID() = %q want default", empty.LogoID())
	}
	p := &Preferences{TitleCard: "o_bird"}
	if p.LogoID() != "people" {
		t.Fatalf("LogoID() = %q", p.LogoID())
	}
	p2 := &Preferences{Logo: "wire"}
	if p2.LogoID() != "goopy" {
		t.Fatalf("LogoID() = %q", p2.LogoID())
	}
	p3 := &Preferences{Logo: "random"}
	if p3.LogoID() != "random" {
		t.Fatalf("LogoID() = %q", p3.LogoID())
	}
}
