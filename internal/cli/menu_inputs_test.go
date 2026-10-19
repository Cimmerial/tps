package cli

import (
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

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
