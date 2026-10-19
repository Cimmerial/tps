package notify

import (
	"os"
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

// Run with: TPS_LIVE_SMS=1 go test ./internal/notify -run TestLiveSMS -v -count=1
// Requires sms_to in preferences.yaml (same dir as config.yaml) and Messages/Shortcuts on macOS.
// PASS means Messages accepted the send command; delivery can still fail in the UI (Text Message
// Forwarding, existing thread, E.164 matching your contact). Use -count=1 so Go does not cache.
func TestLiveSMS(t *testing.T) {
	if os.Getenv("TPS_LIVE_SMS") != "1" {
		t.Skip("set TPS_LIVE_SMS=1 to send a test message")
	}
	prefs, err := config.LoadPreferences()
	if err != nil {
		t.Fatal(err)
	}
	if prefs.SMSTo == "" {
		p, err := config.PreferencesPath()
		if err != nil {
			t.Skip("set sms_to in preferences.yaml to run live SMS test")
		}
		t.Skipf("set sms_to in %s (or options → capture sms) to run live SMS test", p)
	}
	msg := "tps capture sms test"
	if err := SendSMS(prefs.SMSTo, msg, prefs.SMSShortcut); err != nil {
		t.Fatal(err)
	}
	t.Log("queued message via Messages; if the app shows Message Failed, check SMS relay and send once manually to that number in Messages")
}
