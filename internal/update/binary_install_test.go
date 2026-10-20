package update

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestIsFileInUseError(t *testing.T) {
	if !isFileInUseError(errors.New("The process cannot access the file because it is being used by another process.")) {
		t.Fatal("expected in-use")
	}
	if isFileInUseError(errors.New("not found")) {
		t.Fatal("unexpected in-use")
	}
}

func TestPathsEquivalent(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "tps.exe")
	b := filepath.Join(dir, "TPS.EXE")
	if !pathsEquivalent(a, b) {
		t.Fatalf("expected equal paths: %q vs %q", a, b)
	}
}

func TestIsWindowsRunningBinary(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	if !isWindowsRunningBinary(exe) {
		t.Fatalf("expected running binary %q", exe)
	}
}
