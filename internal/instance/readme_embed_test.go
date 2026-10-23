package instance

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSyncProjectUsesEmbeddedWhenNoRepo(t *testing.T) {
	app := t.TempDir()
	wrote, err := SyncProject("", app)
	if err != nil {
		t.Fatal(err)
	}
	if !wrote {
		t.Fatal("expected embedded readme write")
	}
	got, err := os.ReadFile(filepath.Join(app, ReadmeDestName))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("empty TPS.md")
	}
	specs, err := os.ReadFile(filepath.Join(app, SpecsDestName))
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) == 0 {
		t.Fatal("empty TPS_SPECS.md")
	}
}
