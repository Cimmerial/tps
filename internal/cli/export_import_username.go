package cli

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/ui"
)

func promptIncludeUsername(action string, username string) (bool, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return false, nil
	}
	fmt.Println()
	fmt.Printf("    %s %s\n", ui.Dim("username:"), ui.Bold(username))
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "1", Abbrev: "y", Desc: ui.Dim(fmt.Sprintf("include username when %sing options", action))},
		{Num: "2", Abbrev: "n", Desc: ui.Dim(fmt.Sprintf("omit username from %s", action))},
	})
	fmt.Println()
	for {
		line, err := ui.ReadMenuLine()
		if err != nil {
			return false, err
		}
		switch strings.ToLower(strings.TrimSpace(line)) {
		case "1", "y", "yes":
			return true, nil
		case "2", "n", "no":
			return false, nil
		default:
			fmt.Println(ui.Dim("enter 1, 2, y, or n"))
		}
	}
}

func exportIncludeUsername(prefsUsername string) (bool, error) {
	if strings.TrimSpace(prefsUsername) == "" {
		return false, nil
	}
	return promptIncludeUsername("export", prefsUsername)
}

func importIncludeUsername(fileUsername string) (bool, error) {
	if strings.TrimSpace(fileUsername) == "" {
		return false, nil
	}
	return promptIncludeUsername("import", fileUsername)
}
