package runner

import (
	"os/exec"
	"runtime"
)

// scriptCommand returns an exec.Cmd that runs script in the platform default shell.
func scriptCommand(script string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("powershell", "-NoProfile", "-Command", script)
	default:
		return exec.Command("sh", "-c", script)
	}
}

// ShellInvocation describes how RunScript invokes the platform shell (for tests).
func ShellInvocation() (bin string, flag string) {
	if runtime.GOOS == "windows" {
		return "powershell", "-NoProfile"
	}
	return "sh", "-c"
}
