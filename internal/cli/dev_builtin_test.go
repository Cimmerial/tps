package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/instance"
)

func TestRunDevBuiltinSym(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "cmd", "tps"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := &config.Project{Path: root}
	if _, err := instance.EnsureDefaultDevProjectCommands(p); err != nil {
		t.Fatal(err)
	}
	_, _, _ = runDevBuiltin(p, "sym", nil) // prints catalog
	handled, code, err := runDevBuiltin(p, "sym", []string{"block_dark"})
	if err != nil || code != 0 || !handled {
		t.Fatalf("sym: handled=%v code=%d err=%v", handled, code, err)
	}
}

func TestRunDevBuiltinSymFilter(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "cmd", "tps"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := &config.Project{Path: root}
	if _, err := instance.EnsureDefaultDevProjectCommands(p); err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	handled, code, err := runDevBuiltin(p, "sym", []string{"check"})
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	if err != nil || code != 0 || !handled {
		t.Fatalf("sym check: handled=%v code=%d err=%v", handled, code, err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("✓")) {
		t.Fatalf("expected check glyph in output: %s", buf.String())
	}
}

func TestRunDevBuiltinNotDevRepo(t *testing.T) {
	p := &config.Project{Path: t.TempDir()}
	handled, _, err := runDevBuiltin(p, "sym", nil)
	if err != nil || handled {
		t.Fatalf("expected not handled: handled=%v err=%v", handled, err)
	}
}
