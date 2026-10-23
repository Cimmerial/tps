package config

import (
	"fmt"
	"regexp"
	"strings"
)

// Generic HTTP(S) URL token (trim trailing punctuation after match).
const captureURLPattern = `https?://[^\s<>"']+`

// ResolveCapturePattern expands friendly names or returns a regex pattern as-is.
func ResolveCapturePattern(s string) string {
	s = strings.TrimSpace(s)
	switch strings.ToLower(s) {
	case "trycloudflare", "cloudflare", "cf":
		return `https://[a-z0-9-]+\.trycloudflare\.com`
	case "url", "link":
		return captureURLPattern
	default:
		return s
	}
}

// TrimCapturedMatch cleans a regex match (trailing punctuation from log lines).
func TrimCapturedMatch(s string) string {
	s = strings.TrimSpace(s)
	for len(s) > 0 {
		last := s[len(s)-1]
		if last == ')' || last == ']' || last == ',' || last == '.' || last == ';' || last == '"' || last == '\'' {
			s = s[:len(s)-1]
			continue
		}
		break
	}
	return s
}

// ExtractCapture finds the first match for pattern in log text.
func ExtractCapture(logText, pattern string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	m := re.FindString(logText)
	if m == "" {
		return "", nil
	}
	return TrimCapturedMatch(m), nil
}

// Capture action names.
const (
	CapturePrint     = "print"
	CaptureClipboard = "clipboard"
	CaptureSMS       = "sms"
)

// DefaultCaptureActions is used when capture is set but capture_actions is omitted.
func DefaultCaptureActions() []string {
	return []string{CapturePrint, CaptureClipboard}
}

// ParseCaptureActions normalizes a YAML list or comma-separated string.
func ParseCaptureActions(raw interface{}) []string {
	if raw == nil {
		return nil
	}
	var parts []string
	switch v := raw.(type) {
	case string:
		for _, p := range strings.Split(v, ",") {
			parts = append(parts, strings.TrimSpace(p))
		}
	case []interface{}:
		for _, item := range v {
			parts = append(parts, strings.TrimSpace(fmt.Sprint(item)))
		}
	case []string:
		parts = append(parts, v...)
	default:
		parts = append(parts, strings.TrimSpace(fmt.Sprint(v)))
	}
	return normalizeCaptureActions(parts)
}

func normalizeCaptureActions(parts []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, p := range parts {
		p = strings.ToLower(strings.TrimSpace(p))
		switch p {
		case CapturePrint, CaptureClipboard, CaptureSMS:
			if !seen[p] {
				seen[p] = true
				out = append(out, p)
			}
		}
	}
	return out
}

// CaptureActive is true when capture should run (pattern set and not disabled).
func (c *Command) CaptureActive() bool {
	if c == nil {
		return false
	}
	return strings.TrimSpace(c.Capture) != "" && !c.CaptureDisabled
}

// EffectiveCaptureActions returns actions for a command with capture active.
func (c *Command) EffectiveCaptureActions() []string {
	if !c.CaptureActive() {
		return nil
	}
	if len(c.CaptureActions) > 0 {
		return normalizeCaptureActions(c.CaptureActions)
	}
	return DefaultCaptureActions()
}

func (c *Command) HasCaptureAction(name string) bool {
	for _, a := range c.EffectiveCaptureActions() {
		if a == name {
			return true
		}
	}
	return false
}

// HasStoredCaptureAction reports a saved action (for edit UI while capture is off).
func (c *Command) HasStoredCaptureAction(name string) bool {
	for _, a := range c.storedCaptureActions() {
		if a == name {
			return true
		}
	}
	return false
}

func (c *Command) storedCaptureActions() []string {
	if c == nil {
		return nil
	}
	if len(c.CaptureActions) > 0 {
		return normalizeCaptureActions(c.CaptureActions)
	}
	if strings.TrimSpace(c.Capture) != "" {
		return DefaultCaptureActions()
	}
	return nil
}
