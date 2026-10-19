//go:build darwin

package jobs

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Cimmerial/tps/internal/config"
)

const terminalTabFieldSep = "\x1f"

type terminalTab struct {
	TTY   string
	Label string
}

// CloseAllTerminalWindows stops processes in every Terminal.app tab, then closes tabs
// (same kill-then-close order as tk, but for all windows — not only tps-tracked jobs).
func CloseAllTerminalWindows() CloseAllTermsResult {
	cfg, _ := config.Load()
	tracked := RegisteredJobLabels(cfg)
	before := listTerminalTabs()
	_, _ = KillAllRegistered(cfg)

	closeEveryTerminalTab(before)
	return combineCloseResults(tracked, closedTabsSince(before, listTerminalTabs()))
}

func combineCloseResults(tracked []string, tabs CloseAllTermsResult) CloseAllTermsResult {
	if tabs.Count > 0 {
		return tabs
	}
	if len(tracked) == 0 {
		return tabs
	}
	return CloseAllTermsResult{Count: len(tracked), Labels: tracked}
}

func closedTabsSince(before, after []terminalTab) CloseAllTermsResult {
	still := map[string]bool{}
	for _, t := range after {
		key := normalizeTTY(t.TTY)
		if key != "" {
			still[key] = true
		}
	}
	var labels []string
	for _, t := range before {
		key := normalizeTTY(t.TTY)
		if key == "" || still[key] {
			continue
		}
		labels = append(labels, t.displayLabel())
	}
	return CloseAllTermsResult{Count: len(labels), Labels: labels}
}

func (t terminalTab) displayLabel() string {
	label := strings.TrimSpace(t.Label)
	if label != "" {
		return label
	}
	short := normalizeTTY(t.TTY)
	if short != "" {
		return short
	}
	return strings.TrimSpace(t.TTY)
}

func closeEveryTerminalTab(before []terminalTab) {
	ttys := tabTTYs(before)
	for _, tty := range ttys {
		killProcessesOnTTY(tty)
	}
	time.Sleep(300 * time.Millisecond)
	for _, tty := range ttys {
		closeTerminalTabOK(tty, "")
	}
	for _, tty := range listTerminalTabTTYs() {
		killProcessesOnTTY(tty)
		closeTerminalTabOK(tty, "")
	}
}

func tabTTYs(tabs []terminalTab) []string {
	out := make([]string, 0, len(tabs))
	for _, t := range tabs {
		if strings.TrimSpace(t.TTY) != "" {
			out = append(out, t.TTY)
		}
	}
	return out
}

func listTerminalTabTTYs() []string {
	return tabTTYs(listTerminalTabs())
}

func listTerminalTabs() []terminalTab {
	script := `tell application "Terminal"
  set delim to (ASCII character 31)
  set out to ""
  repeat with w in windows
    repeat with t in tabs of w
      try
        set out to out & (tty of t) & delim & (name of t) & linefeed
      end try
    end repeat
  end repeat
  return out
end tell`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return nil
	}
	var tabs []terminalTab
	for _, line := range strings.FieldsFunc(string(out), func(r rune) bool { return r == '\n' || r == '\r' }) {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tty, label, _ := strings.Cut(line, terminalTabFieldSep)
		tty = strings.TrimSpace(tty)
		if tty == "" {
			continue
		}
		tabs = append(tabs, terminalTab{TTY: tty, Label: strings.TrimSpace(label)})
	}
	return tabs
}

func killProcessesOnTTY(ttyDev string) {
	short := normalizeTTY(ttyDev)
	if short == "" {
		return
	}
	dev := "/dev/" + short
	raw, err := exec.Command("lsof", "-t", dev).Output()
	if err != nil {
		return
	}
	seen := map[int]bool{}
	for _, field := range strings.Fields(string(raw)) {
		pid, err := strconv.Atoi(field)
		if err != nil || pid <= 0 || seen[pid] {
			continue
		}
		seen[pid] = true
		_, _ = killProcessTree(pid)
	}
}
