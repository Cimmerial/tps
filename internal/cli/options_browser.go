package cli

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/browser"
	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func optionsPreferredBrowserRowDesc(prefs *config.Preferences) string {
	current := browser.NormalizeID(prefs.PreferredBrowser)
	label := "default (system browser)"
	if current != browser.IDDefault {
		if p, ok := browser.ProfileByID(current); ok {
			label = p.LabelOrID()
		} else {
			label = string(current)
		}
	}
	return ui.Dim("preferred browser") + " — " + optionsRowDescCurrent(label)
}

func optionsPreferredBrowserMenu(prefs *config.Preferences, projectID string, statusTick int) (bool, error) {
	ui.ClearScreen()
	current := browser.NormalizeID(prefs.PreferredBrowser)
	choices := []struct {
		id    browser.ID
		label string
	}{
		{id: browser.IDDefault, label: "default (system browser)"},
	}
	for _, p := range browser.All() {
		choices = append(choices, struct {
			id    browser.ID
			label string
		}{id: p.ID, label: p.LabelOrID()})
	}

	fmt.Println()
	fmt.Printf("%s\n", ui.Header("options"))
	fmt.Println()
	fmt.Printf("%s\n", ui.SubHeader("preferred browser"))
	fmt.Printf("    %s %s\n", ui.Dim("current:"), ui.Bold(optionsPreferredBrowserLabel(current)))
	fmt.Println()
	fmt.Printf("    %s\n", ui.Dim("used by: doc"))
	fmt.Println()
	rows := make([]numberedMenuRow, len(choices))
	for i, c := range choices {
		desc := ui.Dim(c.label)
		if c.id == current {
			desc = optionsRowDescCurrent(c.label)
		}
		rows[i] = numberedMenuRow{Num: fmt.Sprintf("%d", i+1), Desc: desc}
	}
	printNumberedMenuRows(rows)
	fmt.Println()
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "0", Abbrev: "q", Desc: ui.Dim("no change / back")},
	})
	fmt.Println()
	choice, back, err := readOptionsMenuChoice()
	if err != nil || back {
		return false, err
	}
	if choice < 1 || choice > len(choices) {
		fmt.Println(ui.Dim("invalid choice"))
		return false, nil
	}
	id := choices[choice-1].id
	if id == current {
		fmt.Println(ui.Dim("already set"))
		return false, nil
	}
	if id == browser.IDDefault {
		prefs.PreferredBrowser = ""
	} else {
		prefs.PreferredBrowser = string(id)
	}
	if err := config.SavePreferences(prefs); err != nil {
		return false, err
	}
	fmt.Println(ui.Dim(fmt.Sprintf("preferred browser set to %s", optionsPreferredBrowserLabel(id))))
	if err := refreshUISession(projectID, statusTick); err != nil {
		return true, err
	}
	return true, nil
}

func optionsPreferredBrowserLabel(id browser.ID) string {
	if id == browser.IDDefault {
		return "default"
	}
	if p, ok := browser.ProfileByID(id); ok {
		return p.LabelOrID()
	}
	return string(id)
}
