package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ide"
	"github.com/Cimmerial/tps/internal/ui"
	"github.com/Cimmerial/tps/internal/ui/logos"
)

type prefSnapshot struct {
	logoID         string
	titleCardColor string
	accent         string
	colorLabel     string
	smsTo          string
	smsShortcut    string
	preferredIDE   string
}

func snapshotPreferences(p *config.Preferences) prefSnapshot {
	return prefSnapshot{
		logoID:         p.LogoID(),
		titleCardColor: p.TitleCardColorMode(),
		accent:         p.AccentMode(),
		colorLabel:     colorPrefLabel(p.ColorEnabled),
		smsTo:          strings.TrimSpace(p.SMSTo),
		smsShortcut:    strings.TrimSpace(p.SMSShortcut),
		preferredIDE:   string(ide.NormalizeID(p.PreferredIDE)),
	}
}

func formatPrefChange(label, oldVal, newVal string) string {
	return label + ": " + ui.Red(oldVal) + ui.Dim(" -> ") + ui.Green(newVal)
}

func optionsSessionChanges(before, after prefSnapshot) []string {
	var lines []string
	if before.logoID != after.logoID {
		lines = append(lines, formatPrefChange("title card", before.logoID, after.logoID))
	}
	if before.titleCardColor != after.titleCardColor {
		lines = append(lines, formatPrefChange("title card color", before.titleCardColor, after.titleCardColor))
	}
	if before.accent != after.accent {
		lines = append(lines, formatPrefChange("accent", before.accent, after.accent))
	}
	if before.colorLabel != after.colorLabel {
		lines = append(lines, formatPrefChange("color", before.colorLabel, after.colorLabel))
	}
	if before.smsTo != after.smsTo {
		lines = append(lines, formatPrefChange("sms_to", smsToChangeDisplay(before.smsTo), smsToChangeDisplay(after.smsTo)))
	}
	if before.smsShortcut != after.smsShortcut {
		lines = append(lines, formatPrefChange("sms_shortcut", prefDisplay(before.smsShortcut), prefDisplay(after.smsShortcut)))
	}
	if before.preferredIDE != after.preferredIDE {
		lines = append(lines, formatPrefChange("default IDE", optionsPreferredIDELabel(ide.ID(before.preferredIDE)), optionsPreferredIDELabel(ide.ID(after.preferredIDE))))
	}
	return lines
}

func prefDisplay(v string) string {
	if v == "" {
		return "(empty)"
	}
	return v
}

func finishOptionsSession(projectID string, statusTick int, start prefSnapshot) error {
	endPrefs, err := config.LoadPreferences()
	if err != nil {
		return err
	}
	if err := refreshInteractiveSession(projectID, statusTick, "", false); err != nil {
		return err
	}
	changes := optionsSessionChanges(start, snapshotPreferences(endPrefs))
	if len(changes) == 0 {
		return nil
	}
	fmt.Println()
	fmt.Println(ui.Header("OPTIONS CHANGES:"))
	for _, line := range changes {
		fmt.Printf("  - %s\n", line)
	}
	return nil
}

