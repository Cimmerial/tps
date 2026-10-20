package cli

import (
	"testing"
)

func TestSplitInputFieldsQuoted(t *testing.T) {
	fields := splitInputFields(`nc -n "push release" -a rel -c "bash scripts/push-release.sh [INPUT1]" -tags tps`)
	if len(fields) < 2 || fields[0] != "nc" {
		t.Fatalf("fields: %v", fields)
	}
	cmd, err := parseNewCommandFlags(fields[1:])
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Alias != "push release" {
		t.Fatalf("alias: %q", cmd.Alias)
	}
	if cmd.Abbrev != "rel" {
		t.Fatalf("abbrev: %q", cmd.Abbrev)
	}
	if cmd.Shell != "bash scripts/push-release.sh [INPUT1]" {
		t.Fatalf("shell: %q", cmd.Shell)
	}
	if len(cmd.Tags) != 1 || cmd.Tags[0] != "tps" {
		t.Fatalf("tags: %v", cmd.Tags)
	}
}

func TestParseNewCommandUnquoted(t *testing.T) {
	cmd, err := parseNewCommandFlags(splitInputFields(`nc -n push-release -a rel -c bash scripts/push-release.sh -tags tps`)[1:])
	if err != nil {
		t.Fatal(err)
	}
	if cmd.Alias != "push-release" {
		t.Fatalf("alias: %q", cmd.Alias)
	}
	if cmd.Shell != "bash scripts/push-release.sh" {
		t.Fatalf("shell: %q", cmd.Shell)
	}
}
