package notify

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

// SendSMS delivers body to to (E.164 or local format). Uses Shortcuts when shortcutName is set.
func SendSMS(to, body, shortcutName string) error {
	to = strings.TrimSpace(to)
	body = strings.TrimSpace(body)
	if to == "" {
		p, err := config.PreferencesPath()
		if err != nil {
			return fmt.Errorf("sms_to is not configured (set sms_to in preferences.yaml)")
		}
		return fmt.Errorf("sms_to is not configured (set sms_to in %s or options → capture sms)", p)
	}
	if body == "" {
		return fmt.Errorf("empty message")
	}
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("sms is only supported on macOS")
	}
	shortcutName = strings.TrimSpace(shortcutName)
	if shortcutName != "" {
		if err := runShortcut(shortcutName, body); err == nil {
			return nil
		}
		// fall through to Messages
	}
	return sendMessagesAppleScript(to, body)
}

func runShortcut(name, input string) error {
	cmd := exec.Command("shortcuts", "run", name, "--input", input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// messagesSendScript is executed via osascript - with argv [to, body].
const messagesSendScript = `
on run argv
	set theBuddy to item 1 of argv
	set theMessage to item 2 of argv
	tell application "Messages"
		launch
	end tell
	tell application "Messages"
		try
			send theMessage to buddy theBuddy of service "SMS"
		on error
			try
				send theMessage to buddy theBuddy of (first account whose service type is iMessage)
			on error
				send theMessage to buddy theBuddy of (first service whose service type is iMessage)
			end try
		end try
	end tell
end run
`

func sendMessagesAppleScript(to, body string) error {
	cmd := exec.Command("osascript", "-", to, body)
	cmd.Stdin = strings.NewReader(messagesSendScript)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
