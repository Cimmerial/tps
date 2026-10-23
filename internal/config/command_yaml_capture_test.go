package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestCommandYAMLCaptureDisabledPreservesFields(t *testing.T) {
	c := Command{
		Alias:           "tunnel",
		Capture:         "trycloudflare",
		CaptureDisabled: true,
		CaptureActions:  []string{"print", "sms"},
		Detach:          true,
	}
	raw, err := yaml.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	var back Command
	if err := yaml.Unmarshal(raw, &back); err != nil {
		t.Fatal(err)
	}
	if back.Capture != "trycloudflare" || !back.CaptureDisabled {
		t.Fatalf("capture: %q disabled=%v", back.Capture, back.CaptureDisabled)
	}
	if len(back.CaptureActions) != 2 || back.CaptureActions[1] != "sms" {
		t.Fatalf("actions: %v", back.CaptureActions)
	}
}
