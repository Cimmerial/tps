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
	ignore, err := os.ReadFile(filepath.Join(app, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	wantIgnore := "TPS.md\nTPS_SPECS.md\n"
	if string(ignore) != wantIgnore {
		t.Fatalf(".gitignore %q", ignore)
	}
	if _, err := os.Stat(filepath.Join(app, SpecsDestName)); err != nil {
		t.Fatalf("expected %s from embed: %v", SpecsDestName, err)
	}
}

func TestSyncProjectCopiesSpecsFromDisk(t *testing.T) {
	tpsRepo := t.TempDir()
	srcDir := filepath.Join(tpsRepo, "instance")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}
	specBody := "# specs from disk\n"
	if err := os.WriteFile(filepath.Join(srcDir, "SPECS.md"), []byte(specBody), 0o644); err != nil {
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
	got, err := os.ReadFile(filepath.Join(app, SpecsDestName))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != specBody {
		t.Fatalf("got %q", got)
	}
}

func TestSyncProjectGitignoreWhenAlreadyPresent(t *testing.T) {
	app := t.TempDir()
	if err := os.WriteFile(filepath.Join(app, ".gitignore"), []byte("node_modules/\nTPS.md\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := SyncProject("", app); err != nil {
		t.Fatal(err)
	}
	ignore, err := os.ReadFile(filepath.Join(app, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	if string(ignore) != "node_modules/\nTPS.md\nTPS_SPECS.md\n" {
		t.Fatalf(".gitignore %q", ignore)
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
