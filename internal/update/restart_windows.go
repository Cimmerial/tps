//go:build windows

package update

import (
	"fmt"
	"os"

	"github.com/Cimmerial/tps/internal/ui"
)

// RestartInteractive on Windows exits to the shell so the next `tps` starts with a clean console.
func RestartInteractive(updated bool) error {
	exitWindowsAfterUpdate(updated)
	return nil
}

func exitWindowsAfterUpdate(updated bool) {
	ui.RestoreStdin()
	fmt.Print("\r\n")
	if updated {
		fmt.Println(ui.Green("tps updated — run ") + ui.Cyan("tps") + ui.Green(" again to continue."))
	}
	os.Exit(0)
}
