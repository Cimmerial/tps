package cli

import (
	"os"
	"strings"

	"github.com/Cimmerial/tps/internal/update"
)

const flagUpdated = "--tps-updated"
const flagUpToDate = "--tps-uptodate"
const flagSessionProjectPrefix = "--tps-session-project="

var sessionRestartBanner string
var sessionRestartProjectID string

// peelRestartFlagsFromArgv returns "updated", "uptodate", or "" and strips internal flags from os.Args.
// Call from Execute before cobra parses flags.
func peelRestartFlagsFromArgv() string {
	status := ""
	var out []string
	out = append(out, os.Args[0])
	for i := 1; i < len(os.Args); i++ {
		switch {
		case os.Args[i] == flagUpdated:
			status = "updated"
		case os.Args[i] == flagUpToDate:
			status = "uptodate"
		case strings.HasPrefix(os.Args[i], flagSessionProjectPrefix):
			sessionRestartProjectID = strings.TrimPrefix(os.Args[i], flagSessionProjectPrefix)
		case os.Args[i] == update.FlagApplyStaged:
			continue
		default:
			out = append(out, os.Args[i])
		}
	}
	os.Args = out
	return status
}

func takeSessionRestartBanner() string {
	b := sessionRestartBanner
	sessionRestartBanner = ""
	return b
}

func takeSessionRestartProjectID() string {
	id := sessionRestartProjectID
	sessionRestartProjectID = ""
	return id
}
