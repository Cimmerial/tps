package ui

import "sync"

type promptNoticeSession struct {
	mu     sync.Mutex
	insert func(lines ...string)
}

var (
	promptNoticeMu      sync.Mutex
	activePromptNotice  *promptNoticeSession
)

// BeginPromptNotices registers a handler that prints lines above the active edited prompt.
// Call the returned function when ReadEditedLine returns.
func BeginPromptNotices(insert func(lines ...string)) func() {
	s := &promptNoticeSession{insert: insert}
	promptNoticeMu.Lock()
	activePromptNotice = s
	promptNoticeMu.Unlock()
	return func() {
		promptNoticeMu.Lock()
		if activePromptNotice == s {
			activePromptNotice = nil
		}
		promptNoticeMu.Unlock()
	}
}

// PrintLinesAbovePrompt prints immediately when a prompt session is active.
// Returns false if nothing was listening (caller should queue for later).
func PrintLinesAbovePrompt(lines ...string) bool {
	if len(lines) == 0 {
		return true
	}
	promptNoticeMu.Lock()
	s := activePromptNotice
	promptNoticeMu.Unlock()
	if s == nil || s.insert == nil {
		return false
	}
	s.mu.Lock()
	s.insert(lines...)
	s.mu.Unlock()
	return true
}
