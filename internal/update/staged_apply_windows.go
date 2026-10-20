//go:build windows

package update

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	envStagedFrom   = "TPS_APPLY_STAGED_FROM"
	envStagedTo     = "TPS_APPLY_STAGED_TO"
	envStagedDest   = "TPS_STAGED_DEST"
	envStagedParent = "TPS_STAGED_PARENT_PID"
	FlagApplyStaged = "--tps-apply-staged"
)

// ApplyStagedWindowsBinaryIfNeeded swaps tps.exe.new into place (after u or install).
// When true, the on-disk binary was updated but this process may still be the previous build — run tps again.
func ApplyStagedWindowsBinaryIfNeeded() (needReload bool, err error) {
	dest, err := windowsExecutablePath()
	if err != nil {
		return false, err
	}
	cleanupLeftoverWindowsBackup(dest)

	to, from := stagedWindowsPaths(dest)
	if from == "" || to == "" {
		return false, nil
	}
	handoff := os.Getenv(envStagedFrom) != ""
	if handoff {
		time.Sleep(350 * time.Millisecond)
	}
	tries := 8
	pause := 80 * time.Millisecond
	if handoff {
		tries = 25
		pause = 150 * time.Millisecond
	}
	var lastErr error
	for i := 0; i < tries; i++ {
		err := swapStagedWhileRunning(to)
		if err == nil {
			_ = os.Unsetenv(envStagedFrom)
			_ = os.Unsetenv(envStagedTo)
			return true, nil
		}
		lastErr = err
		time.Sleep(pause)
	}
	return false, fmt.Errorf("apply staged binary %s -> %s: %w", from, to, lastErr)
}

// ScheduleWindowsStagedApplyAfterExit waits for this process to exit, then replaces tps.exe (no file lock).
func ScheduleWindowsStagedApplyAfterExit() {
	dest, err := windowsExecutablePath()
	if err != nil {
		return
	}
	if _, err := os.Stat(dest + ".new"); err != nil {
		return
	}
	cmd := exec.Command("powershell.exe", "-NoProfile", "-WindowStyle", "Hidden", "-Command", windowsStagedApplyAfterExitPS)
	cmd.Env = append(os.Environ(),
		envStagedDest+"="+dest,
		fmt.Sprintf("%s=%d", envStagedParent, os.Getpid()),
	)
	_ = cmd.Start()
}

const windowsStagedApplyAfterExitPS = `
$dest = $env:TPS_STAGED_DEST
$ppid = [int]$env:TPS_STAGED_PARENT_PID
if (-not $dest) { exit 0 }
Wait-Process -Id $ppid -ErrorAction SilentlyContinue
Start-Sleep -Milliseconds 450
$new = $dest + '.new'
if (-not (Test-Path -LiteralPath $new)) { exit 0 }
$old = $dest + '.old'
Remove-Item -LiteralPath $old -Force -ErrorAction SilentlyContinue
try {
  Move-Item -LiteralPath $new -Destination $dest -Force
} catch {
  try {
    if (Test-Path -LiteralPath $dest) {
      Move-Item -LiteralPath $dest -Destination $old -Force
    }
    Move-Item -LiteralPath $new -Destination $dest -Force
    Remove-Item -LiteralPath $old -Force -ErrorAction SilentlyContinue
  } catch { exit 1 }
}
`

// swapStagedWhileRunning moves dest -> dest.old then dest.new -> dest while this process still runs the old image.
func swapStagedWhileRunning(dest string) error {
	staging := dest + ".new"
	if _, err := os.Stat(staging); err != nil {
		return nil
	}
	backup := dest + ".old"
	_ = os.Remove(backup)
	if err := os.Rename(dest, backup); err != nil {
		return fmt.Errorf("rename running binary: %w", err)
	}
	if err := os.Rename(staging, dest); err != nil {
		_ = os.Rename(backup, dest)
		return err
	}
	return nil
}

func cleanupLeftoverWindowsBackup(dest string) {
	if dest == "" {
		return
	}
	if _, err := os.Stat(dest + ".new"); err == nil {
		return
	}
	_ = os.Remove(dest + ".old")
}

func windowsExecutablePath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	to, err := filepath.Abs(exe)
	if err != nil {
		to = exe
	}
	if resolved, err := filepath.EvalSymlinks(to); err == nil {
		to = resolved
	}
	return to, nil
}

func stagedWindowsPaths(dest string) (to, from string) {
	if f := os.Getenv(envStagedFrom); f != "" {
		return os.Getenv(envStagedTo), f
	}
	if dest == "" {
		return "", ""
	}
	from = dest + ".new"
	if _, err := os.Stat(from); err != nil {
		return "", ""
	}
	return dest, from
}
