package runner

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/ui"
)

// summarizeSuccessOutput returns a short line for captured stdout on success, or "" when nothing to print.
func summarizeSuccessOutput(script, combined string) string {
	if strings.Contains(script, "go test") {
		return summarizeGoTest(combined)
	}
	return ""
}

func summarizeGoTest(out string) string {
	var ok, skipped, fail int
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "ok":
			ok++
		case "?":
			skipped++
		case "FAIL":
			fail++
		}
	}
	if fail > 0 {
		return ""
	}
	total := ok + skipped
	if total == 0 {
		return ui.Green("tests ok")
	}
	if skipped > 0 && ok > 0 {
		return ui.Green(fmt.Sprintf("tests ok (%d packages, %d without tests)", ok, skipped))
	}
	if skipped > 0 {
		return ui.Green(fmt.Sprintf("tests ok (%d packages without tests)", skipped))
	}
	return ui.Green(fmt.Sprintf("tests ok (%d packages)", ok))
}
