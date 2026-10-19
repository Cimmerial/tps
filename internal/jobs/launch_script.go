package jobs

import (
	"os"
	"strings"
)

func writeRunScript(path, workDir, shellCmd, logFile string) error {
	var b strings.Builder
	b.WriteString("#!/usr/bin/env bash\nset -e\n")
	b.WriteString("cd ")
	b.WriteString(shSingleQuote(workDir))
	b.WriteString(" || exit 1\n")
	if logFile != "" {
		b.WriteString("export TPS_CAPTURE_LOG=")
		b.WriteString(shSingleQuote(logFile))
		b.WriteString("\n")
		b.WriteString("if [[ -n \"$TPS_CAPTURE_LOG\" ]]; then exec > >(tee \"$TPS_CAPTURE_LOG\") 2>&1; fi\n")
	}
	b.WriteString(strings.TrimSpace(shellCmd))
	b.WriteString("\n")
	return os.WriteFile(path, []byte(b.String()), 0o755)
}

func shSingleQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\"'\"'") + "'"
}
