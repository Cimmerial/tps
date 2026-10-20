package git

import (
	"github.com/Cimmerial/tps/internal/ui"
)

func quietGitArgs(args []string) []string {
	if !ui.SimplifiedOutput() || len(args) == 0 {
		return args
	}
	switch args[0] {
	case "commit":
		if !gitArgPresent(args, "-q", "--quiet") {
			return insertGitFlag(args, 1, "-q")
		}
	case "push":
		if !gitArgPresent(args, "-q", "--quiet") {
			return insertGitFlag(args, 1, "-q")
		}
	}
	return args
}

func gitArgPresent(args []string, flags ...string) bool {
	for _, a := range args[1:] {
		for _, f := range flags {
			if a == f {
				return true
			}
		}
	}
	return false
}

func insertGitFlag(args []string, after int, flag string) []string {
	if after < 0 || after > len(args) {
		return args
	}
	out := make([]string, 0, len(args)+1)
	out = append(out, args[:after]...)
	out = append(out, flag)
	out = append(out, args[after:]...)
	return out
}
