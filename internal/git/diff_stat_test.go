package git

import "testing"

func TestParseDiffNumstat(t *testing.T) {
	out := "12\t3\tinternal/foo.go\n-\t-\timage.png\n0\t5\tdeleted.go\n"
	m := parseDiffNumstat([]byte(out))
	if d := m["internal/foo.go"]; d.Added != 12 || d.Removed != 3 || d.Binary {
		t.Fatalf("foo: %+v", d)
	}
	if !m["image.png"].Binary {
		t.Fatal("expected binary")
	}
	if d := m["deleted.go"]; d.Added != 0 || d.Removed != 5 {
		t.Fatalf("deleted: %+v", d)
	}
}
