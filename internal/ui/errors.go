package ui

import "errors"

// ErrPromptCancelled is returned when the user presses Esc on a cancelable menu prompt.
var ErrPromptCancelled = errors.New("prompt cancelled")
