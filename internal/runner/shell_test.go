package runner

import (
	"runtime"
	"testing"
)

func TestShellInvocation(t *testing.T) {
	bin, flag := ShellInvocation()
	switch runtime.GOOS {
	case "windows":
		if bin != "powershell" || flag != "-NoProfile" {
			t.Fatalf("got %s %s", bin, flag)
		}
	default:
		if bin != "sh" || flag != "-c" {
			t.Fatalf("got %s %s", bin, flag)
		}
	}
}
