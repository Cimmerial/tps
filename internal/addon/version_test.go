package addon

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectVersionFromManifest(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ManifestName)
	if err := os.WriteFile(path, []byte("id: test\nversion: 2.3.4\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if v := ProjectVersion(dir); v != "2.3.4" {
		t.Fatalf("ProjectVersion = %q", v)
	}
}

func TestIsAddonProject(t *testing.T) {
	dir := t.TempDir()
	if IsAddonProject(dir) {
		t.Fatal("expected false without manifest")
	}
	path := filepath.Join(dir, ManifestName)
	if err := os.WriteFile(path, []byte("id: test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !IsAddonProject(dir) {
		t.Fatal("expected true with manifest")
	}
}
