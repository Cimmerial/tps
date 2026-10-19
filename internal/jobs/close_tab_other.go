//go:build !darwin

package jobs

func tryCloseDetachedTab(tty, marker string) {}
