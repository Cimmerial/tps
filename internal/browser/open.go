package browser

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/Cimmerial/tps/internal/openurl"
)

// Open launches url in the browser chosen by preferred (empty = system default).
func Open(url string, preferred string) error {
	id := NormalizeID(preferred)
	if id == IDDefault {
		return openurl.SystemOpen(url)
	}
	p, ok := ProfileByID(id)
	if !ok {
		return fmt.Errorf("unknown browser %q", preferred)
	}
	return p.openURL(url)
}

func (p Profile) openURL(url string) error {
	switch runtime.GOOS {
	case "darwin":
		if name, ok := p.ResolveMacApp(); ok {
			return start(exec.Command("open", "-a", name, url))
		}
	case "windows", "linux":
		if bin, ok := p.ResolveBin(); ok {
			return start(exec.Command(bin, url))
		}
	}
	return fmt.Errorf("%s is not installed", p.LabelOrID())
}

func start(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
