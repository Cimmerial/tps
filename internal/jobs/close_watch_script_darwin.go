//go:build darwin

package jobs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// spawnCloseWatcherScript starts a nohup bash process that waits for pid then closes the Terminal tab.
func spawnCloseWatcherScript(projectID string, job Job) error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	scriptDir := filepath.Join(dir, "tps", "jobs", "watchers")
	if err := os.MkdirAll(scriptDir, 0o700); err != nil {
		return err
	}
	scriptPath := filepath.Join(scriptDir, fmt.Sprintf("%s-%d-close.sh", projectID, job.ID))
	marker := JobWindowMarker(projectID, job.ID)
	ttyPath := mustTTYFile(projectID, job.ID)

	body := fmt.Sprintf(`#!/usr/bin/env bash
set -u
pid=%d
marker=%s
while kill -0 "$pid" 2>/dev/null; do sleep 0.2; done
tty=""
if [[ -f %s ]]; then tty=$(tr -d '[:space:]' < %s); fi
/usr/bin/osascript - "$tty" "$marker" <<'TPS_WATCH_CLOSE'
on run argv
  set targetTTY to item 1 of argv
  set marker to item 2 of argv
  tell application "Terminal"
    if marker is not "" then
      repeat with w in windows
        repeat with t in tabs of w
          try
            set ct to custom title of t
            set tn to name of t
            if ct contains marker or tn contains marker then
              close t saving no
              return
            end if
          end try
        end repeat
      end repeat
    end if
    if targetTTY is not "" then
      set shortTTY to targetTTY
      if not (shortTTY starts with "/dev/") then set shortTTY to "/dev/" & shortTTY
      repeat with w in windows
        repeat with t in tabs of w
          try
            set tt to tty of t
            if tt is shortTTY then
              close t saving no
              return
            end if
          end try
        end repeat
      end repeat
    end if
  end tell
end run
TPS_WATCH_CLOSE
rm -f %s
`, job.PID, shSingleQuote(marker), shSingleQuote(ttyPath), shSingleQuote(ttyPath), shSingleQuote(scriptPath))

	if err := os.WriteFile(scriptPath, []byte(body), 0o755); err != nil {
		return err
	}
	// Detached bash+nohup: Go's exec.Command("nohup", path) often fails to keep the watcher alive on macOS.
	cmd := exec.Command("/bin/bash", "-c", "/usr/bin/nohup "+shSingleQuote(scriptPath)+" </dev/null >/dev/null 2>&1 &")
	return cmd.Start()
}

func mustTTYFile(projectID string, jobID int) string {
	p, err := TTYFile(projectID, jobID)
	if err != nil {
		return ""
	}
	return p
}
