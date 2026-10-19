//go:build darwin

package pick

import (
	"fmt"
	"os/exec"
	"strings"
)

func ChooseFile(prompt string) (string, error) {
	prompt = escapeAppleScript(prompt)
	script := fmt.Sprintf(`POSIX path of (choose file with prompt "%s")`, prompt)
	return runChoose(script)
}

func ChooseFiles(prompt string) ([]string, error) {
	prompt = escapeAppleScript(prompt)
	script := fmt.Sprintf(`set theFiles to choose file with prompt "%s" with multiple selections allowed true
set out to ""
repeat with f in theFiles
  set out to out & (POSIX path of f) & linefeed
end repeat
return out`, prompt)
	out, err := runChoose(script)
	if err != nil {
		return nil, err
	}
	out = strings.TrimSpace(out)
	if out == "" {
		return nil, ErrCancelled
	}
	var paths []string
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			paths = append(paths, line)
		}
	}
	return paths, nil
}

func ChooseFolder(prompt string) (string, error) {
	prompt = escapeAppleScript(prompt)
	script := fmt.Sprintf(`POSIX path of (choose folder with prompt "%s")`, prompt)
	return runChoose(script)
}

func runChoose(script string) (string, error) {
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return "", ErrCancelled
	}
	path := strings.TrimSpace(string(out))
	if path == "" {
		return "", ErrCancelled
	}
	return path, nil
}

func escapeAppleScript(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
