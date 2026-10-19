package cli

import (
	"fmt"
	"sync"
	"time"

	"github.com/Cimmerial/tps/internal/jobs"
	"github.com/Cimmerial/tps/internal/ui"
)

var (
	detachedNoticeMu sync.Mutex
	detachedNotices  []string
)

func initDetachedCompleteNotices() {
	jobs.OnDetachedJobComplete(func(job jobs.Job, elapsed time.Duration) {
		secs := elapsed.Seconds()
		line := ui.Dim(fmt.Sprintf("%s process completed after %.1fs", job.Alias, secs))
		if !ui.PrintLinesAbovePrompt(line) {
			detachedNoticeMu.Lock()
			detachedNotices = append(detachedNotices, line)
			detachedNoticeMu.Unlock()
		}
	})
}

func drainDetachedCompleteNotices() {
	detachedNoticeMu.Lock()
	lines := detachedNotices
	detachedNotices = nil
	detachedNoticeMu.Unlock()
	for _, line := range lines {
		fmt.Println(line)
	}
}
