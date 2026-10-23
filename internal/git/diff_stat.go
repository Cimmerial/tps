package git

import (
	"os/exec"
	"strconv"
	"strings"
)

// PathLineDiff is insert/delete line counts vs HEAD (working tree + index).
type PathLineDiff struct {
	Added   int
	Removed int
	Binary  bool
}

// DiffNumstatAgainstHEAD returns per-path line changes vs HEAD.
func DiffNumstatAgainstHEAD(dir string) (map[string]PathLineDiff, error) {
	cmd := exec.Command("git", "diff", "HEAD", "--numstat")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseDiffNumstat(out), nil
}

func parseDiffNumstat(out []byte) map[string]PathLineDiff {
	m := make(map[string]PathLineDiff)
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}
		path := strings.TrimSpace(parts[2])
		if path == "" {
			continue
		}
		if parts[0] == "-" || parts[1] == "-" {
			m[path] = PathLineDiff{Binary: true}
			continue
		}
		added, errA := strconv.Atoi(parts[0])
		removed, errR := strconv.Atoi(parts[1])
		if errA != nil || errR != nil {
			continue
		}
		m[path] = PathLineDiff{Added: added, Removed: removed}
	}
	return m
}

// IsUntrackedPath reports whether path is an untracked file under dir.
func IsUntrackedPath(dir, path string) bool {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard", "--", path)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) != ""
}
