//go:build windows

package jobs

import (
	"os/exec"
	"strconv"
	"time"
)

func childPIDs(pid int) []int {
	return nil
}

func processGroup(pid int) int {
	return 0
}

func collectTree(root int) []int {
	if root <= 0 {
		return nil
	}
	if alive(root) {
		return []int{root}
	}
	return nil
}

func killProcessTree(root int) (before, after []int) {
	if root <= 0 {
		return nil, nil
	}
	before = []int{root}
	_ = exec.Command("taskkill", "/PID", strconv.Itoa(root), "/T", "/F").Run()
	time.Sleep(200 * time.Millisecond)
	after = survivors(before)
	return before, after
}

func survivors(from []int) []int {
	var out []int
	for _, pid := range from {
		if alive(pid) {
			out = append(out, pid)
		}
	}
	return out
}

func killSubtree(pid int) {
	_ = exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F").Run()
}
