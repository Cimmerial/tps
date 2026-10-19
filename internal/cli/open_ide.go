package cli

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ide"
	"github.com/Cimmerial/tps/internal/ui"
)

func runOpenInIDE(p *config.Project) error {
	prefs, err := config.LoadPreferences()
	if err != nil {
		return err
	}
	choice, err := ide.Resolve(prefs.PreferredIDE)
	if err != nil {
		return err
	}
	there, err := ide.AlreadyThere(choice, p.Path)
	if err != nil {
		return err
	}
	if there {
		fmt.Println(ui.Dim(fmt.Sprintf("already in %s (%s)", choice.Profile.LabelOrID(), p.Name)))
		return nil
	}
	if err := ide.OpenProject(choice, p.Path); err != nil {
		return err
	}
	fmt.Println(ui.Dim(fmt.Sprintf("opened %s in %s", p.Name, choice.Profile.LabelOrID())))
	return nil
}
