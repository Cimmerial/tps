package pick

import (
	"fmt"
	"path/filepath"
	"strings"
)

// FileFilter limits file picks; zero value allows any extension.
type FileFilter struct {
	Extensions []string
}

func (f FileFilter) allowed(path string) bool {
	if len(f.Extensions) == 0 {
		return true
	}
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	for _, want := range f.Extensions {
		want = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(want), "."))
		if want != "" && ext == want {
			return true
		}
	}
	return false
}

// Validate reports whether path matches the filter extensions.
func (f FileFilter) Validate(path string) error {
	if f.allowed(path) {
		return nil
	}
	var parts []string
	for _, e := range f.Extensions {
		e = strings.TrimPrefix(strings.TrimSpace(e), ".")
		if e != "" {
			parts = append(parts, "."+e)
		}
	}
	return fmt.Errorf("choose a %s file (got %s)", strings.Join(parts, " or "), filepath.Base(path))
}

func applescriptTypeSuffix(extensions []string) string {
	if len(extensions) == 0 {
		return ""
	}
	seen := map[string]bool{}
	var items []string
	add := func(s string) {
		if s == "" || seen[s] {
			return
		}
		seen[s] = true
		items = append(items, `"`+s+`"`)
	}
	for _, e := range extensions {
		e = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(e), "."))
		add(e)
		if e == "pdf" {
			add("com.adobe.pdf")
		}
	}
	if len(items) == 0 {
		return ""
	}
	return " of type {" + strings.Join(items, ", ") + "}"
}
