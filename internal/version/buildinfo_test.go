package version

import "testing"

func TestFormatModuleVersion(t *testing.T) {
	tests := []struct {
		mod   string
		dirty bool
		want  string
	}{
		{"v0.1.0", false, "v0.1.0"},
		{"v0.1.1-0.20261020005918-e971f11656c0", true, "v0.1.1-dirty"},
		{"v0.1.1-0.20261020005918-e971f11656c0+dirty", false, "v0.1.1-dirty"},
		{"0.1.0", false, "v0.1.0"},
	}
	for _, tc := range tests {
		if got := formatModuleVersion(tc.mod, tc.dirty); got != tc.want {
			t.Fatalf("formatModuleVersion(%q, %v) = %q, want %q", tc.mod, tc.dirty, got, tc.want)
		}
	}
}
