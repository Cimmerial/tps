package browser

import "testing"

func TestNormalizeID(t *testing.T) {
	if NormalizeID("") != IDDefault {
		t.Fatal("empty")
	}
	if NormalizeID("default") != IDDefault {
		t.Fatal("default")
	}
	if NormalizeID("chrome") != IDChrome {
		t.Fatal("chrome")
	}
	if NormalizeID("Google Chrome") != IDChrome {
		t.Fatal("google chrome phrase")
	}
	if NormalizeID("msedge") != IDEdge {
		t.Fatal("msedge")
	}
}
