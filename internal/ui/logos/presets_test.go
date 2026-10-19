package logos

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Cimmerial/tps/internal/config"
)

func TestArtResolution(t *testing.T) {
	if Art("default") != Art("") {
		t.Fatal("empty should use default art")
	}
	if Art("unknown") == "" {
		t.Fatal("unknown should fall back to default art")
	}
	if Art("block") == Art("default") {
		t.Fatal("block should differ from default art")
	}
	if config.NormalizeLogoID("pixel") != "block" {
		t.Fatal("pixel legacy alias")
	}
}

func TestPresetsOrder(t *testing.T) {
	list := Presets()
	if len(list) < 3 {
		t.Fatalf("len=%d", len(list))
	}
	if list[0].ID != "default" || list[1].ID != "random" {
		t.Fatalf("order: %s, %s", list[0].ID, list[1].ID)
	}
}

func TestREADMECatalogCount(t *testing.T) {
	readme := readmeCatalog()
	if len(readme) != 23 {
		t.Fatalf("readme logos=%d want 23", len(readme))
	}
}

func TestREADMECatalogMatchesFile(t *testing.T) {
	path := repoREADMEPath()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skip("README not available:", err)
	}
	parsed := parseREADMELogos(string(data))
	if len(parsed) != len(readmeCatalog()) {
		t.Fatalf("parsed=%d embedded=%d — regenerate catalog_readme.go from README.md",
			len(parsed), len(readmeCatalog()))
	}
	for i := range parsed {
		if parsed[i].ID != readmeCatalog()[i].ID {
			t.Fatalf("id[%d]=%q embedded=%q", i, parsed[i].ID, readmeCatalog()[i].ID)
		}
		if strings.TrimSpace(parsed[i].Art) != strings.TrimSpace(readmeCatalog()[i].Art) {
			t.Fatalf("art drift at %q — regenerate catalog_readme.go", parsed[i].ID)
		}
	}
}

func repoREADMEPath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "README.md"
	}
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "README.md")
}
