package openurl

import (
	"fmt"
	"os/exec"
	"runtime"
)

// SystemOpen launches the OS default handler for url (typically the default browser).
func SystemOpen(url string) error {
	switch runtime.GOOS {
	case "darwin":
		return start(exec.Command("open", url))
	case "windows":
		return start(exec.Command("cmd", "/c", "start", "", url))
	default:
		bin, err := exec.LookPath("xdg-open")
		if err != nil {
			return fmt.Errorf("no browser launcher (install xdg-open)")
		}
		return start(exec.Command(bin, url))
	}
}

func start(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
