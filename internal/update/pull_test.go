package update

import (
	"path/filepath"
	"strings"
	"testing"
)

func gitBaseName(path string) string {
	base := strings.ToLower(filepath.Base(path))
	if base == "git.exe" {
		return "git"
	}
	return base
}

func TestGitPullCommand(t *testing.T) {
	root := filepath.Join("/tmp", "tps")
	cmd := gitPullCommand(root)
	if gitBaseName(cmd.Path) != "git" {
		t.Fatalf("path: got %q want git binary", cmd.Path)
	}
	want := []string{"-C", root, "pull", "--ff-only"}
	if len(cmd.Args) != len(want)+1 {
		t.Fatalf("args: got %v", cmd.Args)
	}
	if gitBaseName(cmd.Args[0]) != "git" {
		t.Fatalf("args[0]: got %q want git", cmd.Args[0])
	}
	for i, w := range want {
		if cmd.Args[i+1] != w {
			t.Fatalf("args[%d]: got %q want %q", i+1, cmd.Args[i+1], w)
		}
	}
}
