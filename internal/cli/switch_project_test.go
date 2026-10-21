package cli

import (
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

func TestSwitchProjectByNameArg(t *testing.T) {
	cfg := &config.Config{
		Projects: []config.Project{
			{ID: "a", Name: "alpha", Path: "/a"},
			{ID: "b", Name: "beta", Path: "/b"},
		},
	}
	id, err := switchProject(nil, cfg, "a", "beta")
	if err != nil {
		t.Fatal(err)
	}
	if id != "b" {
		t.Fatalf("id=%q want b", id)
	}
}

func TestSwitchProjectByNameArgSameCurrent(t *testing.T) {
	cfg := &config.Config{
		Projects: []config.Project{
			{ID: "a", Name: "alpha", Path: "/a"},
		},
	}
	id, err := switchProject(nil, cfg, "a", "alpha")
	if err != nil {
		t.Fatal(err)
	}
	if id != "a" {
		t.Fatalf("id=%q want a (re-select current)", id)
	}
}

func TestIsProjectPickBack(t *testing.T) {
	for _, line := range []string{"", "   ", "q", "Q"} {
		if !isProjectPickBack(line) {
			t.Fatalf("line=%q want back", line)
		}
	}
	for _, line := range []string{"1", "beta", "quit"} {
		if isProjectPickBack(line) {
			t.Fatalf("line=%q want pick, not back", line)
		}
	}
}

func TestParsePickIndex(t *testing.T) {
	n, err := parsePickIndex("2", 3)
	if err != nil || n != 2 {
		t.Fatalf("n=%d err=%v", n, err)
	}
	_, err = parsePickIndex("0", 3)
	if err == nil {
		t.Fatal("want error for 0")
	}
	_, err = parsePickIndex("4", 3)
	if err == nil {
		t.Fatal("want error for 4")
	}
}
