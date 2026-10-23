package cli

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func optionsUsernameRowDesc(prefs *config.Preferences) string {
	base := ui.Dim("username (export filenames)")
	if prefs == nil {
		return base
	}
	u := strings.TrimSpace(prefs.Username)
	if u == "" {
		return base + ui.Dim(" — not set")
	}
	return base + ui.Dim(" — ") + ui.Bold(u)
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
