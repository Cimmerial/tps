package platform

import (
	"runtime"
	"testing"
)

func TestSupportsCaptureSMS(t *testing.T) {
	got := SupportsCaptureSMS()
	want := runtime.GOOS == "darwin"
	if got != want {
		t.Fatalf("SupportsCaptureSMS()=%v want %v", got, want)
	}
}
