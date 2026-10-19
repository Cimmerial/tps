package jobs

import (
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func childPIDs(pid int) []int {
	if pid <= 0 {
		return nil
	}
	out, err := exec.Command("pgrep", "-P", strconv.Itoa(pid)).Output()
	if err != nil {
		return nil
	}
	var kids []int
	for _, line := range strings.Fields(string(out)) {
		if n, err := strconv.Atoi(line); err == nil && n > 0 {
			kids = append(kids, n)
		}
	}
	return kids
}

func processGroup(pid int) int {
	out, err := exec.Command("ps", "-o", "pgid=", "-p", strconv.Itoa(pid)).Output()
	if err != nil {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	return n
}

func collectTree(root int) []int {
	if root <= 0 {
		return nil
	}
	seen := map[int]bool{}
	queue := []int{root}
	for len(queue) > 0 {
		pid := queue[0]
		queue = queue[1:]
		if pid <= 0 || seen[pid] || !alive(pid) {
			continue
		}
		seen[pid] = true
		for _, c := range childPIDs(pid) {
			queue = append(queue, c)
		}
	}
	out := make([]int, 0, len(seen))
	for pid := range seen {
		out = append(out, pid)
	}
	return out
}

func killProcessTree(root int) (before, after []int) {
	if root <= 0 {
		return nil, nil
	}
	before = collectTree(root)
	pgid := processGroup(root)
	if pgid > 0 {
		_ = syscall.Kill(-pgid, syscall.SIGTERM)
	}
	killSubtree(root)
	time.Sleep(250 * time.Millisecond)
	for _, pid := range before {
		if alive(pid) {
			_ = syscall.Kill(pid, syscall.SIGKILL)
		}
	}
	if pgid > 0 {
		_ = syscall.Kill(-pgid, syscall.SIGKILL)
	}
	time.Sleep(100 * time.Millisecond)
	after = survivors(before)
	return before, after
}

func survivors(from []int) []int {
	seen := map[int]bool{}
	var out []int
	add := func(pid int) {
		if pid <= 0 || seen[pid] || !alive(pid) {
			return
		}
		seen[pid] = true
		out = append(out, pid)
	}
	for _, p := range from {
		add(p)
		for _, c := range childPIDs(p) {
			add(c)
			for _, gc := range childPIDs(c) {
				add(gc)
			}
		}
	}
	return out
}

func killSubtree(pid int) {
	for _, c := range childPIDs(pid) {
		killSubtree(c)
	}
	if alive(pid) {
		_ = syscall.Kill(pid, syscall.SIGTERM)
	}
}