func runOptionsMenu(projectID string, statusTick int) error {
	startPrefs, err := config.LoadPreferences()
	if err != nil {
		return err
	}
	startSnap := snapshotPreferences(startPrefs)

	for {
		ui.ClearScreen()
		prefs, err := config.LoadPreferences()
		if err != nil {
			return err
		}
		fmt.Println()
		fmt.Printf("%s\n", ui.Header("options"))
		fmt.Println()
		fmt.Printf("    %s\n", ui.Dim("numbered choices — 0 or q = done"))
		fmt.Println()
		printNumberedMenuRows([]numberedMenuRow{
			{Num: "1", Desc: ui.Dim("title card logo")},
			{Num: "2", Desc: ui.Dim("title card color")},
			{Num: "3", Desc: ui.Dim("terminal colors")},
			{Num: "4", Desc: optionsCaptureSMSRowDesc(prefs)},
			{Num: "5", Desc: optionsPreferredIDERowDesc(prefs)},
		})
		fmt.Println()
		printNumberedMenuRows([]numberedMenuRow{
			{Num: "0", Abbrev: "q", Desc: ui.Dim("done")},
		})
		fmt.Println()
		choice, back, err := readOptionsMenuChoice()
		if err != nil {
			return err
		}
		if back {
			return finishOptionsSession(projectID, statusTick, startSnap)
		}
		switch choice {
		case 1:
			changed, err := optionsLogoMenu(prefs, projectID, statusTick)
			if err != nil {
				return err
			}
			if changed {
				continue
			}
		case 2:
			changed, err := optionsTitleCardColorMenu(prefs, projectID, statusTick)
			if err != nil {
				return err
			}
			if changed {
				continue
			}
		case 3:
			changed, err := optionsColorsMenu(prefs, projectID, statusTick)
			if err != nil {
				return err
			}
			if changed {
				continue
			}
		case 4:
			changed, err := optionsCaptureSMSMenu(prefs, projectID, statusTick)
			if err != nil {
				return err
			}
			if changed {
				continue
			}
		case 5:
			changed, err := optionsPreferredIDEMenu(prefs, projectID, statusTick)
			if err != nil {
				return err
			}
			if changed {
				continue
			}
		default:
			fmt.Println(ui.Dim("invalid choice"))
		}
	}
}

func optionsLogoMenu(prefs *config.Preferences, projectID string, statusTick int) (bool, error) {
	list := logos.Presets()
	for {
		ui.ClearScreen()
		current := prefs.LogoID()
		fmt.Println()
		fmt.Printf("%s\n", ui.Header("options"))
		fmt.Println()
		fmt.Printf("%s\n", ui.SubHeader("title card logo"))
		fmt.Printf("    %s %s\n", ui.Dim("current:"), ui.Bold(current)+" "+ui.Dim("("+logos.Label(current)+")"))
		fmt.Println()
		fmt.Printf("    %s\n", ui.Dim("pick a logo to preview — 0 or q = back"))
		fmt.Println()
		rows := make([]numberedMenuRow, len(list))
		for i, p := range list {
			desc := p.Label
			if p.ID == current {
				desc = optionsRowDescCurrent(p.Label)
			} else {
				desc = ui.Dim(p.Label)
			}
			rows[i] = numberedMenuRow{
				Num:  fmt.Sprintf("%d", i+1),
				Desc: desc,
			}
		}
		printNumberedMenuRows(rows)
		fmt.Println()
		printNumberedMenuRows([]numberedMenuRow{
			{Num: "0", Abbrev: "q", Desc: ui.Dim("back")},
		})
		fmt.Println()

		choice, back, err := readOptionsMenuChoice()
		if err != nil || back {
			return false, err
		}
		if choice < 1 || choice > len(list) {
			fmt.Println(ui.Dim("invalid choice"))
			continue
		}
		id := list[choice-1].ID
		applied, err := optionsLogoPreviewApply(prefs, projectID, statusTick, id, current)
		if err != nil {
			return applied, err
		}
		if applied {
			return true, nil
		}
	}
}

func optionsLogoPreviewApply(prefs *config.Preferences, projectID string, statusTick int, id, current string) (bool, error) {
	ui.ClearScreen()
	fmt.Println()
	fmt.Printf("%s\n", ui.Header("options"))
	fmt.Println()
	fmt.Printf("%s\n", ui.SubHeader("title card logo"))
	fmt.Printf("    %s %s\n", ui.Dim("preview:"), ui.Bold(logos.Label(id)))
	fmt.Println()
	ui.PrintLogoPreview(id)
	fmt.Println()
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "1", Desc: ui.Dim("use this logo")},
		{Num: "0", Abbrev: "q", Desc: ui.Dim("back")},
	})
	fmt.Println()
	choice, back, err := readOptionsMenuChoice()
	if err != nil || back || choice != 1 {
		return false, err
	}
	if id == current {
		fmt.Println(ui.Dim("already set"))
		return false, nil
	}
	if id == "default" {
		prefs.Logo = ""
		prefs.TitleCard = ""
	} else {
		prefs.Logo = id
		prefs.TitleCard = ""
	}
	if err := config.SavePreferences(prefs); err != nil {
		return false, err
	}
	if err := applyUIPreferences(); err != nil {
		return false, err
	}
	fmt.Println(ui.Dim(fmt.Sprintf("title card set to %s", id)))
	if err := refreshUISession(projectID, statusTick); err != nil {
		return true, err
	}
	return true, nil
}

