package git

import (
	"os"
	"strings"
)

// AppendGitignoreLines adds paths to dir/.gitignore when not already present.
func AppendGitignoreLines(dir string, paths []string) error {
	return appendGitignoreLines(dir, paths)
}

func appendGitignoreLines(dir string, paths []string) error {
	if len(paths) == 0 {
		return nil
	}
	gitignorePath := dir + string(os.PathSeparator) + ".gitignore"
	var existing []byte
	if data, err := os.ReadFile(gitignorePath); err == nil {
		existing = data
	} else if !os.IsNotExist(err) {
		return err
	}
	content := string(existing)
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	lines := gitignoreLines(content)
	added := false
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" || gitignoreLinesHas(lines, p) {
			continue
		}
		content += p + "\n"
		lines = append(lines, p)
		added = true
	}
	if !added {
		return nil
	}
	return os.WriteFile(gitignorePath, []byte(content), 0o644)
}

func gitignoreLines(content string) []string {
	var out []string
	for _, l := range strings.Split(content, "\n") {
		out = append(out, strings.TrimSpace(l))
	}
	return out
}

func gitignoreLinesHas(lines []string, path string) bool {
	for _, l := range lines {
		if l == path {
			return true
		}
	}
	return false
}
