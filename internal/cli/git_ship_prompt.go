package cli

import (
	"fmt"

	gittool "github.com/Cimmerial/tps/internal/git"
	"github.com/Cimmerial/tps/internal/ui"
)

func promptShipSuspiciousFiles(paths []string) (gittool.SuspiciousAction, error) {
	fmt.Println()
	fmt.Println(ui.Yellow("ship would commit suspicious files:"))
	for _, p := range paths {
		fmt.Printf("    %s\n", ui.Bold(p))
	}
	fmt.Println()
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "1", Abbrev: "y", Desc: "confirm and ship anyway"},
		{Num: "2", Abbrev: "n", Desc: "cancel"},
		{Num: "3", Abbrev: "", Desc: "commit without these files and add them to .gitignore"},
	})
	fmt.Println()
	for {
		line, err := ui.ReadEditedLine(ui.Prompt())
		if err != nil {
			return 0, err
		}
		if act, ok := gittool.ParseShipSuspiciousChoice(line); ok {
			return act, nil
		}
		fmt.Println(ui.Dim("enter 1, 2, 3, y, or n"))
	}
}

func runGitShip(dir, message string) (int, error) {
	return gittool.ShipWithOptions(dir, message, &gittool.ShipOptions{
		OnSuspicious: promptShipSuspiciousFiles,
	})
}