func colorPresetChoices() []config.ColorPreset {
	return config.ColorPresets()
}

func optionsColorMenuDesc(p config.ColorPreset, selected bool) string {
	label := ui.PaintColorID(p.ID, p.Label)
	if selected {
		return optionsRowDescCurrent(label)
	}
	return ui.Dim(label)
}

func colorPresetHelpLine(id string) string {
	switch id {
	case "orange", "purple", "pink", "teal":
		return "256-color terminal foreground"
	case "lime":
		return "alias of bright green"
	default:
		return ""
	}
}

func optionsTitleCardColorMenu(prefs *config.Preferences, projectID string, statusTick int) (bool, error) {
	ui.ClearScreen()
	current := prefs.TitleCardColorMode()
	choices := colorPresetChoices()
	fmt.Println()
	fmt.Printf("%s\n", ui.Header("options"))
	fmt.Println()
	fmt.Printf("%s\n", ui.SubHeader("title card color"))
	fmt.Printf("    %s %s\n", ui.Dim("current:"), ui.Bold(current)+" "+ui.PaintColorID(current, current))
	fmt.Println()
	fmt.Printf("    %s\n", ui.Dim("logo / title card art only (not menus)"))
	fmt.Println()
	rows := make([]numberedMenuRow, len(choices))
	for i, c := range choices {
		rows[i] = numberedMenuRow{
			Num:  fmt.Sprintf("%d", i+1),
			Desc: optionsColorMenuDesc(c, c.ID == current),
		}
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
	id := choices[choice-1].ID
	if id == current {
		fmt.Println(ui.Dim("already set"))
		return false, nil
	}
	prefs.TitleCardColor = id
	if err := config.SavePreferences(prefs); err != nil {
		return false, err
	}
	if err := applyUIPreferences(); err != nil {
		return false, err
	}
	fmt.Println(ui.Dim(fmt.Sprintf("title card color set to %s", id)))
	if err := refreshUISession(projectID, statusTick); err != nil {
		return true, err
	}
	return true, nil
}

func optionsColorsMenu(prefs *config.Preferences, projectID string, statusTick int) (bool, error) {
	ui.ClearScreen()
	fmt.Println()
	fmt.Printf("%s\n", ui.Header("options"))
	fmt.Println()
	fmt.Printf("%s\n", ui.SubHeader("terminal colors"))
	fmt.Printf("    %s %s\n", ui.Dim("accent:"), ui.Bold(prefs.AccentMode())+" "+ui.PaintColorID(prefs.AccentMode(), prefs.AccentMode()))
	fmt.Printf("    %s %s\n", ui.Dim("color:"), ui.Bold(colorPrefLabel(prefs.ColorEnabled)))
	fmt.Println()
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "1", Desc: ui.Dim("accent color")},
		{Num: "2", Desc: ui.Dim("color on/off (auto when unset)")},
	})
	fmt.Println()
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "0", Abbrev: "q", Desc: ui.Dim("back")},
	})
	fmt.Println()
	choice, back, err := readOptionsMenuChoice()
	if err != nil || back {
		return false, err
	}
	switch choice {
	case 1:
		return optionsAccentMenu(prefs, projectID, statusTick)
	case 2:
		return optionsColorToggleMenu(prefs, projectID, statusTick)
	default:
		fmt.Println(ui.Dim("invalid choice"))
		return false, nil
	}
}

