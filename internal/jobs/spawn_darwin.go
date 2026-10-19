//go:build darwin

package jobs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const closeTerminalTabHook = `
tps_close_tab() {
  local id
  id=$(basename "$(tty)" 2>/dev/null) || return 0
  osascript -e "tell application \"Terminal\"
    repeat with w in windows
      repeat with t in tabs of w
        try
          if tty of t ends with \"$id\" then
            close t
            return
          end if
        end try
      end repeat
    end repeat
  end tell" 2>/dev/null || true
}
trap tps_close_tab EXIT
`

func spawnInTerminal(workDir, shellCmd, pidFile, logFile string) error {
	workDir = filepath.Clean(workDir)
	pidFile = filepath.Clean(pidFile)
	if err := os.MkdirAll(filepath.Dir(pidFile), 0o700); err != nil {
		return err
	}
	if logFile != "" {
		logFile = filepath.Clean(logFile)
		if err := os.MkdirAll(filepath.Dir(logFile), 0o700); err != nil {
			return err
		}
	}

	base := strings.TrimSuffix(pidFile, ".pid")
	runScript := base + "-run.sh"
	launcher := base + ".command"

	if err := writeRunScript(runScript, workDir, shellCmd); err != nil {
		return err
	}

	runInv := "bash -l " + shSingleQuote(runScript)
	var runLine string
	if logFile != "" {
		runLine = fmt.Sprintf("( %s 2>&1 | tee %s )", runInv, shSingleQuote(logFile))
	} else {
		runLine = runInv
	}

	body := fmt.Sprintf(
		"#!/usr/bin/env bash\n%s\n%s &\nchild=$!\necho $child > %s\nwait $child\n",
		closeTerminalTabHook,
		runLine,
		shSingleQuote(pidFile),
	)
	if err := os.WriteFile(launcher, []byte(body), 0o755); err != nil {
		return err
	}
	cmd := exec.Command("open", "-g", "-a", "Terminal", launcher)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeRunScript(path, workDir, shellCmd string) error {
	var b strings.Builder
	b.WriteString("#!/usr/bin/env bash\nset -e\n")
	b.WriteString("cd ")
	b.WriteString(shSingleQuote(workDir))
	b.WriteString(" || exit 1\n")
	b.WriteString(strings.TrimSpace(shellCmd))
	b.WriteString("\n")
	return os.WriteFile(path, []byte(b.String()), 0o755)
}

func shSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
