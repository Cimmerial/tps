package install

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// DefaultBinDir is where install scripts place the tps binary.
func DefaultBinDir() (string, error) {
	if runtime.GOOS == "windows" {
		local := os.Getenv("LOCALAPPDATA")
		if local == "" {
			return "", fmt.Errorf("LOCALAPPDATA not set")
		}
		return filepath.Join(local, "Programs", "tps"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "bin"), nil
}

// EnsureOnPath adds common install locations to the user shell profile when missing.
func EnsureOnPath() error {
	if runtime.GOOS == "windows" {
		return ensureWindowsPath()
	}
	return ensureUnixPath()
}

func ensureUnixPath() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	localBin := filepath.Join(home, ".local", "bin")
	goBin := ""
	if out, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		goBin = filepath.Join(strings.TrimSpace(string(out)), "bin")
	}
	dirs := []string{localBin}
	if goBin != "" {
		dirs = append(dirs, goBin)
	}
	for _, rc := range []string{filepath.Join(home, ".zshrc"), filepath.Join(home, ".bashrc")} {
		if err := appendPathToFile(rc, dirs); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func appendPathToFile(rc string, dirs []string) error {
	var content []byte
	if data, err := os.ReadFile(rc); err != nil {
		if os.IsNotExist(err) {
			if !strings.HasSuffix(rc, ".zshrc") {
				return nil
			}
			content = nil
		} else {
			return err
		}
	} else {
		content = data
	}
	text := string(content)
	changed := false
	for _, d := range dirs {
		if strings.Contains(text, d) {
			continue
		}
		line := fmt.Sprintf("export PATH=\"$PATH:%s\"\n", d)
		text += line
		changed = true
	}
	if !changed {
		return nil
	}
	return os.WriteFile(rc, []byte(text), 0o644)
}

func ensureWindowsPath() error {
	dir, err := DefaultBinDir()
	if err != nil {
		return err
	}
	goBin := ""
	if out, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		goBin = filepath.Join(strings.TrimSpace(string(out)), "bin")
	}
	dirs := []string{dir}
	if goBin != "" {
		dirs = append(dirs, goBin)
	}
	userPath := os.Getenv("Path")
	if userPath == "" {
		out, _ := exec.Command("powershell", "-NoProfile", "-Command", "[Environment]::GetEnvironmentVariable('Path','User')").Output()
		userPath = strings.TrimSpace(string(out))
	}
	for _, d := range dirs {
		if strings.Contains(strings.ToLower(userPath), strings.ToLower(d)) {
			continue
		}
		script := fmt.Sprintf("[Environment]::SetEnvironmentVariable('Path', $env:Path + ';%s', 'User')", d)
		if err := exec.Command("powershell", "-NoProfile", "-Command", script).Run(); err != nil {
			return err
		}
	}
	return nil
}
