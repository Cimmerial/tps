//go:build windows

package jobs

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func spawnInTerminal(workDir, shellCmd, pidFile, logFile, ttyFile string, focus bool) error {
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
	_ = os.Remove(ttyFile)

	base := strings.TrimSuffix(pidFile, ".pid")
	script := base + "-run.ps1"
	launcher := base + "-run.cmd"
	if err := writeWindowsRunScript(script, workDir, shellCmd, logFile, pidFile); err != nil {
		return err
	}
	if err := writeWindowsRunCmd(launcher, script); err != nil {
		return err
	}

	if os.Getenv("TPS_NO_WT") != "1" {
		if wt, err := exec.LookPath("wt.exe"); err == nil {
			cmd := exec.Command(wt, "-w", "0", "new-tab", "--", "cmd.exe", "/k", launcher)
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	}
	// cmd /k + .cmd launcher avoids execution-policy blocks on -run.ps1 (and WT swallowing powershell flags).
	cmd := exec.Command("cmd", "/c", "start", "", "cmd.exe", "/k", launcher)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func writeWindowsRunCmd(path, ps1Path string) error {
	ps1Path = filepath.Clean(ps1Path)
	var b strings.Builder
	b.WriteString("@echo off\r\n")
	b.WriteString("powershell.exe -NoProfile -NoExit -ExecutionPolicy Bypass -File \"")
	b.WriteString(ps1Path)
	b.WriteString("\"\r\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func writeWindowsRunScript(path, workDir, shellCmd, logFile, pidFile string) error {
	var b strings.Builder
	b.WriteString("$ErrorActionPreference = 'Continue'\n")
	b.WriteString("Set-Location -LiteralPath ")
	b.WriteString(psSingleQuote(workDir))
	b.WriteString("\n")
	inner := shellCmd
	if logFile != "" {
		inner = shellCmd + " *>&1 | Tee-Object -FilePath " + psSingleQuote(logFile)
	}
	// Run the project windows_shell line under cmd so ExecutionPolicy Bypass on that line is honored.
	b.WriteString("$jobCmd = ")
	b.WriteString(psSingleQuote(inner))
	b.WriteString("\n")
	b.WriteString("$p = Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', $jobCmd -PassThru -WorkingDirectory (Get-Location)\n")
	b.WriteString("Set-Content -LiteralPath ")
	b.WriteString(psSingleQuote(pidFile))
	b.WriteString(" -Value $p.Id\n")
	b.WriteString("Wait-Process -Id $p.Id\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func psSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
