package cli

import (
	"fmt"
	"runtime"

	"github.com/Cimmerial/tps/internal/jobs"
	"github.com/Cimmerial/tps/internal/ui"
)

func runCloseAllTerms() (jobs.CloseAllTermsResult, error) {
	if runtime.GOOS != "darwin" {
		return jobs.CloseAllTermsResult{}, fmt.Errorf("vite: macOS only (Terminal.app)")
	}
	var result jobs.CloseAllTermsResult
	err := runWithBraille("closing terminals…", func() error {
		result = jobs.CloseAllTerminalWindows()
		return nil
	})
	return result, err
}

func printCloseAllTermsDone(result jobs.CloseAllTermsResult) {
	if result.Count == 0 {
		fmt.Println(ui.Green("closed 0 Terminal.app windows"))
		return
	}
	fmt.Println(ui.Green(fmt.Sprintf("closed %d Terminal.app window(s)", result.Count)))
	for i, label := range result.Labels {
		if i >= 5 {
			break
		}
		fmt.Printf("  %s\n", ui.Dim(label))
	}
	if result.Count > 5 {
		fmt.Println(ui.Dim(fmt.Sprintf("  … and %d more", result.Count-5)))
	}
}
