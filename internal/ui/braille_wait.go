package ui

import (
	"fmt"
	"os"
	"time"
)

// BrailleWait runs fn while animating label on stdout (safe before the next prompt).
func BrailleWait(label string, fn func()) {
	if !stdoutIsTTY() {
		fn()
		return
	}
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		i := 0
		for {
			select {
			case <-stop:
				fmt.Fprint(os.Stdout, "\r\033[K")
				close(done)
				return
			default:
				fmt.Fprintf(os.Stdout, "\r%s %s", BrailleFrame(i), Dim(label))
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
	fn()
	close(stop)
	<-done
}
