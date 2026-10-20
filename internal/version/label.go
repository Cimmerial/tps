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

// Label returns the version string shown to users (release ldflags, else git describe, else "dev").
func Label() string {
	v := strings.TrimSpace(Version)
	if v != "" && v != "dev" {
		return v
	}
	if d := strings.TrimSpace(describeGit()); d != "" {
		return d
	}
	if v != "" {
		return v
	}
	return "dev"
}
