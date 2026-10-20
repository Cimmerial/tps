package runner

import (
	"fmt"
	"strings"
)

type outcomeSummary struct {
	label  string
	detail string
}

// summarizeSuccess returns a short description for a successful captured run.
func summarizeSuccess(script, combined string) outcomeSummary {
	if strings.Contains(script, "go test") {
		return summarizeGoTestSuccess(combined)
	}
	script = strings.TrimSpace(script)
	if script == "" {
		return outcomeSummary{label: "run", detail: "completed"}
	}
	if len(script) > 72 {
		script = script[:69] + "..."
	}
	return outcomeSummary{label: "run", detail: script}
}

func summarizeFailure(script, combined string, exitCode int) outcomeSummary {
	if strings.Contains(script, "go test") {
		if s := summarizeGoTestFailure(combined); s.detail != "" {
			return s
		}
	}
	if exitCode != 0 {
		return outcomeSummary{label: stepLabel(script), detail: fmt.Sprintf("exit %d", exitCode)}
	}
	return outcomeSummary{label: stepLabel(script), detail: "failed"}
}

func stepLabel(script string) string {
	if strings.Contains(script, "go test") {
		return "tests"
	}
	return "run"
}

func summarizeGoTestSuccess(out string) outcomeSummary {
	ok, skipped, fail := countGoTestLines(out)
	if fail > 0 {
		return summarizeGoTestFailure(out)
	}
	detail := "all packages passed"
	if ok > 0 && skipped > 0 {
		detail = fmt.Sprintf("%d packages passed, %d without tests", ok, skipped)
	} else if skipped > 0 {
		detail = fmt.Sprintf("%d packages without tests", skipped)
	} else if ok > 0 {
		detail = fmt.Sprintf("%d packages passed", ok)
	}
	return outcomeSummary{label: "tests", detail: detail}
}

func summarizeGoTestFailure(out string) outcomeSummary {
	_, _, fail := countGoTestLines(out)
	if fail > 0 {
		return outcomeSummary{label: "tests", detail: fmt.Sprintf("%d package(s) failed", fail)}
	}
	return outcomeSummary{label: "tests", detail: "failed"}
}

func countGoTestLines(out string) (ok, skipped, fail int) {
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
	return ok, skipped, fail
}
