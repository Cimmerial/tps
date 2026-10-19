package ui

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

var brailleFrames = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

type Spinner struct {
	stop chan struct{}
	done chan struct{}
}

func NewSpinner(label string) *Spinner {
	s := &Spinner{
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	if !stderrIsTTY() {
		close(s.done)
		return s
	}
	go func() {
		i := 0
		for {
			select {
			case <-s.stop:
				fmt.Fprint(os.Stderr, "\r\033[K")
				close(s.done)
				return
			default:
				frame := brailleFrames[i%len(brailleFrames)]
				fmt.Fprintf(os.Stderr, "\r%s %s", frame, label)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
	return s
}

func (s *Spinner) Stop() {
	if s.done == nil {
		return
	}
	select {
	case <-s.done:
		return
	default:
		close(s.stop)
		<-s.done
	}
}

func stderrIsTTY() bool {
	return term.IsTerminal(int(os.Stderr.Fd()))
}
