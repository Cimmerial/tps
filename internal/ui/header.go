package ui

import "strings"

// Header is a top-level section title (bold, accent color, ALL CAPS).
func Header(s string) string {
	return Bold(paint(accentCode(), strings.ToUpper(strings.TrimSpace(s))))
}

// SubHeader is an indented nested section label (bold, accent color, ALL CAPS; no literal "-").
func SubHeader(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		label = "general"
	}
	return Bold(paint(accentCode(), "  "+strings.ToUpper(label)))
}
