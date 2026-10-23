package cli

import "testing"

func TestSumShipFileDiffs(t *testing.T) {
	diffs := map[string]shipFileDiff{
		"a": {Added: 10, Removed: 2},
		"b": {Added: 3, Removed: 0},
	}
	add, rem := sumShipFileDiffs([]string{"a", "b"}, diffs)
	if add != 13 || rem != 2 {
		t.Fatalf("got +%d -%d", add, rem)
	}
}
