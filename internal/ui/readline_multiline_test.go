package ui

import "testing"

func TestShiftEnterInsertsNewline(t *testing.T) {
	if shiftEnterInsertsNewline("[13;2~", false) {
		t.Fatal("disabled multiline")
	}
	cases := []string{"[13;2~", "[13;2u", "[27;2;13~", "\x1b\r", "[13;5;2~"}
	for _, seq := range cases {
		if !shiftEnterInsertsNewline(seq, true) {
			t.Errorf("want shift+enter for %q", seq)
		}
	}
	if shiftEnterInsertsNewline("[13;5~", true) {
		t.Fatal("ctrl+enter should not match")
	}
}
