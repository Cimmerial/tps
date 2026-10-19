package ui

import (
	"fmt"
	"os"
	"time"
)

// BrailleWait runs fn while animating label on stdout (safe before the next prompt).
func BrailleWait(label string, fn func()) {
	brailleWait(label, os.Stdout, stdoutIsTTY, fn)
}

// BrailleWaitStderr animates on stderr while fn runs (foreground commands writing to stdout).
func BrailleWaitStderr(label string, fn func()) {
	brailleWait(label, os.Stderr, stderrIsTTY, fn)
}

func brailleWait(label string, w *os.File, isTTY func() bool, fn func()) {
	if !isTTY() {
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
				fmt.Fprint(w, "\r\033[K")
				close(done)
				return
			default:
				fmt.Fprintf(w, "\r%s %s", BrailleFrame(i), Dim(label))
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
	fn()
	close(stop)
	<-done
}
