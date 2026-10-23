package ui

import "strings"

// Header is a top-level info section title (bold magenta, ALL CAPS).
func Header(s string) string {
	return HeaderInfo(s)
}

// HeaderInfo is a top-level info section title (bold magenta, ALL CAPS).
func HeaderInfo(s string) string {
	return Bold(paint(magenta, strings.ToUpper(strings.TrimSpace(s))))
}

// HeaderChoice is a top-level choice menu title (bold red, ALL CAPS).
func HeaderChoice(s string) string {
	return Bold(paint(red, strings.ToUpper(strings.TrimSpace(s))))
}

// SubHeader is an indented info section label (bold magenta, ALL CAPS).
func SubHeader(label string) string {
	return SubHeaderInfo(label)
}

// SubHeaderInfo is an indented info section label (bold magenta, ALL CAPS).
func SubHeaderInfo(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		label = "general"
	}
	return Bold(paint(magenta, strings.ToUpper(label)))
}

// SubHeaderChoice is an indented choice section label (bold red, ALL CAPS).
func SubHeaderChoice(label string) string {
	label = strings.TrimSpace(label)
	if label == "" {
		label = "general"
	}
	return Bold(paint(red, strings.ToUpper(label)))
}

// HeaderCommit is a top-level ship-included title (bold green, ALL CAPS).
func HeaderCommit(s string) string {
	return Bold(paint(green, strings.ToUpper(strings.TrimSpace(s))))
}
