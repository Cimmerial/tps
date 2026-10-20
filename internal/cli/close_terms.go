package cli

import (
	"fmt"
	"runtime"

	"github.com/Cimmerial/tps/internal/jobs"
	"github.com/Cimmerial/tps/internal/ui"
)

func runCloseAllTerms() (jobs.CloseAllTermsResult, error) {
	label := "closing terminals…"
	if runtime.GOOS != "darwin" {
		label = "stopping terminals…"
	}
	var result jobs.CloseAllTermsResult
	err := runWithBraille(label, func() error {
		result = jobs.CloseAllTerminalWindows()
		return nil
	})
	return result, err
}

func printCloseAllTermsDone(result jobs.CloseAllTermsResult) {
	if runtime.GOOS == "darwin" {
		if result.Count == 0 {
			fmt.Println(ui.Green("closed 0 Terminal.app windows"))
			return
		}
		fmt.Println(ui.Green(fmt.Sprintf("closed %d Terminal.app window(s)", result.Count)))
	} else {
		if result.Count == 0 {
			fmt.Println(ui.Green("stopped 0 tps terminal job(s)"))
			return
		}
		fmt.Println(ui.Green(fmt.Sprintf("stopped %d tps terminal job(s)", result.Count)))
	}
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
