//go:build darwin

package pick

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Cimmerial/tps/internal/ui"
)

func ChooseFile(prompt string, filter FileFilter) (string, error) {
	path, err := chooseFileOSA(prompt, filter)
	if err == nil {
		if err := filter.Validate(path); err != nil {
			return "", err
		}
		return path, nil
	}
	if err == ErrCancelled {
		return "", err
	}
	printOSAFallback(err)
	path, err = readPathLine(prompt)
	if err != nil {
		return "", err
	}
	if err := filter.Validate(path); err != nil {
		return "", err
	}
	return path, nil
}

func ChooseFiles(prompt string, filter FileFilter) ([]string, error) {
	paths, err := chooseFilesOSA(prompt, filter)
	if err == nil {
		for _, path := range paths {
			if err := filter.Validate(path); err != nil {
				return nil, err
			}
		}
		return paths, nil
	}
	if err == ErrCancelled {
		return nil, err
	}
	printOSAFallback(err)
	paths, err = readPathsLine(prompt)
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		if err := filter.Validate(path); err != nil {
			return nil, err
		}
	}
	return paths, nil
}

func ChooseFolder(prompt string) (string, error) {
	path, err := chooseFolderOSA(prompt)
	if err == nil {
		return path, nil
	}
	if err == ErrCancelled {
		return "", err
	}
	printOSAFallback(err)
	return readPathLine(prompt)
}

func printOSAFallback(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", ui.Yellow(strings.TrimSpace(err.Error())))
	fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("Finder picker unavailable — type path(s) below, or run tps from Terminal.app if you need the dialog."))
	fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("(Cursor’s terminal often blocks AppleScript — System Settings → Privacy & Security → Automation.)"))
}

func chooseFilesOSA(prompt string, filter FileFilter) ([]string, error) {
	prompt = escapeAppleScript(prompt)
	typeSuffix := applescriptTypeSuffix(filter.Extensions)
	lines := []string{
		fmt.Sprintf(`set theFiles to choose file with prompt "%s" with multiple selections allowed%s`, prompt, typeSuffix),
		`set out to ""`,
		`repeat with f in theFiles`,
		`set out to out & (POSIX path of f) & linefeed`,
		`end repeat`,
		`return out`,
	}
	out, err := runOSA(lines)
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

func chooseFileOSA(prompt string, filter FileFilter) (string, error) {
	prompt = escapeAppleScript(prompt)
	typeSuffix := applescriptTypeSuffix(filter.Extensions)
	return runOSA([]string{fmt.Sprintf(`POSIX path of (choose file with prompt "%s"%s)`, prompt, typeSuffix)})
}

func chooseFolderOSA(prompt string) (string, error) {
	prompt = escapeAppleScript(prompt)
	return runOSA([]string{fmt.Sprintf(`POSIX path of (choose folder with prompt "%s")`, prompt)})
}

func runOSA(scriptLines []string) (string, error) {
	fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("Opening Finder picker… (check behind this window if you do not see it)"))
	args := make([]string, 0, len(scriptLines)*2)
	for _, line := range scriptLines {
		args = append(args, "-e", line)
	}
	cmd := exec.Command("osascript", args...)
	combined, err := cmd.CombinedOutput()
	out := strings.TrimSpace(string(combined))
	if err != nil {
		if osaUserCancelled(combined) {
			return "", ErrCancelled
		}
		msg := strings.TrimSpace(string(combined))
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("osascript: %s", msg)
	}
	if out == "" {
		return "", ErrCancelled
	}
	return out, nil
}

func osaUserCancelled(combined []byte) bool {
	s := string(combined)
	return strings.Contains(s, "User canceled") || strings.Contains(s, "User cancelled") || strings.Contains(s, "(-128)")
}

func escapeAppleScript(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

