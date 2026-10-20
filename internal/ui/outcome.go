package ui

import (
	"fmt"
	"strings"
)

// OutcomeLine formats a simplified result: "label — detail  SUCCESS" (green) or FAILURE (red).
func OutcomeLine(ok bool, label, detail string) string {
	badge := outcomeBadge(ok)
	msg := strings.TrimSpace(label)
	if d := strings.TrimSpace(detail); d != "" {
		msg = msg + " — " + d
	}
	return msg + "  " + badge
}

// PrintOutcome prints OutcomeLine to stdout.
func PrintOutcome(ok bool, label, detail string) {
	fmt.Println(OutcomeLine(ok, label, detail))
}

func outcomeBadge(ok bool) string {
	text := "SUCCESS"
	code := green
	if !ok {
		text = "FAILURE"
		code = red
	}
	if ColorEnabled() || SimplifiedOutput() {
		return code + text + reset
	}
	return text
}
