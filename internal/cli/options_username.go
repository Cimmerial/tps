package cli

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func optionsUsernameRowDesc(prefs *config.Preferences) string {
	if prefs == nil {
		return optionsTopMenuItemDesc("username (export filenames)", "not set")
	}
	u := strings.TrimSpace(prefs.Username)
	if u == "" {
		return optionsTopMenuItemDesc("username (export filenames)", "not set")
	}
	return optionsTopMenuItemDesc("username (export filenames)", u)
}

func optionsPromptUsername(prefs *config.Preferences) error {
	fmt.Println()
	fmt.Printf("    %s %s\n", ui.Dim("current:"), usernameDisplay(prefs.Username))
	line, err := ui.ReadEditedLine("username: ")
	if err != nil {
		return err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		fmt.Println(ui.Dim("no change"))
		return nil
	}
	prefs.Username = line
	if err := config.SavePreferences(prefs); err != nil {
		return err
	}
	fmt.Println(ui.Green("saved username"))
	return nil
}

func usernameDisplay(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ui.Dim("(not set)")
	}
	return ui.Bold(v)
}
