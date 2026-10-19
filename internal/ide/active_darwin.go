//go:build darwin

package ide

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// AlreadyThere reports whether the IDE is frontmost and already showing projectRoot.
func AlreadyThere(choice Choice, projectRoot string) (bool, error) {
	abs, err := absClean(projectRoot)
	if err != nil {
		return false, err
	}
	if !ideIsFrontmost(choice.Profile) {
		return false, nil
	}
	open, err := folderOpenInIDE(choice.Profile, abs)
	if err != nil {
		return false, err
	}
	return open, nil
}

func ideIsFrontmost(p Profile) bool {
	name := strings.TrimSpace(p.ProcessName)
	if name == "" {
		return false
	}
	script := `tell application "System Events" to get name of first application process whose frontmost is true`
	out, err := exec.Command("osascript", "-e", script).Output()
	if err != nil {
		return false
	}
	front := strings.TrimSpace(string(out))
	return strings.EqualFold(front, name)
}

type backupWorkspaces struct {
	Workspaces []struct {
		Workspace struct {
			ConfigPath string `json:"configPath"`
		} `json:"workspace"`
	} `json:"workspaces"`
	Folders []struct {
		FolderURI string `json:"folderUri"`
	} `json:"folders"`
}

func folderOpenInIDE(p Profile, absProject string) (bool, error) {
	if p.DataDirName == "" {
		return false, nil
	}
	open, err := folderOpenFromBackupWorkspaces(p, absProject)
	if err != nil || open {
		return open, err
	}
	return folderOpenFromWorkspaceStorage(p, absProject)
}

func folderOpenFromBackupWorkspaces(p Profile, absProject string) (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	storage := filepath.Join(home, "Library", "Application Support", p.DataDirName, "User", "globalStorage", "storage.json")
	data, err := os.ReadFile(storage)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	var bw backupWorkspaces
	if err := json.Unmarshal(data, &bw); err != nil {
		return false, err
	}
	for _, f := range bw.Folders {
		if path, ok := folderURIToPath(f.FolderURI); ok && pathsEqual(path, absProject) {
			return true, nil
		}
	}
	for _, w := range bw.Workspaces {
		if path, ok := folderURIToPath(w.Workspace.ConfigPath); ok && pathsEqual(path, absProject) {
			return true, nil
		}
	}
	return false, nil
}

func folderOpenFromWorkspaceStorage(p Profile, absProject string) (bool, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}
	root := filepath.Join(home, "Library", "Application Support", p.DataDirName, "User", "workspaceStorage")
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		wsFile := filepath.Join(root, ent.Name(), "workspace.json")
		data, err := os.ReadFile(wsFile)
		if err != nil {
			continue
		}
		var meta struct {
			Folder string `json:"folder"`
		}
		if err := json.Unmarshal(data, &meta); err != nil {
			continue
		}
		path, ok := folderURIToPath(meta.Folder)
		if !ok || !pathsEqual(path, absProject) {
			continue
		}
		if workspaceSessionActive(p, ent.Name()) {
			return true, nil
		}
	}
	return false, nil
}

func workspaceSessionActive(p Profile, workspaceID string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	logsRoot := filepath.Join(home, "Library", "Application Support", p.DataDirName, "logs")
	sessions, err := os.ReadDir(logsRoot)
	if err != nil {
		return false
	}
	marker := "workspaceId-" + workspaceID
	cutoff := time.Now().Add(-45 * time.Minute)
	for _, sess := range sessions {
		if !sess.IsDir() {
			continue
		}
		if !dirContainsRecentMarker(filepath.Join(logsRoot, sess.Name()), marker, cutoff) {
			continue
		}
		return true
	}
	return false
}

func dirContainsRecentMarker(dir, marker string, cutoff time.Time) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, ent := range entries {
		path := filepath.Join(dir, ent.Name())
		if ent.IsDir() {
			if dirContainsRecentMarker(path, marker, cutoff) {
				return true
			}
			continue
		}
		if !strings.Contains(ent.Name(), marker) {
			continue
		}
		info, err := ent.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(cutoff) {
			return true
		}
	}
	return false
}
