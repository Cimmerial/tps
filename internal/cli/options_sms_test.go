package cli

import "testing"

func TestMaskSMSTo(t *testing.T) {
	if maskSMSTo("") != "" {
		t.Fatal("empty")
	}
	if maskSMSTo("8373") != "****" {
		t.Fatalf("short: %q", maskSMSTo("8373"))
	}
	if maskSMSTo("+15558373") != "******8373" {
		t.Fatalf("long: %q", maskSMSTo("+15558373"))
	}
}
