package jobs

import (
	"runtime"
	"testing"
)

func TestLaunchReady(t *testing.T) {
	if launchReady(0, "tty1") {
		t.Fatal("pid 0 not ready")
	}
	if runtime.GOOS == "darwin" {
		if launchReady(42, "") {
			t.Fatal("darwin needs tty")
		}
		if !launchReady(42, "ttys001") {
			t.Fatal("darwin with pid+tty")
		}
	} else {
		if !launchReady(42, "") {
			t.Fatal("non-darwin pid only")
		}
	}
}
