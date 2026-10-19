//go:build !darwin

package jobs

func CloseAllTerminalWindows() CloseAllTermsResult {
	return CloseAllTermsResult{}
}
