package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Cimmerial/tps/internal/addon"
	"github.com/Cimmerial/tps/internal/config"
	"gopkg.in/yaml.v3"
)

func TestParseConfirmYes(t *testing.T) {
	if !parseConfirmYes("y") || !parseConfirmYes("Y") || !parseConfirmYes(" yes ") {
		t.Fatal("expected yes")
	}
	if parseConfirmYes("") || parseConfirmYes("n") || parseConfirmYes("no") {
		t.Fatal("expected no")
	}
}

func TestLinkedUtilitiesProjects(t *testing.T) {
	dir := t.TempDir()
	utilRoot := filepath.Join(dir, "tps-utilities")
	if err := os.MkdirAll(utilRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	m := addon.DefaultUtilitiesManifest()
	data, err := yaml.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(utilRoot, addon.ManifestName), data, 0o644); err != nil {
		t.Fatal(err)
	}
	other := filepath.Join(dir, "app")
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := &config.Config{
		Projects: []config.Project{
			{ID: "u", Name: "tps-utilities", Path: utilRoot},
			{ID: "a", Name: "app", Path: other},
		},
	}
	linked := linkedUtilitiesProjects(cfg)
	if len(linked) != 1 || linked[0].ID != "u" {
		t.Fatalf("linked=%v", linked)
	}
}
