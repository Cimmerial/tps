package cli

import "os"

const flagUpdated = "--tps-updated"
const flagUpToDate = "--tps-uptodate"

var sessionRestartBanner string

// peelRestartFlagsFromArgv returns "updated", "uptodate", or "" and strips internal flags from os.Args.
// Call from Execute before cobra parses flags.
func peelRestartFlagsFromArgv() string {
	status := ""
	var out []string
	out = append(out, os.Args[0])
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case flagUpdated:
			status = "updated"
		case flagUpToDate:
			status = "uptodate"
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
