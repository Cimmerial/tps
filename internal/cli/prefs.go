package cli

import (
	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func applyUIPreferences() error {
	prefs, err := config.LoadPreferences()
	if err != nil {
		return err
	}
	ui.ApplySettings(ui.UISettings{
		LogoID:         prefs.LogoID(),
		TitleCardColor: prefs.TitleCardColorMode(),
		Accent:         prefs.AccentMode(),
		ColorEnabled:   prefs.ColorEnabled,
	})
	return nil
}

func refreshUISession(projectID string, statusTick int) error {
	return refreshInteractiveSession(projectID, statusTick, "")
}
