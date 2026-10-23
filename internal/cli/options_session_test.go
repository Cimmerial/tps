package cli

import (
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

func TestOptionsSessionChanges(t *testing.T) {
	start := snapshotPreferences(&config.Preferences{})
	on := true
	after := snapshotPreferences(&config.Preferences{
		Logo:           "wire",
		TitleCardColor: "magenta",
		Accent:         "blue",
		ColorEnabled:   &on,
		SMSTo:          "+15551234567",
	})
	changes := optionsSessionChanges(start, after)
	if len(changes) != 4 {
		t.Fatalf("changes = %v", changes)
	}
	same := optionsSessionChanges(start, start)
	if len(same) != 0 {
		t.Fatalf("expected no changes, got %v", same)
	}
}
