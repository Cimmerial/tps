//go:build darwin

package jobs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func spawnInTerminal(workDir, shellCmd, pidFile, logFile, ttyFile string) error {
	workDir = filepath.Clean(workDir)
	pidFile = filepath.Clean(pidFile)
	ttyFile = filepath.Clean(ttyFile)
	pgidFile := strings.TrimSuffix(pidFile, ".pid") + ".pgid"
	if err := os.MkdirAll(filepath.Dir(pidFile), 0o700); err != nil {
		return err
	}
	if logFile != "" {
		logFile = filepath.Clean(logFile)
		if err := os.MkdirAll(filepath.Dir(logFile), 0o700); err != nil {
			return err
		}
	}
	_ = os.Remove(ttyFile)
	_ = os.Remove(pgidFile)

	base := strings.TrimSuffix(pidFile, ".pid")
	runScript := base + "-run.sh"
	launcher := base + ".command"

	if err := writeRunScript(runScript, workDir, shellCmd, logFile); err != nil {
		return err
	}

	runInv := "bash -l " + shSingleQuote(runScript)

	closeHook := fmt.Sprintf(`
tps_close_tab() {
  local dev
  if [[ -f %s ]]; then dev=$(cat %s); else dev=$(basename "$(tty 2>/dev/null)" 2>/dev/null); fi
  [[ -z "$dev" ]] && return 0
  osascript -e "tell application \"Terminal\"
    set target to \"$dev\"
    repeat with w in windows
      repeat with t in tabs of w
        try
          set tt to tty of t
          if tt is target or tt is (\"/dev/\" & target) then
            close t saving no
            return
          end if
        end try
      end repeat
    end repeat
  end tell" 2>/dev/null || true
}
trap tps_close_tab EXIT
dev=$(basename "$(tty 2>/dev/null)" 2>/dev/null)
echo "$dev" > %s
`, shSingleQuote(ttyFile), shSingleQuote(ttyFile), shSingleQuote(ttyFile))

	body := fmt.Sprintf(
		"#!/usr/bin/env bash\n%s\nset -m\n%s &\nchild=$!\necho $child > %s\npgid=$(ps -o pgid= -p $child 2>/dev/null | tr -d ' ')\necho \"$pgid\" > %s\nwait $child\n",
		closeHook,
		runInv,
		shSingleQuote(pidFile),
		shSingleQuote(pgidFile),
	)
	if err := os.WriteFile(launcher, []byte(body), 0o755); err != nil {
		return err
	}
	cmd := exec.Command("open", "-g", "-a", "Terminal", launcher)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeRunScript(path, workDir, shellCmd, logFile string) error {
	var b strings.Builder
	b.WriteString("#!/usr/bin/env bash\nset -e\n")
	b.WriteString("cd ")
	b.WriteString(shSingleQuote(workDir))
	b.WriteString(" || exit 1\n")
	if logFile != "" {
		b.WriteString("export TPS_CAPTURE_LOG=")
		b.WriteString(shSingleQuote(logFile))
		b.WriteString("\n")
		b.WriteString("if [[ -n \"$TPS_CAPTURE_LOG\" ]]; then exec > >(tee \"$TPS_CAPTURE_LOG\") 2>&1; fi\n")
	}
	b.WriteString(strings.TrimSpace(shellCmd))
	b.WriteString("\n")
	return os.WriteFile(path, []byte(b.String()), 0o755)
}

func shSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
