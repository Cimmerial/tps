package cli

import (
	"strings"
	"testing"
)

func TestFormatSessionProjectLineNoTransition(t *testing.T) {
	line := formatSessionProjectLine("beta", "")
	if !strings.Contains(line, "beta") {
		t.Fatalf("line=%q want current name", line)
	}
	if strings.Contains(line, "->") {
		t.Fatalf("line=%q should not contain transition", line)
	}
}

func TestFormatSessionProjectLineWithTransition(t *testing.T) {
	line := formatSessionProjectLine("beta", "alpha")
	if !strings.Contains(line, "alpha") || !strings.Contains(line, "beta") {
		t.Fatalf("line=%q want both names", line)
	}
	if !strings.Contains(line, "->") {
		t.Fatalf("line=%q want transition arrow", line)
	}
	idxAlpha := strings.Index(line, "alpha")
	idxArrow := strings.Index(line, "->")
	idxBeta := strings.Index(line, "beta")
	if idxAlpha < 0 || idxArrow < 0 || idxBeta < 0 || !(idxAlpha < idxArrow && idxArrow < idxBeta) {
		t.Fatalf("line=%q want alpha -> beta order", line)
	}
}

func TestFormatSessionProjectLineSameName(t *testing.T) {
	line := formatSessionProjectLine("tps", "tps")
	if !strings.Contains(line, "tps") {
		t.Fatalf("line=%q", line)
	}
	if strings.Count(line, "tps") < 2 {
		t.Fatalf("line=%q want both prev and current tps", line)
	}
}
