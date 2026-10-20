package update

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	exe, err := os.Executable()
	if err != nil {
		return false, err
	}
	running, err := absClean(exe)
	if err != nil {
		return false, err
	}
	if !sameFile(dst, running) {
		return false, fmt.Errorf("open %s: %w", dst, copyErr)
	}
	staging := dst + ".new"
	if err := copyBinary(src, staging); err != nil {
		return false, err
	}
	queueWindowsReplace(staging, dst)
	return true, nil
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
