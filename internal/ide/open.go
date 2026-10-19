package ide

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Choice is a resolved editor to use for OpenProject.
type Choice struct {
	Profile Profile
	Bin     string   // CLI launcher
	MacApp  string   // macOS app name for `open -a` when Bin is empty
}

// Resolve picks the IDE from preferences, then env, then first installed GUI editor.
func Resolve(preferred string) (Choice, error) {
	id := NormalizeID(preferred)
	if id != IDAuto {
		p, ok := ProfileByID(id)
		if !ok {
			return Choice{}, fmt.Errorf("unknown IDE %q", preferred)
		}
		bin, macApp, ok := p.ResolveLauncher()
		if !ok {
			return Choice{}, fmt.Errorf("%s is not installed (no launcher on PATH or in Applications)", p.LabelOrID())
		}
		return Choice{Profile: p, Bin: bin, MacApp: macApp}, nil
	}
	if ch, ok := choiceFromEnv(); ok {
		return ch, nil
	}
	if p, ok := DetectInstalled(); ok {
		bin, macApp, _ := p.ResolveLauncher()
		return Choice{Profile: p, Bin: bin, MacApp: macApp}, nil
	}
	return Choice{}, fmt.Errorf("no supported IDE found (install Cursor, VS Code, … or set options → default IDE)")
}

func choiceFromEnv() (Choice, bool) {
	raw := strings.TrimSpace(os.Getenv("TPS_EDITOR"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("EDITOR"))
	}
	if raw == "" {
		return Choice{}, false
	}
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return Choice{}, false
	}
	bin := parts[0]
	if looked, err := exec.LookPath(bin); err == nil {
		bin = looked
	}
	if p, ok := profileFromEditorBin(bin); ok {
		return Choice{Profile: p, Bin: bin, MacApp: ""}, true
	}
	return Choice{}, false
}

// OpenProject launches the IDE on projectRoot (non-blocking).
func OpenProject(choice Choice, projectRoot string) error {
	abs, err := absClean(projectRoot)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); err != nil {
		return fmt.Errorf("project path: %w", err)
	}
	if choice.MacApp != "" {
		cmd := exec.Command("open", "-a", choice.MacApp, abs)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("open -a %q: %w", choice.MacApp, err)
		}
		return nil
	}
	cmd := exec.Command(choice.Bin, abs)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("%s: %w", filepath.Base(choice.Bin), err)
	}
	return nil
}
