//go:build linux

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
	launcher := base + ".sh"

	if err := writeRunScript(runScript, workDir, shellCmd, logFile); err != nil {
		return err
	}

	body := fmt.Sprintf(
		"#!/usr/bin/env bash\nset -m\nbash -l %s &\nchild=$!\necho $child > %s\npgid=$(ps -o pgid= -p $child 2>/dev/null | tr -d ' ')\necho \"$pgid\" > %s\nwait $child\n",
		shSingleQuote(runScript),
		shSingleQuote(pidFile),
		shSingleQuote(pgidFile),
	)
	if err := os.WriteFile(launcher, []byte(body), 0o755); err != nil {
		return err
	}

	term, args := linuxTerminalLauncher(launcher)
	if term == "" {
		return fmt.Errorf("no supported terminal emulator found (install wezterm, kitty, gnome-terminal, or xterm)")
	}
	cmd := exec.Command(term, args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func linuxTerminalLauncher(script string) (string, []string) {
	candidates := []struct {
		bin  string
		args func(string) []string
	}{
		{"wezterm", func(s string) []string { return []string{"start", "--", "bash", "-l", s} }},
		{"kitty", func(s string) []string { return []string{"bash", "-l", s} }},
		{"gnome-terminal", func(s string) []string { return []string{"--", "bash", "-l", s} }},
		{"x-terminal-emulator", func(s string) []string { return []string{"-e", "bash", "-l", s} }},
		{"xterm", func(s string) []string { return []string{"-e", "bash", "-l", s} }},
	}
	for _, c := range candidates {
		if path, err := exec.LookPath(c.bin); err == nil {
			return path, c.args(script)
		}
	}
	return "", nil
}
