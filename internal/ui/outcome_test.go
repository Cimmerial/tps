package ui

import (
	"strings"
	"testing"
)

func TestOutcomeBadgeSimplifiedForcesColor(t *testing.T) {
	ApplySettings(UISettings{SimplifiedOutput: true, ColorEnabled: boolPtr(false)})
	line := OutcomeLine(true, "tests", "3 packages passed")
	if !strings.Contains(line, "\033[32m") || !strings.Contains(line, "SUCCESS") {
		t.Fatalf("expected green SUCCESS, got %q", line)
	}
	line = OutcomeLine(false, "push", "exit 1")
	if !strings.Contains(line, "\033[31m") || !strings.Contains(line, "FAILURE") {
		t.Fatalf("expected red FAILURE, got %q", line)
	}
	ApplySettings(UISettings{})
}

func boolPtr(v bool) *bool { return &v }
