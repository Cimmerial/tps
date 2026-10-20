package cli

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/browser"
	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/repo"
	"github.com/Cimmerial/tps/internal/ui"
)

func runOpenDocs() error {
	prefs, err := config.LoadPreferences()
	if err != nil {
		return err
	}
	if err := browser.Open(repo.ReadmeURL, prefs.PreferredBrowser); err != nil {
		return err
	}
	fmt.Println(ui.Dim(fmt.Sprintf("opened %s", repo.ReadmeURL)))
	return nil
}
