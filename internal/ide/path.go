package ide

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func absClean(path string) (string, error) {
	if path == "" {
		return "", os.ErrInvalid
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	clean := filepath.Clean(abs)
	resolved, err := filepath.EvalSymlinks(clean)
	if err != nil {
		resolved = clean
	}
	return filepath.Clean(resolved), nil
}

func folderURIToPath(uri string) (string, bool) {
	uri = strings.TrimSpace(uri)
	if uri == "" {
		return "", false
	}
	u, err := url.Parse(uri)
	if err != nil || u.Scheme != "file" {
		return "", false
	}
	p, err := url.PathUnescape(u.Path)
	if err != nil {
		return "", false
	}
	if p == "" {
		return "", false
	}
	abs, err := absClean(p)
	if err != nil {
		return filepath.Clean(p), true
	}
	return abs, true
}

func pathsEqual(a, b string) bool {
	aa, errA := absClean(a)
	bb, errB := absClean(b)
	if errA == nil && errB == nil {
		return aa == bb
	}
	return filepath.Clean(a) == filepath.Clean(b)
}
