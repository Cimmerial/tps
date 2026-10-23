package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// ErrShipSelectCancel is returned when the user cancels the ship file picker.
var ErrShipSelectCancel = errors.New("ship file selection cancelled")

// ShipFileSelection tracks which changed paths are included in the commit.
type ShipFileSelection struct {
	Paths    []string
	Included []bool
}

// NewShipFileSelection starts with every path included.
func NewShipFileSelection(paths []string) *ShipFileSelection {
	inc := make([]bool, len(paths))
	for i := range inc {
		inc[i] = true
	}
	return &ShipFileSelection{Paths: paths, Included: inc}
}

// Toggle flips inclusion for the 1-based index. Returns false when out of range.
func (s *ShipFileSelection) Toggle(oneBased int) bool {
	i := oneBased - 1
	if i < 0 || i >= len(s.Paths) {
		return false
	}
	s.Included[i] = !s.Included[i]
	return true
}

// TogglePath flips inclusion for a repo-relative path.
func (s *ShipFileSelection) TogglePath(path string) bool {
	for i, p := range s.Paths {
		if p == path {
			s.Included[i] = !s.Included[i]
			return true
		}
	}
	return false
}

// IncludedPaths returns paths marked for commit.
func (s *ShipFileSelection) IncludedPaths() []string {
	var out []string
	for i, p := range s.Paths {
		if s.Included[i] {
			out = append(out, p)
		}
	}
	return out
}

// ExcludedPaths returns paths not marked for commit.
func (s *ShipFileSelection) ExcludedPaths() []string {
	var out []string
	for i, p := range s.Paths {
		if !s.Included[i] {
			out = append(out, p)
		}
	}
	return out
}

// ParseShipFileSelectInput maps interactive input: empty or 0 = quit, c = commit, number = toggle.
func ParseShipFileSelectInput(line string) (commit, quit bool, toggleOneBased int, ok bool) {
	line = strings.TrimSpace(strings.ToLower(line))
	if line == "" || line == "0" {
		return false, true, 0, true
	}
	if line == "c" {
		return true, false, 0, true
	}
	n, err := strconv.Atoi(line)
	if err != nil || n < 1 {
		return false, false, 0, false
	}
	return false, false, n, true
}

// ListChangedPaths returns paths with local changes (porcelain v1, same order as git status).
func ListChangedPaths(dir string) ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain=v1", "-z", "-uall")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parsePorcelainV1Null(out), nil
}

func parsePorcelainV1Null(data []byte) []string {
	var paths []string
	seen := make(map[string]bool)
	i := 0
	for i < len(data) {
		if data[i] == 0 {
			i++
			continue
		}
		j := i
		for j < len(data) && data[j] != 0 {
			j++
		}
		line := string(data[i:j])
		i = j + 1
		if len(line) < 3 {
			continue
		}
		xy := line[0:2]
		rest := strings.TrimSpace(line[2:])
		if xy[0] == 'R' || xy[0] == 'C' {
			if i >= len(data) {
				continue
			}
			k := i
			for k < len(data) && data[k] != 0 {
				k++
			}
			newPath := unquotePorcelainPath(string(data[i:k]))
			i = k + 1
			if newPath != "" && !seen[newPath] {
				seen[newPath] = true
				paths = append(paths, newPath)
			}
			continue
		}
		p := unquotePorcelainPath(rest)
		if p != "" && !seen[p] {
			seen[p] = true
			paths = append(paths, p)
		}
	}
	return paths
}

func unquotePorcelainPath(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
		s = strings.ReplaceAll(s, "\\\"", "\"")
		s = strings.ReplaceAll(s, "\\\\", "\\")
	}
	return s
}

// StageShipFiles stages included paths and unstages excluded paths from the index.
func StageShipFiles(dir string, allPaths []string, include []string) error {
	includeSet := make(map[string]bool, len(include))
	for _, p := range include {
		includeSet[p] = true
	}
	var excluded []string
	for _, p := range allPaths {
		if !includeSet[p] {
			excluded = append(excluded, p)
		}
	}
	if len(excluded) > 0 {
		if err := unstagePaths(dir, excluded); err != nil {
			return err
		}
	}
	if len(include) > 0 {
		args := []string{"add", "--"}
		args = append(args, include...)
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git add: %w", err)
		}
	}
	return nil
}
