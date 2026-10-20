//go:build darwin

package runner

import (
	"strings"
	"testing"
)

func TestAugmentEnvironPrependsHomebrew(t *testing.T) {
	env := augmentEnviron()
	var path string
	for _, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			path = strings.TrimPrefix(e, "PATH=")
			break
		}
	}
	if path == "" {
		t.Fatal("no PATH")
	}
	if !strings.HasPrefix(path, "/opt/homebrew/bin") {
		t.Fatalf("PATH=%s", path)
	}
}
