//go:build !windows

package jobs

import "os/exec"

func killPort(port string) bool {
	script := "lsof -ti tcp:" + port + " 2>/dev/null | xargs kill -9 2>/dev/null"
	err := exec.Command("sh", "-c", script).Run()
	return err == nil
}
