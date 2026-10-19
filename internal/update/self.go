package update

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func InstallFromProjectRoot(root string) error {
	mainDir := filepath.Join(root, "cmd", "tps")
	if _, err := os.Stat(mainDir); err != nil {
		return fmt.Errorf("update needs a Go CLI at cmd/tps (run from the tps dev repo)")
	}
	cmd := exec.Command("go", "install", "./cmd/tps")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func InstallAndRestart(root string) error {
	if err := InstallFromProjectRoot(root); err != nil {
		return err
	}
	if err := RestartInteractive(); err != nil {
		return fmt.Errorf("installed but restart failed: %w (run tps again)", err)
	}
	return nil
}
