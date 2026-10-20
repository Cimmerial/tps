package addonmenu

import (
	"regexp"
	"strings"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func plainLine(s string) string {
	return strings.TrimSpace(ansiEscape.ReplaceAllString(s, ""))
}

// actionSummaryLines picks user-visible outcome lines from a command's stderr.
func actionSummaryLines(stderr string) []string {
	var out []string
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		p := plainLine(line)
		if strings.HasPrefix(p, "wrote ") ||
			strings.HasPrefix(p, "unchanged:") ||
			strings.HasPrefix(p, "skipped ") ||
			strings.HasPrefix(p, "mass convert:") ||
			strings.Contains(p, "SUCCESS:") ||
			strings.Contains(p, "FAILED:") {
			out = append(out, line)
		}
	}
	return out
}
