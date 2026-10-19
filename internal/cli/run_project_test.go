package cli

import "testing"

func TestJoinDetachedShellPrecursors(t *testing.T) {
	got := joinDetachedShellPrecursors([]string{"export FOO=1", "export BAR=2"}, "npm run x")
	want := "export FOO=1 && export BAR=2 && npm run x"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if joinDetachedShellPrecursors(nil, "npm run x") != "npm run x" {
		t.Fatal("expected shell unchanged with no precursors")
	}
}
