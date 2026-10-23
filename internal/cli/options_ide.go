package cli

import (
	"fmt"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ide"
	"github.com/Cimmerial/tps/internal/ui"
)

func optionsPreferredIDERowDesc(prefs *config.Preferences) string {
	current := ide.NormalizeID(prefs.PreferredIDE)
	label := "auto (detect installed IDE)"
	if current != ide.IDAuto {
		if p, ok := ide.ProfileByID(current); ok {
			label = p.LabelOrID()
		} else {
			label = string(current)
		}
	}
	return optionsTopMenuItemDesc("default IDE", label)
}

func optionsPreferredIDEMenu(prefs *config.Preferences, projectID string, statusTick int) (bool, error) {
	ui.ClearScreen()
	current := ide.NormalizeID(prefs.PreferredIDE)
	choices := []struct {
		id    ide.ID
		label string
	}{
		{id: ide.IDAuto, label: "auto (detect installed IDE)"},
	}
	for _, p := range ide.All() {
		choices = append(choices, struct {
			id    ide.ID
			label string
		}{id: p.ID, label: p.LabelOrID()})
	}

	fmt.Println()
	fmt.Printf("%s\n", ui.HeaderChoice("options"))
	fmt.Println()
	fmt.Printf("%s\n", ui.SubHeader("default IDE"))
	fmt.Printf("    %s %s\n", ui.Dim("current:"), ui.Bold(optionsPreferredIDELabel(current)))
	fmt.Println()
	fmt.Printf("    %s\n", ui.Dim("used by: open in desired ide"))
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
	if id == ide.IDAuto {
		prefs.PreferredIDE = ""
	} else {
		prefs.PreferredIDE = string(id)
	}
	if err := config.SavePreferences(prefs); err != nil {
		return false, err
	}
	fmt.Println(ui.Dim(fmt.Sprintf("default IDE set to %s", optionsPreferredIDELabel(id))))
	if err := refreshUISession(projectID, statusTick); err != nil {
		return true, err
	}
	return true, nil
}

func optionsPreferredIDELabel(id ide.ID) string {
	if id == ide.IDAuto {
		return "auto"
	}
	if p, ok := ide.ProfileByID(id); ok {
		return p.LabelOrID()
	}
	return string(id)
}
