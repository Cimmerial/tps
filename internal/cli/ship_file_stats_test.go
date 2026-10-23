package cli

import "testing"

func TestCountFileLines(t *testing.T) {
	if countFileLines(nil) != 0 {
		t.Fatal("empty")
	}
	if countFileLines([]byte("a\nb\n")) != 2 {
		t.Fatal("two lines")
	}
	if countFileLines([]byte("a\nb")) != 2 {
		t.Fatal("no trailing newline")
	}
}

func TestSumShipFileStats(t *testing.T) {
	stats := map[string]shipFileStat{
		"a": {Lines: 10, Bytes: 100},
		"b": {Lines: 5, Bytes: 50},
	}
	lines, bytes := sumShipFileStats([]string{"a", "b"}, stats)
	if lines != 15 || bytes != 150 {
		t.Fatalf("got %d lines %d bytes", lines, bytes)
	}
}

func TestFormatShipFileSize(t *testing.T) {
	if formatShipFileSize(512) != "512 B" {
		t.Fatal("bytes")
	}
	if formatShipFileSize(2048) != "2 KB" {
		t.Fatal("kb")
	}
}
