//go:build darwin

package jobs

import (
	"fmt"
	"os/exec"
	"strings"
)

func CloseTerminalTab(ttyDev, windowMarker string) {
	short := normalizeTTY(ttyDev)
	if short != "" {
		script := ttyMatchScript(short, true)
		out, err := exec.Command("osascript", "-e", script).Output()
		if err == nil && strings.TrimSpace(string(out)) == "closed tab" {
			return
		}
	}
	if windowMarker != "" {
		closeWindowsByMarker(windowMarker)
	}
}

func ttyMatchScript(shortTTY string, close bool) string {
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
	esc := strings.ReplaceAll(marker, "\"", "\\\"")
	script := fmt.Sprintf(`tell application "Terminal"
  repeat with w in windows
    set wn to name of w
    if wn contains "%s" then
      close w saving no
      return
    end if
    repeat with t in tabs of w
      try
        set ct to custom title of t
        set tn to name of t
        if ct contains "%s" or tn contains "%s" then
          close t saving no
          return
        end if
      end try
    end repeat
  end repeat
end tell`, esc, esc, esc)
	_ = exec.Command("osascript", "-e", script).Run()
}
