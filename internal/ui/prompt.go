package ui

import (
	"fmt"
	"os"
)

// ClearPromptLine removes the current input line (for async messages above the prompt).
func ClearPromptLine() {
	fmt.Fprint(os.Stdout, "\r\033[K")
}

// PrintCaptureResult prints a captured value on its own line (call while not at the prompt).
func PrintCaptureResult(value, copiedMsg string) {
	fmt.Println(value)
	if copiedMsg != "" {
		fmt.Println(copiedMsg)
	}
}
