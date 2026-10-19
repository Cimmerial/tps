package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func OpenConfigInEditor() error {
	p, err := Path()
	if err != nil {
		return err
	}
	if _, err := os.Stat(p); err != nil {
		if err := Save(&Config{}); err != nil {
			return err
		}
	}
	return openInEditor(p)
}

func findGUIEditor() (bin string, args []string) {
	if bin, args, ok := editorFromEnv(); ok {
		return bin, args
	}

	candidates := []struct {
		path string
		args []string
	}{
		{"code", []string{"-w"}},
		{"/usr/local/bin/code", []string{"-w"}},
		{filepath.Join(os.Getenv("HOME"), "bin", "code"), []string{"-w"}},
		{"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code", []string{"-w"}},
		{"cursor", []string{"-w"}},
		{filepath.Join(os.Getenv("HOME"), "bin", "cursor"), []string{"-w"}},
		{"/Applications/Cursor.app/Contents/Resources/app/bin/cursor", []string{"-w"}},
	}

	for _, c := range candidates {
		p := c.path
		if !strings.Contains(p, "/") {
			if looked, err := exec.LookPath(p); err == nil {
				p = looked
			} else {
				continue
			}
		} else if _, err := os.Stat(p); err != nil {
			continue
		}
		return p, c.args
	}
	return "", nil
}

func editorFromEnv() (bin string, args []string, ok bool) {
	raw := strings.TrimSpace(os.Getenv("TPS_EDITOR"))
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv("EDITOR"))
	}
	if raw == "" {
		return "", nil, false
	}
	parts := strings.Fields(raw)
	return parts[0], parts[1:], true
}

func openInEditor(path string) error {
	if bin, args, ok := editorFromEnv(); ok {
		return runEditor(bin, append(args, path))
	}
	if bin, args := findGUIEditor(); bin != "" {
		return runEditor(bin, append(args, path))
	}
	if runtime.GOOS == "darwin" {
		return exec.Command("open", "-t", path).Run()
	}
	return runEditor("nano", []string{path})
}

func runEditor(bin string, args []string) error {
	cmd := exec.Command(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor: %w", err)
	}
	return nil
}
