package update

import (
	"runtime"
	"strings"
	"testing"
)

func TestReleaseAssetName(t *testing.T) {
	name := ReleaseAssetName("1.2.3")
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(name, "windows_amd64.zip") {
			t.Fatalf("got %s", name)
		}
	} else {
		if !strings.Contains(name, runtime.GOOS) || !strings.HasSuffix(name, ".tar.gz") {
			t.Fatalf("got %s", name)
		}
	}
	if !strings.HasPrefix(name, "tps_1.2.3_") {
		t.Fatalf("got %s", name)
	}
}
