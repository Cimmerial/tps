package instance

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

func TestSyncProjectSkipsDevRepo(t *testing.T) {
	root := t.TempDir()
	dev := filepath.Join(root, "cmd", "tps")
	if err := os.MkdirAll(dev, 0o755); err != nil {
		t.Fatal(err)
	}
	wrote, err := SyncProject(root, root)
	if err != nil {
		t.Fatal(err)
	}
	if wrote {
		t.Fatal("expected no write into dev repo")
	}
}

func TestSyncProjectCopiesReadme(t *testing.T) {
	tpsRepo := t.TempDir()
	srcDir := filepath.Join(tpsRepo, "instance")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "# hello from tps\n"
	if err := os.WriteFile(filepath.Join(srcDir, "README.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	app := t.TempDir()
	wrote, err := SyncProject(tpsRepo, app)
	if err != nil {
		t.Fatal(err)
	}
	if !wrote {
		t.Fatal("expected write")
	}
	got, err := os.ReadFile(filepath.Join(app, ReadmeDestName))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != body {
		t.Fatalf("got %q", got)
	}
}

func TestSyncAllLinked(t *testing.T) {
	tpsRepo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tpsRepo, "cmd", "tps"), 0o755); err != nil {
		t.Fatal(err)
	}
	srcDir := filepath.Join(tpsRepo, "instance")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "README.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	app := t.TempDir()
	cfg := &config.Config{
		Projects: []config.Project{
			{Name: "tps", Path: tpsRepo},
			{Name: "app", Path: app},
		},
	}
	n, err := SyncAllLinked(tpsRepo, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("synced %d want 1", n)
	}
}
