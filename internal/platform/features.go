package platform

import "runtime"

// SupportsCaptureSMS is true where capture can send SMS (macOS Messages / Shortcuts).
func SupportsCaptureSMS() bool {
	return runtime.GOOS == "darwin"
}

// CloseAllTermsAbbrev is the interactive shortcut for closing/killing detached terminals.
func CloseAllTermsAbbrev() string {
	return "vite"
}

// CloseAllTermsHelp returns help text for the full help menu (session section).
func CloseAllTermsHelp() (command, abbrev, description string) {
	switch runtime.GOOS {
	case "darwin":
		return "close all terms", "vite", "kill tab processes and close every Terminal.app window"
	case "windows":
		return "close all terms", "vite", "stop every tps detached terminal job (all projects)"
	default:
		return "close all terms", "vite", "stop every tps detached terminal job (all projects)"
	}
}

// OptionsMenuDescription is the options (o) row in full help.
func OptionsMenuDescription() string {
	base := "logo/color, terminal colors, default IDE, preferred browser"
	if SupportsCaptureSMS() {
		return "logo/color, terminal colors, capture sms, " + base
	}
	return base
}

// TerminalsHelpLine documents t/tk/vite in full help.
func TerminalsHelpLine() string {
	abbr := CloseAllTermsAbbrev()
	switch runtime.GOOS {
	case "darwin":
		return "t · tk · tk -a · " + abbr + " · t -k 1 · tr"
	default:
		return "t · tk · tk -a · " + abbr + " · t -k 1 · tr"
	}
}
