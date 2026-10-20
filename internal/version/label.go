package version

import (
	"os/exec"
	"strings"
)

// describeGit is overridden in tests.
var describeGit = defaultDescribeGit

func defaultDescribeGit() string {
	out, err := exec.Command("git", "describe", "--tags", "--always", "--dirty").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// Label returns the version string shown to users (release ldflags, else module vX.Y.Z from the binary).
func Label() string {
	v := strings.TrimSpace(Version)
	if v != "" && v != "dev" {
		return v
	}
	if d := strings.TrimSpace(moduleVersionLabel()); d != "" {
		return d
	}
	if d := strings.TrimSpace(describeGit()); d != "" && !isBareCommitHash(d) {
		return ensureVPrefix(d)
	}
	if DevVersion != "" {
		return "v" + strings.TrimPrefix(DevVersion, "v")
	}
	if v != "" {
		return v
	}
	return "dev"
}

func ensureVPrefix(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || strings.HasPrefix(s, "v") {
		return s
	}
	if strings.Contains(s, ".") {
		return "v" + s
	}
	return s
}

func isBareCommitHash(s string) bool {
	if len(s) < 7 || len(s) > 40 {
		return false
	}
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
		case c >= 'a' && c <= 'f':
		case c >= 'A' && c <= 'F':
		default:
			return false
		}
	}
	return true
}
