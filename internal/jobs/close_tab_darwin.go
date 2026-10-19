//go:build darwin

package jobs

import "time"

func tryCloseDetachedTab(tty, marker string) {
	for i := 0; i < 5; i++ {
		if closeTerminalTabOK(tty, marker) {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
}
