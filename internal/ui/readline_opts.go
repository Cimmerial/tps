package ui

// ReadLineOptions configures ReadEditedLineFuncHistoryWithOpts.
type ReadLineOptions struct {
	Boxed                bool
	CancelOnEsc          bool
	ResponseRuleOnSubmit bool
}

// DefaultReadLineOptions is single-line input with no special behavior.
func DefaultReadLineOptions() ReadLineOptions {
	return ReadLineOptions{}
}

// ReadMenuLine reads a boxed choice prompt (red chrome, Esc = cancel).
func ReadMenuLine() (string, error) {
	SetPromptChrome(ChromeChoice)
	return ReadEditedLineFuncHistoryWithOpts(func() string { return Prompt() }, nil, ReadLineOptions{
		CancelOnEsc: true,
	})
}
