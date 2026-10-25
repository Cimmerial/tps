package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Executable returns a git binary path (PATH or common install locations on Windows).
func Executable() (string, error) {
	if p, err := exec.LookPath("git"); err == nil {
		return p, nil
	}
	if runtime.GOOS == "windows" {
		for _, p := range windowsGitCandidates() {
			if st, err := os.Stat(p); err == nil && !st.IsDir() {
				return p, nil
			}
		}
	}
	return "", fmt.Errorf("git not found in PATH")
}

func windowsGitCandidates() []string {
	var out []string
	if pf := os.Getenv("ProgramFiles"); pf != "" {
		out = append(out, filepath.Join(pf, "Git", "cmd", "git.exe"))
	}
	if pf86 := os.Getenv("ProgramFiles(x86)"); pf86 != "" {
		out = append(out, filepath.Join(pf86, "Git", "cmd", "git.exe"))
	}
	if local := os.Getenv("LOCALAPPDATA"); local != "" {
		matches, _ := filepath.Glob(filepath.Join(local, "GitHubDesktop", "app-*", "resources", "app", "git", "cmd", "git.exe"))
		out = append(out, matches...)
	}
	return out
}
