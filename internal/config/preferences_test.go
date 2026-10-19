package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPreferencesRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	on := true
	prefs := &Preferences{
		Logo:         "pixel",
		Accent:       "magenta",
		ColorEnabled: &on,
	}
	if err := SavePreferences(prefs); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "tps", "preferences.yaml")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file: %v", err)
	}
	loaded, err := LoadPreferences()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Logo != "pixel" || loaded.Accent != "magenta" {
		t.Fatalf("loaded %+v", loaded)
	}
	if loaded.ColorEnabled == nil || !*loaded.ColorEnabled {
		t.Fatal("color_enabled not preserved")
	}
}

func TestNormalizeLogoID(t *testing.T) {
	tests := map[string]string{
		"":          "default",
		"default":   "default",
		"titlecard": "default",
		"classic":   "classic",
		"tps_dots":  "tps_slash",
		"pixel":     "pixel",
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
	if p.LogoID() != "o_bird" {
		t.Fatalf("LogoID() = %q", p.LogoID())
	}
	p2 := &Preferences{Logo: "wire"}
	if p2.LogoID() != "wire" {
		t.Fatalf("LogoID() = %q", p2.LogoID())
	}
}
