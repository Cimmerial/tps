package jobs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteRunScript(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "run.sh")
	work := "/tmp/project"
	if err := writeRunScript(path, work, "echo hi", ""); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	if !strings.Contains(text, work) || !strings.Contains(text, "echo hi") {
		t.Fatalf("script: %s", text)
	}
}
