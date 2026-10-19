package config

import "strings"

// ResolveCapturePattern expands friendly names or returns a regex pattern as-is.
func ResolveCapturePattern(s string) string {
	s = strings.TrimSpace(s)
	switch strings.ToLower(s) {
	case "trycloudflare", "cloudflare", "cf":
		return `https://[a-z0-9-]+\.trycloudflare\.com`
	default:
		return s
	}
}
