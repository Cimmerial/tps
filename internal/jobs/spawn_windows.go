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
	if err := writeWindowsRunScript(script, workDir, shellCmd, logFile, pidFile); err != nil {
		return err
	}

	if wt, err := exec.LookPath("wt.exe"); err == nil {
		cmd := exec.Command(wt, "-w", "0", "new-tab", "powershell", "-NoExit", "-NoProfile", "-File", script)
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	cmd := exec.Command("cmd", "/c", "start", "", "powershell", "-NoExit", "-NoProfile", "-File", script)
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
	b.WriteString("$arg = ")
	b.WriteString(psSingleQuote(inner))
	b.WriteString("\n")
	b.WriteString("$p = Start-Process powershell -ArgumentList '-NoProfile','-NoExit','-Command',$arg -PassThru\n")
	b.WriteString("Set-Content -LiteralPath ")
	b.WriteString(psSingleQuote(pidFile))
	b.WriteString(" -Value $p.Id\n")
	b.WriteString("Wait-Process -Id $p.Id\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func psSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
