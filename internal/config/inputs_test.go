package config

import "testing"

func TestParseShellInputTokens(t *testing.T) {
	got, err := ParseShellInputTokens([]string{"-Cimmerial", "-two"})
	if err != nil || len(got) != 2 || got[0] != "Cimmerial" || got[1] != "two" {
		t.Fatalf("parse: %v %v", got, err)
	}
	_, err = ParseShellInputTokens([]string{"nope"})
	if err == nil {
		t.Fatal("expected error for non - token")
	}
}

func TestSubstituteShellInputs(t *testing.T) {
	shell := "npm run inspect:library-classes -- --username [INPUT1]"
	got, err := SubstituteShellInputs(shell, []string{"Cimmerial"})
	if err != nil {
		t.Fatal(err)
	}
	want := "npm run inspect:library-classes -- --username Cimmerial"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestValidateCommandInputs(t *testing.T) {
	c := &Command{
		Alias: "ic",
		Shell: "echo [INPUT1] [INPUT2]",
	}
	if err := ValidateCommandInputs(c, []string{"a"}); err == nil {
		t.Fatal("expected missing input error")
	}
	if err := ValidateCommandInputs(c, []string{"a", "b", "c"}); err == nil {
		t.Fatal("expected too many inputs")
	}
	if err := ValidateCommandInputs(c, []string{"a", "b"}); err != nil {
		t.Fatal(err)
	}
	plain := &Command{Alias: "r", Shell: "npm run dev"}
	if err := ValidateCommandInputs(plain, []string{"x"}); err == nil {
		t.Fatal("expected extra input error")
	}
}

func TestShellsWithInputs(t *testing.T) {
	c := &Command{
		Alias: "ic",
		Shell: "npm run inspect -- [INPUT1]",
	}
	got, err := c.ShellsWithInputs([]string{"user1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "npm run inspect -- user1" {
		t.Fatalf("got %v", got)
	}
}
