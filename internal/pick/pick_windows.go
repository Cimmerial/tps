//go:build windows

package pick

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps/internal/ui"
)

func ChooseFile(prompt string, filter FileFilter) (string, error) {
	path, err := chooseFileWin(prompt, filter)
	if err == nil {
		if err := filter.Validate(path); err != nil {
			return "", err
		}
		return path, nil
	}
	if err == ErrCancelled {
		return "", err
	}
	printWinFallback(err)
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
	paths, err := chooseFilesWin(prompt, filter)
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
	printWinFallback(err)
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
	path, err := chooseFolderWin(prompt)
	if err == nil {
		return path, nil
	}
	if err == ErrCancelled {
		return "", err
	}
	printWinFallback(err)
	return readPathLine(prompt)
}

func ChooseSaveFile(prompt string, defaultName string) (string, error) {
	defaultName = strings.TrimSpace(defaultName)
	if defaultName == "" {
		defaultName = "export.tps.yaml"
	}
	path, err := chooseSaveFileWin(prompt, defaultName)
	if err == nil {
		return path, nil
	}
	if err == ErrCancelled {
		return "", err
	}
	printWinFallback(err)
	path, err = readPathLine(prompt + " (path for " + defaultName + ")")
	if err != nil {
		return "", err
	}
	return ResolveSavePath(path, defaultName), nil
}

func printWinFallback(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", ui.Yellow(strings.TrimSpace(err.Error())))
	fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("File picker unavailable — type path(s) below."))
}

func chooseFileWin(prompt string, filter FileFilter) (string, error) {
	lines, err := runWinPickerScript(winOpenFileScript(false), prompt, filter)
	if err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "", ErrCancelled
	}
	return lines[0], nil
}

func chooseFilesWin(prompt string, filter FileFilter) ([]string, error) {
	return runWinPickerScript(winOpenFileScript(true), prompt, filter)
}

func chooseFolderWin(prompt string) (string, error) {
	lines, err := runWinPickerScript(winFolderScript(), prompt, FileFilter{})
	if err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "", ErrCancelled
	}
	return lines[0], nil
}

func chooseSaveFileWin(prompt, defaultName string) (string, error) {
	prevDefault := os.Getenv("TPS_PICK_DEFAULT")
	_ = os.Setenv("TPS_PICK_DEFAULT", defaultName)
	defer func() { _ = os.Setenv("TPS_PICK_DEFAULT", prevDefault) }()
	lines, err := runWinPickerScript(winSaveFileScript(defaultName), prompt, FileFilter{})
	if err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "", ErrCancelled
	}
	return lines[0], nil
}

func runWinPickerScript(script, prompt string, filter FileFilter) ([]string, error) {
	fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("Opening file picker… (check behind this window if you do not see it)"))
	prevPrompt := os.Getenv("TPS_PICK_PROMPT")
	prevFilter := os.Getenv("TPS_PICK_FILTER")
	_ = os.Setenv("TPS_PICK_PROMPT", prompt)
	_ = os.Setenv("TPS_PICK_FILTER", windowsFileFilter(filter.Extensions))
	defer func() {
		_ = os.Setenv("TPS_PICK_PROMPT", prevPrompt)
		_ = os.Setenv("TPS_PICK_FILTER", prevFilter)
	}()

	cmd := exec.Command("powershell.exe", "-NoProfile", "-STA", "-Command", script)
	out, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err != nil {
		if winPickerCancelled(text, err) {
			return nil, ErrCancelled
		}
		msg := text
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("file picker: %s", msg)
	}
	if text == "" {
		return nil, ErrCancelled
	}
	var lines []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

func winPickerCancelled(out string, err error) bool {
	if strings.Contains(strings.ToLower(out), "canceled") || strings.Contains(strings.ToLower(out), "cancelled") {
		return true
	}
	if ee, ok := err.(*exec.ExitError); ok && ee.ExitCode() == 1 && strings.TrimSpace(out) == "" {
		return true
	}
	return false
}

func windowsFileFilter(extensions []string) string {
	if len(extensions) == 0 {
		return "All files (*.*)|*.*"
	}
	var globs []string
	for _, e := range extensions {
		e = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(e), "."))
		if e != "" {
			globs = append(globs, "*."+e)
		}
	}
	if len(globs) == 0 {
		return "All files (*.*)|*.*"
	}
	glob := strings.Join(globs, ";")
	label := strings.Join(globs, ", ")
	return label + "|" + glob + "|All files (*.*)|*.*"
}

func winOpenFileScript(multi bool) string {
	multiLine := "$false"
	if multi {
		multiLine = "$true"
	}
	return `
Add-Type -AssemblyName System.Windows.Forms
$d = New-Object System.Windows.Forms.OpenFileDialog
$d.Title = $env:TPS_PICK_PROMPT
$d.Filter = $env:TPS_PICK_FILTER
$d.Multiselect = ` + multiLine + `
$r = $d.ShowDialog()
if ($r -eq [System.Windows.Forms.DialogResult]::OK) {
  if ($d.Multiselect) {
    foreach ($f in $d.FileNames) { [Console]::Out.WriteLine($f) }
  } else {
    [Console]::Out.WriteLine($d.FileName)
  }
  exit 0
}
exit 1
`
}

func winFolderScript() string {
	return `
Add-Type -AssemblyName System.Windows.Forms
$d = New-Object System.Windows.Forms.FolderBrowserDialog
$d.Description = $env:TPS_PICK_PROMPT
$d.UseDescriptionForTitle = $true
$r = $d.ShowDialog()
if ($r -eq [System.Windows.Forms.DialogResult]::OK) {
  [Console]::Out.WriteLine($d.SelectedPath)
  exit 0
}
exit 1
`
}

func winSaveFileScript(defaultName string) string {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(defaultName)), ".")
	addExt := ""
	if ext != "" {
		addExt = fmt.Sprintf("$d.DefaultExt = '%s'\n", escapePowerShellSingleQuoted(ext))
	}
	return `
Add-Type -AssemblyName System.Windows.Forms
$d = New-Object System.Windows.Forms.SaveFileDialog
$d.Title = $env:TPS_PICK_PROMPT
$d.FileName = $env:TPS_PICK_DEFAULT
` + addExt + `
$r = $d.ShowDialog()
if ($r -eq [System.Windows.Forms.DialogResult]::OK) {
  [Console]::Out.WriteLine($d.FileName)
  exit 0
}
exit 1
`
}

func escapePowerShellSingleQuoted(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
