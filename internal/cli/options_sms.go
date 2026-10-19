package cli

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func optionsCaptureSMSRowDesc(prefs *config.Preferences) string {
	base := ui.Dim("capture sms (phone / Shortcuts)")
	if prefs == nil {
		return base
	}
	to := strings.TrimSpace(prefs.SMSTo)
	if to == "" {
		return base + ui.Dim(" — not set")
	}
	return base + ui.Dim(" — ") + ui.Bold(maskSMSTo(to))
}

func optionsCaptureSMSMenu(_ *config.Preferences, projectID string, statusTick int) (bool, error) {
	prefsPath, err := config.PreferencesPath()
	if err != nil {
		return false, err
	}
	showSMSTo := false
	for {
		prefs, err := config.LoadPreferences()
		if err != nil {
			return false, err
		}
		ui.ClearScreen()
		fmt.Println()
		fmt.Printf("%s\n", ui.Header("options"))
		fmt.Println()
		fmt.Printf("%s\n", ui.SubHeader("capture sms"))
		fmt.Printf("    %s %s\n", ui.Dim("saved to:"), ui.Dim(prefsPath))
		fmt.Printf("    %s %s\n", ui.Dim("sms_to:"), smsToDisplay(prefs.SMSTo, showSMSTo))
		fmt.Printf("    %s %s\n", ui.Dim("sms_shortcut:"), shortcutDisplay(prefs.SMSShortcut))
		fmt.Println()
		fmt.Printf("    %s\n", ui.Dim("then enable sms on a command: edit → capture → toggle sms"))
		fmt.Println()
		printNumberedMenuRows([]numberedMenuRow{
			{Num: "1", Desc: ui.Dim("set sms_to (phone number)")},
			{Num: "2", Desc: ui.Dim("set sms_shortcut (optional Shortcuts name)")},
			{Num: "3", Desc: ui.Dim("clear sms_to and sms_shortcut")},
			{Num: "4", Desc: optionsToggleSMSToVisibilityRow(showSMSTo)},
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
			if err := optionsPromptSMSTo(prefs); err != nil {
				return false, err
			}
		case 2:
			if err := optionsPromptSMSShortcut(prefs); err != nil {
				return false, err
			}
		case 3:
			if prefs.SMSTo == "" && prefs.SMSShortcut == "" {
				fmt.Println(ui.Dim("already empty"))
				continue
			}
			prefs.SMSTo = ""
			prefs.SMSShortcut = ""
			if err := config.SavePreferences(prefs); err != nil {
				return false, err
			}
			fmt.Println(ui.Dim("cleared sms preferences"))
		case 4:
			if strings.TrimSpace(prefs.SMSTo) == "" {
				fmt.Println(ui.Dim("sms_to not set"))
				continue
			}
			showSMSTo = !showSMSTo
		case 0:
			return false, nil
		default:
			fmt.Println(ui.Dim("invalid choice"))
		}
	}
}

func smsToChangeDisplay(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "(empty)"
	}
	return maskSMSTo(v)
}

func maskSMSTo(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return strings.Repeat("*", 6) + v[len(v)-4:]
}

func optionsToggleSMSToVisibilityRow(visible bool) string {
	if visible {
		return ui.Dim("hide sms_to (phone number)")
	}
	return ui.Dim("show sms_to (phone number)")
}

func smsToDisplay(v string, reveal bool) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ui.Dim("(not set)")
	}
	if reveal {
		return ui.Bold(v)
	}
	return ui.Bold(maskSMSTo(v))
}

func shortcutDisplay(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ui.Dim("(not set — uses Messages)")
	}
	return ui.Bold(v)
}

func optionsPromptSMSTo(prefs *config.Preferences) error {
	fmt.Println()
	fmt.Printf("    %s %s\n", ui.Dim("current:"), smsToDisplay(prefs.SMSTo, false))
	line, err := ui.ReadEditedLine("sms_to: ")
	if err != nil {
		return err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		fmt.Println(ui.Dim("no change"))
		return nil
	}
	if line == prefs.SMSTo {
		fmt.Println(ui.Dim("already set"))
		return nil
	}
	prefs.SMSTo = line
	if err := config.SavePreferences(prefs); err != nil {
		return err
	}
	fmt.Println(ui.Dim("sms_to saved"))
	return nil
}

func optionsPromptSMSShortcut(prefs *config.Preferences) error {
	fmt.Println()
	fmt.Printf("    %s %s\n", ui.Dim("current:"), shortcutDisplay(prefs.SMSShortcut))
	line, err := ui.ReadEditedLine("sms_shortcut: ")
	if err != nil {
		return err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		fmt.Println(ui.Dim("no change"))
		return nil
	}
	if line == prefs.SMSShortcut {
		fmt.Println(ui.Dim("already set"))
		return nil
	}
	prefs.SMSShortcut = line
	if err := config.SavePreferences(prefs); err != nil {
		return err
	}
	fmt.Println(ui.Dim("sms_shortcut saved"))
	return nil
}
