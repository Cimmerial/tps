package ui

import "strings"

// Header is a top-level section title (bold, blue, ALL CAPS).
func Header(s string) string {
	return Bold(Blue(strings.ToUpper(strings.TrimSpace(s))))
}

// SubHeader is an indented nested section label (bold, blue, ALL CAPS; no literal "-").
func SubHeader(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		label = "general"
	}
	return Bold(Blue("  " + strings.ToUpper(label)))
}
