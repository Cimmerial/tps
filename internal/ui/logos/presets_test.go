package logos

import "testing"

func TestArtResolution(t *testing.T) {
	if Art("default") != Art("") {
		t.Fatal("empty should use default art")
	}
	if Art("unknown") == "" {
		t.Fatal("unknown should fall back to default art")
	}
	if Art("tps_dots") != Art("tps_slash") {
		t.Fatal("tps_dots alias")
	}
	if Art("pixel") == Art("default") {
		t.Fatal("pixel should differ from default")
	}
}
