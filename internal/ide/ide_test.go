package ide

import (
	"strings"
	"testing"
)

func TestNormalizeID(t *testing.T) {
	if NormalizeID("") != IDAuto {
		t.Fatal("empty")
	}
	if NormalizeID("vscode") != IDVSCode {
		t.Fatal("vscode")
	}
	if NormalizeID("VS Code") != IDVSCode {
		t.Fatal("vs code phrase")
	}
	if NormalizeID("idea") != IDIntelliJ {
		t.Fatal("idea")
	}
	if NormalizeID("Android Studio") != IDAndroidStudio {
		t.Fatal("android studio")
	}
}

func TestPathsEqual(t *testing.T) {
	dir := t.TempDir()
	if !pathsEqual(dir, dir+"/") {
		t.Fatal("same dir")
	}
}

func TestFolderURIToPath(t *testing.T) {
	p, ok := folderURIToPath("file:///tmp/foo%20bar")
	if !ok || !strings.HasSuffix(p, "foo bar") {
		t.Fatalf("decode: %q ok=%v", p, ok)
	}
}