func optionsAccentMenu(prefs *config.Preferences, projectID string, statusTick int) (bool, error) {
	ui.ClearScreen()
	current := prefs.AccentMode()
	choices := colorPresetChoices()
	fmt.Println()
	fmt.Printf("%s\n", ui.Header("options"))
	fmt.Println()
	fmt.Printf("%s\n", ui.SubHeader("accent"))
	fmt.Printf("    %s %s\n", ui.Dim("current:"), ui.Bold(current)+" "+ui.PaintColorID(current, current))
	fmt.Println()
	fmt.Printf("    %s\n", ui.Dim("headers, prompts, and menu shortcut keys"))
	fmt.Println()
	rows := make([]numberedMenuRow, len(choices))
	for i, c := range choices {
		rows[i] = numberedMenuRow{
			Num:      fmt.Sprintf("%d", i+1),
			Desc:     optionsColorMenuDesc(c, c.ID == current),
			HelpLine: colorPresetHelpLine(c.ID),
		}
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
	id := choices[choice-1].ID
	if id == current {
		fmt.Println(ui.Dim("already set"))
		return false, nil
	}
	prefs.Accent = id
	if err := config.SavePreferences(prefs); err != nil {
		return false, err
	}
	if err := applyUIPreferences(); err != nil {
		return false, err
	}
	fmt.Println(ui.Dim(fmt.Sprintf("accent set to %s", id)))
	if err := refreshUISession(projectID, statusTick); err != nil {
		return true, err
	}
	return true, nil
}

func optionsColorToggleMenu(prefs *config.Preferences, projectID string, statusTick int) (bool, error) {
	ui.ClearScreen()
	current := prefs.ColorEnabled
	fmt.Println()
	fmt.Printf("%s\n", ui.Header("options"))
	fmt.Println()
	fmt.Printf("%s\n", ui.SubHeader("color mode"))
	fmt.Printf("    %s %s\n", ui.Dim("current:"), ui.Bold(colorPrefLabel(current)))
	fmt.Println()
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "1", Desc: optionsRowDescPlain("Auto (TTY)", sameColorPref(current, nil)), HelpLine: "color when stdout is a terminal, plain when piped"},
		{Num: "2", Desc: optionsRowDescPlain("Force on", sameColorPref(current, boolPtr(true))), HelpLine: "always ANSI colors"},
		{Num: "3", Desc: optionsRowDescPlain("Force off", sameColorPref(current, boolPtr(false))), HelpLine: "never color (plain text)"},
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
	var next *bool
	switch choice {
	case 1:
		next = nil
	case 2:
		on := true
		next = &on
	case 3:
		off := false
		next = &off
	default:
		fmt.Println(ui.Dim("invalid choice"))
		return false, nil
	}
	if sameColorPref(prefs.ColorEnabled, next) {
		fmt.Println(ui.Dim("already set"))
		return false, nil
	}
	prefs.ColorEnabled = next
	if err := config.SavePreferences(prefs); err != nil {
		return false, err
	}
	if err := applyUIPreferences(); err != nil {
		return false, err
	}
	fmt.Println(ui.Dim(fmt.Sprintf("color mode: %s", colorPrefLabel(next))))
	if err := refreshUISession(projectID, statusTick); err != nil {
		return true, err
	}
	return true, nil
}

func optionsRowDescCurrent(label string) string {
	return ui.Green(label + " (current)")
}

func optionsRowDescIdle(label string, selected bool) string {
	if selected {
		return optionsRowDescCurrent(label)
	}
	return ui.Dim(label)
}

func optionsRowDescPlain(label string, selected bool) string {
	if selected {
		return optionsRowDescCurrent(label)
	}
	return label
}

func boolPtr(v bool) *bool {
	return &v
}

func colorPrefLabel(v *bool) string {
	if v == nil {
		return "auto"
	}
	if *v {
		return "on"
	}
	return "off"
}

func sameColorPref(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// readOptionsMenuChoice reads a numbered menu line. back is true for 0, q, quit, or exit.
func readOptionsMenuChoice() (choice int, back bool, err error) {
	line, err := ui.ReadEditedLine(ui.Prompt())
	if err != nil {
		return 0, false, err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, true, nil
	}
	low := strings.ToLower(line)
	if low == "0" || low == "q" || low == "quit" || low == "exit" {
		return 0, true, nil
	}
	n, err := strconv.Atoi(line)
	if err != nil {
		fmt.Println(ui.Dim("enter a number"))
		return 0, true, nil
	}
	return n, false, nil
}
