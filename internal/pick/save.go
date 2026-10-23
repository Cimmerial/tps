package pick

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolveSavePath appends defaultName when path is an existing directory.
func ResolveSavePath(path, defaultName string) string {
	path = strings.TrimSpace(path)
	defaultName = strings.TrimSpace(defaultName)
	if path == "" {
		return path
	}
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return filepath.Join(path, defaultName)
	}
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, string(os.PathSeparator)) {
		return filepath.Join(strings.TrimRight(path, "/\\"), defaultName)
	}
	return path
}
