package update

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type windowsReplace struct {
	Staging string
	Dest    string
}

var windowsReplaceQueue []windowsReplace

func pendingWindowsReplacements() []windowsReplace {
	return windowsReplaceQueue
}

func clearWindowsReplaceQueue() {
	windowsReplaceQueue = nil
}

func queueWindowsReplace(staging, dest string) {
	windowsReplaceQueue = append(windowsReplaceQueue, windowsReplace{Staging: staging, Dest: dest})
}

// installBinary copies src onto dst, staging to dst.new on Windows when dst is the running exe.
// staged is true when the live binary will be swapped in after RestartInteractive exits.
func installBinary(src, dst string) (staged bool, err error) {
	dst, err = absClean(dst)
	if err != nil {
		return false, err
	}
	if err := copyBinary(src, dst); err == nil {
		return false, nil
	}
	copyErr := err
	if runtime.GOOS != "windows" {
		return false, copyErr
	}
	if !shouldStageWindowsSelfUpdate(dst, copyErr) {
		return false, fmt.Errorf("open %s: %w", dst, copyErr)
	}
	staging := dst + ".new"
	if err := copyBinary(src, staging); err != nil {
		return false, err
	}
	queueWindowsReplace(staging, dst)
	return true, nil
}

func shouldStageWindowsSelfUpdate(dst string, copyErr error) bool {
	if isWindowsRunningBinary(dst) {
		return true
	}
	return isFileInUseError(copyErr)
}

func isWindowsRunningBinary(dst string) bool {
	exe, err := os.Executable()
	if err != nil {
		return false
	}
	return pathsEquivalent(dst, exe)
}

func pathsEquivalent(a, b string) bool {
	if strings.EqualFold(filepath.Clean(a), filepath.Clean(b)) {
		return true
	}
	aa, ea := absClean(a)
	bb, eb := absClean(b)
	if ea == nil && eb == nil && strings.EqualFold(aa, bb) {
		return true
	}
	if sameFile(a, b) {
		return true
	}
	if ea == nil && eb == nil && sameFile(aa, bb) {
		return true
	}
	return false
}

func isFileInUseError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "being used by another process") ||
		strings.Contains(msg, "the process cannot access the file")
}

func absClean(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		abs = filepath.Clean(abs)
	}
	return abs, nil
}
