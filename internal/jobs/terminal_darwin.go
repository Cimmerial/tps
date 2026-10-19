//go:build darwin

package jobs

import (
	"fmt"
	"os/exec"
	"strings"
)

func CloseTerminalTab(ttyDev, windowMarker string) {
	closeTerminalTabOK(ttyDev, windowMarker)
}

func closeTerminalTabOK(ttyDev, windowMarker string) bool {
	short := normalizeTTY(ttyDev)
	if short != "" {
		script := ttyMatchScript(short)
		out, err := exec.Command("osascript", "-e", script).Output()
		if err == nil && strings.TrimSpace(string(out)) == "closed tab" {
			return true
		}
	}
	if windowMarker != "" {
		out, err := exec.Command("osascript", "-e", closeTabByMarkerScript(windowMarker)).Output()
		if err == nil {
			switch strings.TrimSpace(string(out)) {
			case "closed tab", "closed window":
				return true
			}
		}
	}
	return false
}

func closeTabByMarkerScript(marker string) string {
	esc := strings.ReplaceAll(marker, "\"", "\\\"")
	return fmt.Sprintf(`tell application "Terminal"
  repeat with w in windows
    set wn to name of w
    if wn contains "%s" then
      close w saving no
      return "closed window"
    end if
    repeat with t in tabs of w
      try
        set ct to custom title of t
        set tn to name of t
        if ct contains "%s" or tn contains "%s" then
          close t saving no
          return "closed tab"
        end if
      end try
    end repeat
  end repeat
end tell
return "no tab matched marker"`, esc, esc, esc)
}

func ttyMatchScript(shortTTY string) string {
	esc := strings.ReplaceAll(shortTTY, "\"", "\\\"")
	return fmt.Sprintf(`tell application "Terminal"
  set target to "%s"
  repeat with w in windows
    repeat with t in tabs of w
      try
        set tt to tty of t
        if tt is target or tt is ("/dev/" & target) then
          close t saving no
          return "closed tab"
        end if
      end try
    end repeat
  end repeat
end tell
return "no tab matched"`, esc)
}

func closeWindowsByMarker(marker string) {
	_, _ = exec.Command("osascript", "-e", closeTabByMarkerScript(marker)).Output()
}
