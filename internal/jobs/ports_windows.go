//go:build windows

package jobs

import "os/exec"

func killPort(port string) bool {
	script := "Get-NetTCPConnection -LocalPort " + port + " -ErrorAction SilentlyContinue | ForEach-Object { Stop-Process -Id $_.OwningProcess -Force -ErrorAction SilentlyContinue }"
	err := exec.Command("powershell", "-NoProfile", "-Command", script).Run()
	return err == nil
}
