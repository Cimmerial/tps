package git

import (
	"path"
	"strings"
)

// SuspiciousPaths returns paths from names that look like secrets or credentials.
// Only paths that would be included in a commit should be passed in (e.g. staged files).
func SuspiciousPaths(names []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" || seen[name] {
			continue
		}
		if pathLooksSuspicious(name) {
			seen[name] = true
			out = append(out, name)
		}
	}
	return out
}

func pathLooksSuspicious(p string) bool {
	p = strings.TrimSpace(p)
	if p == "" {
		return false
	}
	p = path.Clean(p)
	slash := strings.ReplaceAll(p, "\\", "/")
	lower := strings.ToLower(slash)

	base := path.Base(slash)
	baseLower := strings.ToLower(base)
	if base == ".env" || strings.HasPrefix(base, ".env.") {
		return true
	}
	if strings.HasSuffix(lower, ".pem") {
		return true
	}
	if strings.HasSuffix(lower, ".p12") || strings.HasSuffix(lower, ".pfx") {
		return true
	}
	if baseLower == "credentials" || strings.HasPrefix(baseLower, "credentials.") {
		return true
	}
	if strings.Contains(baseLower, "id_rsa") || strings.Contains(baseLower, "id_dsa") {
		return true
	}

	for _, seg := range strings.Split(slash, "/") {
		if seg == "" || seg == "." {
			continue
		}
		sl := strings.ToLower(seg)
		if sl == "credentials" || strings.HasPrefix(sl, "credentials.") {
			return true
		}
		if strings.Contains(sl, "secret") || strings.Contains(sl, "private") {
			return true
		}
	}
	return false
}
