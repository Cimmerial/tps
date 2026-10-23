//go:build darwin

package jobs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func spawnInTerminal(workDir, shellCmd, pidFile, logFile, ttyFile string, focus bool) error {
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

	recordTTY := fmt.Sprintf(
		"dev=$(basename \"$(tty 2>/dev/null)\" 2>/dev/null | tr -d '[:space:]')\nprintf '%%s' \"$dev\" > %s\n",
		shSingleQuote(ttyFile),
	)
	body := fmt.Sprintf(
		"#!/usr/bin/env bash\n%s\nset -m\n%s &\nchild=$!\necho $child > %s\npgid=$(ps -o pgid= -p $child 2>/dev/null | tr -d ' ')\necho \"$pgid\" > %s\nwait $child\n",
		recordTTY,
		runInv,
		shSingleQuote(pidFile),
		shSingleQuote(pgidFile),
	)
	if err := os.WriteFile(launcher, []byte(body), 0o755); err != nil {
		return err
	}
	openArgs := []string{"-a", "Terminal", launcher}
	if !focus {
		openArgs = append([]string{"-g"}, openArgs...)
	}
	cmd := exec.Command("open", openArgs...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
