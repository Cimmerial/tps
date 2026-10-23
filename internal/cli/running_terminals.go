package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/Cimmerial/tps/internal/jobs"
	"github.com/Cimmerial/tps/internal/ui"
)

func formatRunningTerminalChip(j jobs.Job) string {
	age := time.Since(j.Started).Round(time.Second)
	return ui.Bold(j.Alias) + ui.Dim(fmt.Sprintf("  pid %d  %s", j.PID, age))
}

func buildRunningTerminalsLine(list []jobs.Job) string {
	if len(list) == 0 {
		return ""
	}
	sep := ui.Dim(" · ")
	parts := make([]string, len(list))
	for i, j := range list {
		parts[i] = formatRunningTerminalChip(j)
	}
	return strings.Join(parts, sep)
}

func runningTerminalsTitle(n int) string {
	word := "TERMINAL"
	if n != 1 {
		word += "S"
	}
	return ui.Blue(fmt.Sprintf("%d", n)) + " " + ui.HeaderInfo("RUNNING "+word+":")
}

func printRunningTerminalsSection(projectID string) {
	list, err := jobs.Reap(projectID)
	if err != nil || len(list) == 0 {
		return
	}
	line := buildRunningTerminalsLine(list)
	box := ui.NewMenuBox(ui.ChromeSession)
	box.Title(runningTerminalsTitle(len(list)))
	box.Plain(ui.SectionTrail{LastSection: true, LastRow: true}, line)
	box.Flush()
}

func printRunningTerminalsDetail(projectID string) {
	list, err := jobs.Reap(projectID)
	if err != nil {
		return
	}
	if len(list) == 0 {
		fmt.Println(ui.Dim("no terminals running (tk · tk -a · t -k 1)"))
		return
	}
	var rows []numberedMenuRow
	for i, j := range list {
		age := time.Since(j.Started).Round(time.Second)
		rows = append(rows, numberedMenuRow{
			Desc:     fmt.Sprintf("%d) %s", i+1, j.Alias),
			Abbrev:   fmt.Sprintf("pid %d", j.PID),
			HelpLine: age.String() + "  " + truncateRunes(j.Shell, 48),
		})
	}
	printInfoMenuPanelTitle(runningTerminalsTitle(len(list)), []menuPanelSection{{Rows: rows}})
	fmt.Println(ui.Dim("tk · tk -a · t -k 1 · tr"))
}
