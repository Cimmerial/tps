package ui

import "testing"

func TestMoveWordLeft(t *testing.T) {
	buf := []rune("npm run build")
	cases := []struct {
		pos  int
		want int
	}{
		{0, 0},
		{3, 0},
		{4, 0},
		{7, 4},
		{8, 4},
		{13, 8},
		{len(buf), 8},
	}
	for _, c := range cases {
		got := moveWordLeft(buf, c.pos)
		if got != c.want {
			t.Fatalf("moveWordLeft(%q, %d) = %d want %d", string(buf), c.pos, got, c.want)
		}
	}
}

func TestMoveWordRight(t *testing.T) {
	buf := []rune("npm run build")
	cases := []struct {
		pos  int
		want int
	}{
		{0, 3},
		{3, 7},
		{4, 7},
		{7, len(buf)},
		{8, len(buf)},
		{len(buf), len(buf)},
	}
	for _, c := range cases {
		got := moveWordRight(buf, c.pos)
		if got != c.want {
			t.Fatalf("moveWordRight(%q, %d) = %d want %d", string(buf), c.pos, got, c.want)
		}
	}
}

func TestAppendLineHistory(t *testing.T) {
	var h []string
	appendLineHistory(&h, "a")
	appendLineHistory(&h, "a")
	appendLineHistory(&h, "b")
	appendLineHistory(&h, "")
	if len(h) != 2 || h[0] != "a" || h[1] != "b" {
		t.Fatalf("history: %v", h)
	}
}
