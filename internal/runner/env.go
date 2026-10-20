package runner

import (
	"os"
	"runtime"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

// augmentEnviron returns os.Environ() with common tool paths prepended to PATH (macOS Homebrew, etc.).
func augmentEnviron() []string {
	if runtime.GOOS == "windows" {
		return os.Environ()
	}
	extra := strings.Join([]string{
		"/opt/homebrew/bin",
		"/usr/local/bin",
		"/opt/local/bin",
	}, ":")
	env := os.Environ()
	for i, e := range env {
		if !strings.HasPrefix(e, "PATH=") {
			continue
		}
		cur := strings.TrimPrefix(e, "PATH=")
		if pathHasPrefix(cur, extra) {
			return env
		}
		env[i] = "PATH=" + extra + ":" + cur
		return env
	}
	return append(env, "PATH="+extra+":/usr/bin:/bin")
}

func scriptEnviron(p *config.Project) []string {
	env := augmentEnviron()
	if p != nil && strings.TrimSpace(p.Path) != "" {
		env = append(env, "TPSU_ROOT="+p.Path)
	}
	if ui.ColorEnabled() {
		env = append(env, "TPSU_COLOR=1")
	}
	if ui.SimplifiedOutput() {
		env = append(env, "TPS_SIMPLIFIED_OUTPUT=1")
	}
	return env
}

func pathHasPrefix(pathEnv, prefixes string) bool {
	for _, p := range strings.Split(prefixes, ":") {
		if p == "" {
			continue
		}
		if strings.HasPrefix(pathEnv, p+":") || pathEnv == p {
			return true
		}
	}
	return false
}
