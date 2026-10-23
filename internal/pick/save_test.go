package pick

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSavePath(t *testing.T) {
	dir := t.TempDir()
	name := "tps-demo-user.tps.yaml"
	got := ResolveSavePath(dir, name)
	want := filepath.Join(dir, name)
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	file := filepath.Join(dir, "out.yaml")
	if ResolveSavePath(file, name) != file {
		t.Fatal("file path unchanged")
	}
	trail := dir + string(os.PathSeparator)
	if ResolveSavePath(trail, name) != want {
		t.Fatalf("trailing sep: %q", ResolveSavePath(trail, name))
	}
}
