package ui

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func CopyToClipboard(text string) error {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	case "linux":
		if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd := exec.Command("wl-copy")
			cmd.Stdin = strings.NewReader(text)
			return cmd.Run()
		}
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd := exec.Command("xclip", "-selection", "clipboard")
			cmd.Stdin = strings.NewReader(text)
			return cmd.Run()
		}
		return fmt.Errorf("install wl-copy or xclip for clipboard support")
	case "windows":
		cmd := exec.Command("powershell", "-NoProfile", "-Command", "Set-Clipboard -Value $input")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err != nil {
			clip := exec.Command("clip")
			clip.Stdin = strings.NewReader(text)
			return clip.Run()
		}
		return nil
	default:
		return fmt.Errorf("clipboard copy not supported on %s", runtime.GOOS)
	}
}
