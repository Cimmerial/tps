package addonmenu

import "testing"

func TestParseChoice(t *testing.T) {
	n, err := parseChoice("2", 3)
	if err != nil || n != 2 {
		t.Fatalf("got %d err=%v", n, err)
	}
	_, err = parseChoice("9", 3)
	if err == nil {
		t.Fatal("expected range error")
	}
}
