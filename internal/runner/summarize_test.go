package runner

import (
	"strings"
	"testing"

	"github.com/Cimmerial/tps/internal/ui"
)

func TestSummarizeGoTestSuccess(t *testing.T) {
	out := `ok  	github.com/Cimmerial/tps/internal/cli	(cached)
?   	github.com/Cimmerial/tps/cmd/tps	[no test files]
ok  	github.com/Cimmerial/tps/internal/ui	0.026s`
	sum := summarizeGoTestSuccess(out)
	if sum.label != "tests" || !strings.Contains(sum.detail, "2 packages passed") {
		t.Fatalf("summary = %+v", sum)
	}
	line := ui.OutcomeLine(true, sum.label, sum.detail)
	if !strings.Contains(line, "SUCCESS") {
		t.Fatalf("line = %q", line)
	}
}

func TestSummarizeGoTestFailure(t *testing.T) {
	out := `FAIL	github.com/Cimmerial/tps/internal/cli	0.1s`
	sum := summarizeGoTestFailure(out)
	if !strings.Contains(sum.detail, "failed") {
		t.Fatalf("summary = %+v", sum)
	}
	line := ui.OutcomeLine(false, sum.label, sum.detail)
	if !strings.Contains(line, "FAILURE") {
		t.Fatalf("line = %q", line)
	}
}
