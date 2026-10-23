package git

import "testing"

func TestParsePorcelainV1Null(t *testing.T) {
	data := " M internal/foo.go\000?? new.txt\000"
	got := parsePorcelainV1Null([]byte(data))
	want := []string{"internal/foo.go", "new.txt"}
	if len(got) != len(want) {
		t.Fatalf("got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("[%d] got %q want %q", i, got[i], want[i])
		}
	}
}

func TestParsePorcelainV1NullRename(t *testing.T) {
	data := "R  old.go\000new.go\000"
	got := parsePorcelainV1Null([]byte(data))
	if len(got) != 1 || got[0] != "new.go" {
		t.Fatalf("got %v", got)
	}
}

func TestParsePorcelainV1NullQuoted(t *testing.T) {
	data := "?? \"path with spaces.txt\"\000"
	got := parsePorcelainV1Null([]byte(data))
	if len(got) != 1 || got[0] != "path with spaces.txt" {
		t.Fatalf("got %v", got)
	}
}

func TestShipFileSelectionToggle(t *testing.T) {
	sel := NewShipFileSelection([]string{"a", "b", "c"})
	if !sel.Toggle(2) {
		t.Fatal("toggle 2")
	}
	if sel.Included[1] {
		t.Fatal("b should be excluded")
	}
	inc := sel.IncludedPaths()
	if len(inc) != 2 || inc[0] != "a" || inc[1] != "c" {
		t.Fatalf("included %v", inc)
	}
	exc := sel.ExcludedPaths()
	if len(exc) != 1 || exc[0] != "b" {
		t.Fatalf("excluded %v", exc)
	}
}

func TestParseShipFileSelectInput(t *testing.T) {
	cases := []struct {
		in     string
		commit bool
		quit   bool
		num    int
		ok     bool
	}{
		{"", false, true, 0, true},
		{"  ", false, true, 0, true},
		{"0", false, true, 0, true},
		{"c", true, false, 0, true},
		{"C", true, false, 0, true},
		{"1", false, false, 1, true},
		{"12", false, false, 12, true},
		{"x", false, false, 0, false},
	}
	for _, tc := range cases {
		commit, quit, num, ok := ParseShipFileSelectInput(tc.in)
		if commit != tc.commit || quit != tc.quit || num != tc.num || ok != tc.ok {
			t.Errorf("ParseShipFileSelectInput(%q) = (%v,%v,%d,%v) want (%v,%v,%d,%v)",
				tc.in, commit, quit, num, ok, tc.commit, tc.quit, tc.num, tc.ok)
		}
	}
}
