//go:build windows

package update

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// RestartInteractive exits and re-runs tps (applies any staged .new binary first).
func RestartInteractive(updated bool) error {
	bin, argv, err := buildRestartArgv()
	if err != nil {
		return fmt.Errorf("locate tps: %w", err)
	}
	if updated {
		argv = append(argv, flagUpdated)
	} else {
		argv = append(argv, flagUpToDate)
	}

	parent := os.Getpid()
	script, err := writeWindowsRestartScript(parent, bin, argv)
	if err != nil {
		return err
	}

	cmd := exec.Command("powershell.exe", "-NoProfile", "-ExecutionPolicy", "Bypass", "-WindowStyle", "Hidden", "-File", script)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start updater: %w", err)
	}
	os.Exit(0)
	return nil
}

func writeWindowsRestartScript(parent int, bin string, argv []string) (string, error) {
	tmp, err := os.MkdirTemp("", "tps-restart-*")
	if err != nil {
		return "", err
	}
	script := filepath.Join(tmp, "restart.ps1")

	var b strings.Builder
	b.WriteString("$ErrorActionPreference = 'Stop'\n")
	b.WriteString(fmt.Sprintf("$parent = %d\n", parent))
	b.WriteString("Wait-Process -Id $parent -ErrorAction SilentlyContinue\n")
	b.WriteString("Start-Sleep -Milliseconds 400\n")
	for _, pair := range pendingWindowsReplacements() {
		b.WriteString(fmt.Sprintf("Move-Item -LiteralPath %s -Destination %s -Force\n", psLiteral(pair.Staging), psLiteral(pair.Dest)))
	}
	b.WriteString(fmt.Sprintf("$bin = %s\n", psLiteral(bin)))
	argParts := make([]string, 0, len(argv)-1)
	for _, a := range argv[1:] {
		argParts = append(argParts, psLiteral(a))
	}
	if len(argParts) == 0 {
		b.WriteString("Start-Process -FilePath $bin -NoNewWindow\n")
	} else {
		b.WriteString(fmt.Sprintf("$args = @(%s)\n", strings.Join(argParts, ", ")))
		b.WriteString("Start-Process -FilePath $bin -ArgumentList $args -NoNewWindow\n")
	}

	if err := os.WriteFile(script, []byte(b.String()), 0o600); err != nil {
		return "", err
	}
	return script, nil
}

func psLiteral(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
