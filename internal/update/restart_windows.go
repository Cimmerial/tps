//go:build windows

package update

import (
	"fmt"
	"os"
	"strings"

	"github.com/Cimmerial/tps/internal/ui"
)

// RestartInteractive on Windows exits to the shell so the next `tps` starts with a clean console.
func RestartInteractive(updated bool, newVersion string) error {
	exitWindowsAfterUpdate(updated, newVersion)
	return nil
}

func exitWindowsAfterUpdate(updated bool, newVersion string) {
	ui.RestoreStdin()
	fmt.Print("\r\n")
	if updated {
		ScheduleWindowsStagedApplyAfterExit()
		ver := formatReleaseVersion(newVersion)
		fmt.Println(ui.Green("tps updated to "+ver+" — run ") + ui.Cyan("tps") + ui.Green(" again to continue."))
	}
	os.Exit(0)
}

func formatReleaseVersion(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "the new version"
	}
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}
