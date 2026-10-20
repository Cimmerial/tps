package update

import (
	"os"
	"os/exec"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

const flagUpdated = "--tps-updated"
const flagUpToDate = "--tps-uptodate"
const flagSessionProjectPrefix = "--tps-session-project="

func buildRestartArgv() (bin string, argv []string, err error) {
	bin, err = exec.LookPath("tps")
	if err != nil {
		exe, err := os.Executable()
		if err != nil {
			return "", nil, err
		}
		bin = exe
	}

	argv = []string{bin}
	sessionProject := ""
	if prefs, err := config.LoadPreferences(); err == nil && prefs != nil {
		sessionProject = strings.TrimSpace(prefs.SessionProjectID)
	}
	for i := 1; i < len(os.Args); i++ {
		a := os.Args[i]
		if a == flagUpdated || a == flagUpToDate {
			continue
		}
		if strings.HasPrefix(a, flagSessionProjectPrefix) {
			continue
		}
		argv = append(argv, a)
	}
	if sessionProject != "" {
		argv = append(argv, flagSessionProjectPrefix+sessionProject)
	}
	return bin, argv, nil
}
