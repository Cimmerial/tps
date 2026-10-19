package install

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestDefaultBinDir(t *testing.T) {
	dir, err := DefaultBinDir()
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		if !strings.Contains(dir, "Programs") {
			t.Fatalf("got %s", dir)
		}
	} else {
		if !strings.Contains(dir, ".local") {
			t.Fatalf("got %s", dir)
		}
	}
}

func TestDefaultBinDirUnixHome(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix only")
	}
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	got, err := DefaultBinDir()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, ".local", "bin")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
