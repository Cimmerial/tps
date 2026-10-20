package version

import (
	"runtime/debug"
	"strings"
)

// moduleVersionLabel is overridden in tests.
var moduleVersionLabel = defaultModuleVersionLabel

func defaultModuleVersionLabel() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	dirty := vcsModified(info)
	v := strings.TrimSpace(info.Main.Version)
	if v != "" && v != "(devel)" {
		return formatModuleVersion(v, dirty)
	}
	if v == "(devel)" && DevVersion != "" {
		return formatDevVersion(dirty)
	}
	return ""
}

func vcsModified(info *debug.BuildInfo) bool {
	for _, s := range info.Settings {
		if s.Key == "vcs.modified" && s.Value == "true" {
			return true
		}
	}
	return false
}

// formatModuleVersion turns go module versions (including pseudo-versions) into vX.Y.Z[-dirty].
func formatModuleVersion(mod string, dirty bool) string {
	mod = strings.TrimSpace(mod)
	if mod == "" {
		return ""
	}
	if i := strings.Index(mod, "+"); i > 0 {
		mod = mod[:i]
		dirty = true
	}
	ver := mod
	if i := strings.Index(ver, "-0."); i > 0 {
		ver = ver[:i]
	}
	if !strings.HasPrefix(ver, "v") {
		ver = "v" + ver
	}
	if dirty {
		return ver + "-dirty"
	}
	return ver
}

func formatDevVersion(dirty bool) string {
	ver := "v" + strings.TrimPrefix(strings.TrimSpace(DevVersion), "v")
	if dirty {
		return ver + "-dirty"
	}
	return ver
}
