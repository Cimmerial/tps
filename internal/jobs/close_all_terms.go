package jobs

// CloseAllTermsResult reports how many Terminal.app tabs were closed and their labels.
type CloseAllTermsResult struct {
	Count  int
	Labels []string
}
