package jobs

import (
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

func TestJobDisplayLabel(t *testing.T) {
	if jobDisplayLabel(Job{Alias: "tunnel"}) != "tunnel" {
		t.Fatal("alias")
	}
	if jobDisplayLabel(Job{Shell: "npm run dev"}) != "npm run dev" {
		t.Fatal("shell")
	}
}

func TestRegisteredJobLabelsEmpty(t *testing.T) {
	if RegisteredJobLabels(nil) != nil {
		t.Fatal("nil cfg")
	}
	if len(RegisteredJobLabels(&config.Config{})) != 0 {
		t.Fatal("no projects")
	}
}
