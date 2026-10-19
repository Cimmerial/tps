package git

import "testing"

func TestParseShipSuspiciousChoice(t *testing.T) {
	cases := []struct {
		in      string
		want    SuspiciousAction
		wantOK  bool
	}{
		{"1", SuspiciousProceed, true},
		{"y", SuspiciousProceed, true},
		{"Y", SuspiciousProceed, true},
		{"yes", SuspiciousProceed, true},
		{"2", SuspiciousCancel, true},
		{"n", SuspiciousCancel, true},
		{"3", SuspiciousExcludeGitignore, true},
		{"", 0, false},
		{"maybe", 0, false},
	}
	for _, tc := range cases {
		got, ok := ParseShipSuspiciousChoice(tc.in)
		if ok != tc.wantOK || got != tc.want {
			t.Errorf("ParseShipSuspiciousChoice(%q) = (%v, %v), want (%v, %v)", tc.in, got, ok, tc.want, tc.wantOK)
		}
	}
}
