package jobs

import "regexp"

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

