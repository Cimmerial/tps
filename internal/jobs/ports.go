package jobs

import (
	"os/exec"
	"regexp"
)

var portFlagRE = regexp.MustCompile(`(?:--port|-p)\s+(\d+)`)

func killPortsFromShell(shell string) []string {
	seen := map[string]bool{}
	var freed []string
	for _, m := range portFlagRE.FindAllStringSubmatch(shell, -1) {
		if len(m) < 2 || seen[m[1]] {
			continue
		}
		seen[m[1]] = true
		if killPort(m[1]) {
			freed = append(freed, m[1])
		}
	}
	return freed
}

func killPort(port string) bool {
	script := "lsof -ti tcp:" + port + " 2>/dev/null | xargs kill -9 2>/dev/null"
	err := exec.Command("sh", "-c", script).Run()
	return err == nil
}

