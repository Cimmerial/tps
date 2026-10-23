package cli

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func optionsOutputRowDesc(prefs *config.Preferences) string {
	return optionsTopMenuItemDesc("command output", outputModeLabel(prefs.SimplifiedOutput))
}

func outputModeLabel(simplified bool) string {
	if simplified {
		return "simplified"
	}
	return "verbose"
}

func optionsOutputMenu(prefs *config.Preferences, projectID string, statusTick int) (bool, error) {
	ui.ClearScreen()
	current := prefs.SimplifiedOutput
	fmt.Println()
	fmt.Printf("%s\n", ui.HeaderChoice("options"))
	fmt.Println()
	fmt.Printf("%s\n", ui.SubHeader("command output"))
	fmt.Printf("    %s %s\n", ui.Dim("current:"), ui.Bold(outputModeLabel(current)))
	fmt.Println()
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "1", Desc: optionsRowDescPlain("Verbose", !current), HelpLine: "full git, test, and script output"},
		{Num: "2", Desc: optionsRowDescPlain("Simplified", current), HelpLine: "key results only (ship, release, tests)"},
	})
	fmt.Println()
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "0", Abbrev: "q", Desc: ui.Dim("no change / back")},
	})
	fmt.Println()
	choice, back, err := readOptionsMenuChoice()
	if err != nil || back {
		return false, err
	}
	var next bool
	switch choice {
	case 1:
		next = false
	case 2:
		next = true
	default:
		fmt.Println(ui.Dim("invalid choice"))
		return false, nil
	}
	if prefs.SimplifiedOutput == next {
		fmt.Println(ui.Dim("already set"))
		return false, nil
	}
	prefs.SimplifiedOutput = next
	if err := config.SavePreferences(prefs); err != nil {
		return false, err
	}
	if err := applyUIPreferences(); err != nil {
		return false, err
	}
	fmt.Println(ui.Dim(fmt.Sprintf("command output: %s", outputModeLabel(next))))
	if err := refreshUISession(projectID, statusTick); err != nil {
		return true, err
	}
	return true, nil
}
