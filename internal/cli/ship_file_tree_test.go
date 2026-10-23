package cli

import (
	"strings"
	"testing"
)

func TestBuildShipFileTree(t *testing.T) {
	entries := []shipFileTreeEntry{
		{Path: "internal/cli/git.go"},
		{Path: "internal/cli/git_ship_prompt.go"},
		{Path: "internal/git/git.go"},
		{Path: "internal/ui/chrome.go"},
		{Path: "internal/ui/header.go"},
		{Path: "internal/ui/errors.go"},
		{Path: "README.md"},
	}
	lines := shipFileTreeLines(buildShipFileTree(entries), func(name string) string { return name }, func(n int, name, _ string) string {
		return name
	})
	joined := strings.Join(lines, "\n")
	for _, want := range []string{"internal", "├─", "└─", "│", "git.go", "README.md"} {
		if !strings.Contains(joined, want) {
			t.Fatalf("missing %q in tree:\n%s", want, joined)
		}
	}
	if strings.Contains(joined, "internal/cli/git.go") {
		t.Fatalf("expected tree layout, not flat paths:\n%s", joined)
	}
}

func TestShipFileDisplayNumbers(t *testing.T) {
	paths := []string{
		"internal/cli/git.go",
		"internal/cli/git_ship_prompt.go",
		"internal/git/git.go",
		"internal/ui/chrome.go",
		"internal/ui/header.go",
		"internal/ui/errors.go",
	}
	nums := shipFileDisplayNumbers(paths)
	if nums["internal/cli/git.go"] != 1 {
		t.Fatalf("cli git.go should be 1, got %v", nums)
	}
	if nums["internal/git/git.go"] != 3 {
		t.Fatalf("git package git.go should be 3, got %d", nums["internal/git/git.go"])
	}
	if nums["internal/ui/errors.go"] != 5 {
		t.Fatalf("ui errors should be 5 in tree order, got %d", nums["internal/ui/errors.go"])
	}
	p, ok := shipFilePathForDisplayNumber(paths, 3)
	if !ok || p != "internal/git/git.go" {
		t.Fatalf("display 3 -> %q ok=%v", p, ok)
	}
}
