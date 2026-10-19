package runner

import (
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

func TestDetachInfoDevMobile(t *testing.T) {
	c := config.Command{
		Alias: "run",
		Shell: "npm run dev:mobile",
	}
	if !ShouldDetach(c) {
		t.Fatal("expected detach")
	}
	_, hint := DetachInfo(c)
	if hint == "" || hint == "runs in this session (foreground)" {
		t.Fatalf("hint=%q", hint)
	}
}

func TestDetachInfoForceForeground(t *testing.T) {
	c := config.Command{
		Alias:           "run",
		Shell:           "npm run dev:mobile",
		ForceForeground: true,
	}
	if ShouldDetach(c) {
		t.Fatal("foreground should override auto detach")
	}
}
