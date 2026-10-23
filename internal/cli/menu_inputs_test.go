package cli

import (
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

func TestResolveMenuInputClearVsTps(t *testing.T) {
	in, err := resolveMenuInput(nil, &config.Project{}, "clear")
	if err != nil {
		t.Fatal(err)
	}
	if in.Action != ActionClear || !in.ClearScrollback {
		t.Fatalf("clear: %+v", in)
	}
	in, err = resolveMenuInput(nil, &config.Project{}, "tps")
	if err != nil {
		t.Fatal(err)
	}
	if in.Action != ActionClear || in.ClearScrollback {
		t.Fatalf("tps: %+v", in)
	}
	in, err = resolveMenuInput(nil, &config.Project{}, "")
	if err != nil {
		t.Fatal(err)
	}
	if in.Action != ActionClear || in.ClearScrollback {
		t.Fatalf("empty line: %+v", in)
	}
}

func TestResolveMenuInputVersion(t *testing.T) {
	for _, line := range []string{"ver", "version", "VER"} {
		in, err := resolveMenuInput(nil, &config.Project{}, line)
		if err != nil || in.Action != ActionVersion {
			t.Fatalf("%q: %+v err=%v", line, in, err)
		}
	}
}

func TestResolveMenuInputViteBuiltin(t *testing.T) {
	in, err := resolveMenuInput(nil, &config.Project{}, "vite")
	if err != nil || in.Action != ActionCloseAllTerms {
		t.Fatalf("vite: %+v err=%v", in, err)
	}
}

func TestResolveMenuInputProjectByIndex(t *testing.T) {
	cfg := &config.Config{
		Projects: []config.Project{
			{ID: "a", Name: "alpha", Path: "/a"},
			{ID: "b", Name: "beta", Path: "/b"},
		},
	}
	in, err := resolveMenuInput(cfg, &config.Project{}, "p1")
	if err != nil {
		t.Fatal(err)
	}
	if in.Action != ActionSwitchProject || in.Text != "a" {
		t.Fatalf("p1: %+v", in)
	}
	in, err = resolveMenuInput(cfg, &config.Project{}, "P2")
	if err != nil || in.Text != "b" {
		t.Fatalf("P2: %+v err=%v", in, err)
	}
	_, err = resolveMenuInput(cfg, &config.Project{}, "p9")
	if err == nil {
		t.Fatal("expected error for out of range project")
	}
}

func TestResolveMenuInputOpenInIDE(t *testing.T) {
	for _, line := range []string{"ide", "open in desired ide", "OPEN IDE"} {
		in, err := resolveMenuInput(nil, &config.Project{}, line)
		if err != nil || in.Action != ActionOpenInIDE {
			t.Fatalf("%q: %+v err=%v", line, in, err)
		}
	}
}

func TestResolveMenuInputProjectBack(t *testing.T) {
	in, err := resolveMenuInput(nil, &config.Project{}, "pb")
	if err != nil {
		t.Fatal(err)
	}
	if in.Action != ActionSwitchProjectBack {
		t.Fatalf("pb: %+v", in)
	}
}

func TestResolveMenuInputCommandWithShellInputs(t *testing.T) {
	p := &config.Project{
		Commands: []config.Command{
			{Alias: "inspect classes", Abbrev: "ic", Shell: "npm run inspect -- [INPUT1]"},
		},
	}
	in, err := resolveMenuInput(nil, p, "ic -Cimmerial")
	if err != nil {
		t.Fatal(err)
	}
	if in.Action != ActionRun || in.Alias != "inspect classes" {
		t.Fatalf("action/alias: %+v", in)
	}
	if len(in.ShellInputs) != 1 || in.ShellInputs[0] != "Cimmerial" {
		t.Fatalf("inputs: %v", in.ShellInputs)
	}
	_, err = resolveMenuInput(nil, p, "ic -a -b -c")
	if err == nil {
		t.Fatal("expected too many inputs")
	}
}

func TestResolveMenuInputDoc(t *testing.T) {
	for _, line := range []string{"doc", "documentation", "DOC"} {
		in, err := resolveMenuInput(nil, &config.Project{}, line)
		if err != nil || in.Action != ActionDoc {
			t.Fatalf("%q: %+v err=%v", line, in, err)
		}
	}
}

func TestResolveMenuInputNakedNC(t *testing.T) {
	for _, line := range []string{"nc", "NC", "new command"} {
		in, err := resolveMenuInput(nil, &config.Project{}, line)
		if err != nil {
			t.Fatalf("%q: %v", line, err)
		}
		if in.Action != ActionNewCommand || !in.NewCmdEditor {
			t.Fatalf("%q: %+v", line, in)
		}
	}
}

func TestResolveMenuInputInstallUtilities(t *testing.T) {
	in, err := resolveMenuInput(nil, &config.Project{}, "iu")
	if err != nil || in.Action != ActionInstallUtilities {
		t.Fatalf("iu: %+v err=%v", in, err)
	}
	in, err = resolveMenuInput(nil, &config.Project{}, "install utilities")
	if err != nil || in.Action != ActionInstallUtilities {
		t.Fatalf("install utilities: %+v err=%v", in, err)
	}
}
