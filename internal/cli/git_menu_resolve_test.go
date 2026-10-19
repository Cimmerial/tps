package cli

import (
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

func TestResolveMenuInputGitShortcuts(t *testing.T) {
	p := &config.Project{Commands: []config.Command{}}
	cases := []struct {
		line       string
		wantAction InputAction
		wantText   string
		wantErr    bool
	}{
		{"g improvements", ActionGitShip, "improvements", false},
		{"g ship fix bug", ActionGitShip, "fix bug", false},
		{"g status", ActionGit, "status", false},
		{"gs legacy msg", ActionGitShip, "legacy msg", false},
	}
	for _, tc := range cases {
		in, err := resolveMenuInput(nil, p, tc.line)
		if tc.wantErr {
			if err == nil {
				t.Errorf("%q: want error", tc.line)
			}
			continue
		}
		if err != nil {
			t.Errorf("%q: %v", tc.line, err)
			continue
		}
		if in.Action != tc.wantAction {
			t.Errorf("%q: action=%v want %v text=%q", tc.line, in.Action, tc.wantAction, in.Text)
		}
		if in.Text != tc.wantText {
			t.Errorf("%q: text=%q want %q", tc.line, in.Text, tc.wantText)
		}
	}
}
