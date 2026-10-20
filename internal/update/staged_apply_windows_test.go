//go:build windows

package update

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSwapStagedWhileRunning(t *testing.T) {
	dir := t.TempDir()
	dest := filepath.Join(dir, "tps.exe")
	staging := dest + ".new"
	if err := os.WriteFile(dest, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(staging, []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := swapStagedWhileRunning(dest); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "new" {
		t.Fatalf("dest content: %q", got)
	}
	if _, err := os.Stat(staging); err == nil {
		t.Fatal("staging file should be gone")
	}
}
