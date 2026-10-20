package runner

import (
	"strings"
	"testing"
)

func TestSummarizeGoTest(t *testing.T) {
	out := `ok  	github.com/Cimmerial/tps/internal/cli	(cached)
?   	github.com/Cimmerial/tps/cmd/tps	[no test files]
ok  	github.com/Cimmerial/tps/internal/ui	0.026s`
	line := summarizeGoTest(out)
	if line == "" {
		t.Fatal("expected summary line")
	}
	if !strings.Contains(line, "tests ok") || !strings.Contains(line, "2 packages") {
		t.Fatalf("summary = %q", line)
	}
}
