package instance

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

func TestEnsureDefaultDevProjectCommands(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "cmd", "tps"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := &config.Project{Path: root, Commands: nil}
	changed, err := EnsureDefaultDevProjectCommands(p)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected changed on first ensure")
	}
	if p.FindCommand("sym") == nil {
		t.Fatal("expected sym command")
	}
	pr := p.FindCommand("push-release")
	if pr == nil || pr.Handler != "push-release" || pr.Shell != "" {
		t.Fatalf("expected push-release builtin handler, got %+v", pr)
	}
	changed, err = EnsureDefaultDevProjectCommands(p)
	if err != nil {
		t.Fatal(err)
	}
	if changed {
		t.Fatal("expected no change on second ensure")
	}
}

func TestPruneLegacyDevShellBuiltinCommands(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "cmd", "tps"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := &config.Project{
		Path: root,
		Commands: []config.Command{
			{Alias: "run", Shell: "go test ./..."},
			{Alias: "update", Shell: "go install ./cmd/tps"},
			{Alias: "close all terms", Abbrev: "vite", Shell: "osascript -e 'close'"},
		},
	}
	changed, err := EnsureDefaultDevProjectCommands(p)
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("expected prune + sym add")
	}
	if p.FindCommand("update") != nil || p.FindCommand("close all terms") != nil {
		t.Fatalf("legacy builtins should be removed: %+v", p.Commands)
	}
	if p.FindCommand("run") == nil || p.FindCommand("sym") == nil {
		t.Fatalf("expected run and sym: %+v", p.Commands)
	}
}

func TestIsNonEditableDevProjectCommand(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "cmd", "tps"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := &config.Project{Path: root, Commands: []config.Command{{Alias: "sym", Handler: "sym"}}}
	if !IsNonEditableDevProjectCommand(p, "sym") {
		t.Fatal("sym should be non-editable")
	}
	if IsNonEditableDevProjectCommand(p, "run") {
		t.Fatal("run should be editable")
	}
	other := &config.Project{Path: t.TempDir(), Commands: []config.Command{{Alias: "update", Shell: "x"}}}
	if IsNonEditableDevProjectCommand(other, "update") {
		t.Fatal("update on app repo should be editable")
	}
}

func TestEnsureDefaultDevProjectCommandsSkipsAppRepo(t *testing.T) {
	p := &config.Project{Path: t.TempDir()}
	changed, err := EnsureDefaultDevProjectCommands(p)
	if err != nil {
		t.Fatal(err)
	}
	if changed || len(p.Commands) != 0 {
		t.Fatalf("expected no commands for non-dev repo: %+v", p.Commands)
	}
}
