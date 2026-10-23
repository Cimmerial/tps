package cli

import "testing"

func TestParseGitShipInput(t *testing.T) {
	tests := []struct {
		line    string
		wantMsg string
		wantOK  bool
		wantErr bool
	}{
		{"g improvements", "improvements", true, false},
		{"g ship fix bug", "fix bug", true, false},
		{"g status", "", false, false},
		{"g diff --staged", "", false, false},
		{"gs legacy msg", "legacy msg", true, false},
		{"git ship hello", "hello", true, false},
		{"g ship", "", true, false},
		{"g", "", true, false},
		{"git ship", "", true, false},
	}
	for _, tc := range tests {
		msg, ok, err := parseGitShipInput(fieldsFromLine(tc.line))
		if tc.wantErr {
			if err == nil {
				t.Errorf("%q: want error", tc.line)
			}
			continue
		}
		if err != nil {
			t.Errorf("%q: unexpected error: %v", tc.line, err)
			continue
		}
		if ok != tc.wantOK {
			t.Errorf("%q: ok=%v want %v", tc.line, ok, tc.wantOK)
		}
		if msg != tc.wantMsg {
			t.Errorf("%q: msg=%q want %q", tc.line, msg, tc.wantMsg)
		}
	}
}

func TestIsKnownGitSubcommand(t *testing.T) {
	if !isKnownGitSubcommand("status") || isKnownGitSubcommand("improvements") {
		t.Fatal("known git subcommand mismatch")
	}
}
