package config

import (
	"fmt"
	"runtime"
	"strings"
)

// UnmarshalYAML accepts abbreviation (preferred) and abbrev (legacy), plus shell1 / linux_shell / …
func (c *Command) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw map[string]interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	c.Alias = stringField(raw, "alias")
	c.Description = stringField(raw, "description")
	c.Abbrev = stringField(raw, "abbreviation")
	if c.Abbrev == "" {
		c.Abbrev = stringField(raw, "abbrev")
	}
	c.Detach = boolField(raw, "detach")
	c.ForceForeground = boolField(raw, "foreground")
	c.CloseOnComplete = boolField(raw, "close_on_complete")
	c.Capture = stringField(raw, "capture")
	c.CaptureActions = ParseCaptureActions(raw["capture_actions"])
	c.Tags = parseTagsField(raw["tags"])
	c.Precursor = stringField(raw, "precursor")
	c.Handler = stringField(raw, "handler")

	st := parseShellStore(raw)
	if len(st.generic) > 0 || len(st.windows) > 0 || len(st.linux) > 0 || len(st.darwin) > 0 {
		c.shells = &st
	}
	resolved := st.resolved(runtime.GOOS)
	if len(resolved) == 0 {
		c.Shell = stringField(raw, "shell")
	} else {
		c.Shell = resolved[0]
	}
	return nil
}

// MarshalYAML writes abbreviation when abbrev is set and preserves multi-shell keys.
func (c Command) MarshalYAML() (interface{}, error) {
	out := map[string]interface{}{
		"alias": c.Alias,
	}
	if c.Description != "" {
		out["description"] = c.Description
	}
	if c.Abbrev != "" {
		out["abbreviation"] = c.Abbrev
	}
	if c.Detach {
		out["detach"] = true
	}
	if c.ForceForeground {
		out["foreground"] = true
	}
	if c.CloseOnComplete {
		out["close_on_complete"] = true
	}
	if c.Capture != "" {
		out["capture"] = c.Capture
	}
	if len(c.CaptureActions) > 0 {
		out["capture_actions"] = c.CaptureActions
	}
	if len(c.Tags) > 0 {
		out["tags"] = c.Tags
	}
	if strings.TrimSpace(c.Precursor) != "" {
		out["precursor"] = strings.TrimSpace(c.Precursor)
	}
	if strings.TrimSpace(c.Handler) != "" {
		out["handler"] = strings.TrimSpace(c.Handler)
	}
	if c.shells != nil {
		for k, v := range c.shells.marshalFields() {
			out[k] = v
		}
	} else if strings.TrimSpace(c.Shell) != "" {
		out["shell"] = c.Shell
	}
	return out, nil
}

func stringField(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

func boolField(m map[string]interface{}, key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	switch x := v.(type) {
	case bool:
		return x
	default:
		s := strings.ToLower(strings.TrimSpace(fmt.Sprint(x)))
		return s == "true" || s == "yes" || s == "1"
	}
}
