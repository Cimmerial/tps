package jobs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

func TestWatchCaptureActionsFromLog(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "job.log")
	if err := os.WriteFile(logPath, []byte("ready https://example.com/foo)\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Exercise extract path used by WatchCapture (no clipboard in unit test).
	m, err := config.ExtractCapture(string(mustRead(logPath)), config.ResolveCapturePattern("url"))
	if err != nil || m != "https://example.com/foo" {
		t.Fatalf("match=%q err=%v", m, err)
	}
}

func mustRead(path string) []byte {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}
