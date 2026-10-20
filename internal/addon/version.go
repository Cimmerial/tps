package addon

import (
	"os/exec"
	"strings"
)

// IsAddonProject reports whether projectRoot contains a tps-addon.yaml manifest.
func IsAddonProject(projectRoot string) bool {
	_, err := LoadManifest(projectRoot)
	return err == nil
}

// ProjectVersion returns the addon version from the manifest, or git describe in projectRoot.
func ProjectVersion(projectRoot string) string {
	if m, err := LoadManifest(projectRoot); err == nil {
		if v := strings.TrimSpace(m.Version); v != "" {
			return v
		}
	}
	if d := gitDescribe(projectRoot); d != "" {
		return d
	}
	return ""
}

func gitDescribe(dir string) string {
	out, err := exec.Command("git", "-C", dir, "describe", "--tags", "--always", "--dirty").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
