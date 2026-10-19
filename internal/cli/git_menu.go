package cli

import (
	"fmt"
	"strings"
)

func fieldsFromLine(line string) []string {
	return strings.Fields(strings.TrimSpace(line))
}

// isKnownGitSubcommand reports whether word is a real git first argument (passthrough via g/cmd).
func isKnownGitSubcommand(word string) bool {
	switch strings.ToLower(strings.TrimSpace(word)) {
	case "status", "diff", "log", "add", "commit", "push", "pull", "checkout", "branch",
		"fetch", "merge", "rebase", "stash", "show", "tag", "remote", "clone", "init",
		"reset", "restore", "rev-parse", "cherry-pick", "clean", "blame", "bisect",
		"grep", "mv", "rm", "switch", "worktree", "config", "describe", "shortlog",
		"reflog", "archive", "bundle", "fsck", "gc", "maintenance", "pull-request",
		"submodule", "subtree", "whatchanged":
		return true
	default:
		return false
	}
}

// parseGitShipInput handles g ship, git ship, and legacy gs.
func parseGitShipInput(fields []string) (msg string, ok bool, err error) {
	if len(fields) == 0 {
		return "", false, nil
	}
	switch strings.ToLower(fields[0]) {
	case "gs":
		if len(fields) < 2 {
			return "", true, fmt.Errorf("g ship needs a message (example: g ship fix the menu)")
		}
		return strings.TrimSpace(strings.Join(fields[1:], " ")), true, nil
	case "g":
		if len(fields) < 2 {
			return "", true, fmt.Errorf("g needs a git command (example: g status) or ship message (example: g fix the menu)")
		}
		if strings.EqualFold(fields[1], "ship") {
			msg = strings.TrimSpace(strings.Join(fields[2:], " "))
			if msg == "" {
				return "", true, fmt.Errorf("g ship needs a message (example: g ship fix the menu)")
			}
			return msg, true, nil
		}
		if isKnownGitSubcommand(fields[1]) {
			return "", false, nil
		}
		msg = strings.TrimSpace(strings.Join(fields[1:], " "))
		if msg == "" {
			return "", true, fmt.Errorf("g ship needs a message (example: g fix the menu)")
		}
		return msg, true, nil
	case "git":
		if len(fields) < 3 || !strings.EqualFold(fields[1], "ship") {
			return "", false, nil
		}
		msg = strings.TrimSpace(strings.Join(fields[2:], " "))
		if msg == "" {
			return "", true, fmt.Errorf("git ship needs a message (example: git ship fix the menu)")
		}
		return msg, true, nil
	default:
		return "", false, nil
	}
}
