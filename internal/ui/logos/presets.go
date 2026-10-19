package logos

import (
	"crypto/rand"
	"math/big"

	"github.com/Cimmerial/tps/internal/config"
)

// Preset describes a title-card logo choice in the options menu.
type Preset struct {
	ID    string
	Label string
	Art   string
}

var catalog []Preset

func init() {
	catalog = buildCatalog()
}

func buildCatalog() []Preset {
	readme := readmeCatalog()
	out := make([]Preset, 0, len(readme)+1)
	for _, p := range readme {
		if p.ID == "default" {
			out = append(out, p)
			break
		}
	}
	out = append(out, Preset{
		ID:    "random",
		Label: "random (new pick each time)",
		Art:   "",
	})
	for _, p := range readme {
		if p.ID == "default" {
			continue
		}
		out = append(out, p)
	}
	return out
}

// Presets returns selectable logo presets: default, random, then README logos.
func Presets() []Preset {
	return catalog
}

// PickRandom chooses a logo id from the README catalog (not default/random).
func PickRandom() string {
	readme := readmeCatalog()
	if len(readme) == 0 {
		return "default"
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(readme))))
	if err != nil {
		return readme[0].ID
	}
	return readme[n.Int64()].ID
}

// Art returns ASCII art for id, or default when unknown. id "random" is resolved by caller.
func Art(id string) string {
	id = config.NormalizeLogoID(id)
	if id == "random" {
		id = PickRandom()
	}
	for _, p := range readmeCatalog() {
		if p.ID == id {
			return p.Art
		}
	}
	return readmeCatalog()[0].Art
}

// Label returns a human label for id.
func Label(id string) string {
	id = config.NormalizeLogoID(id)
	if id == "random" {
		return "random (new pick each time)"
	}
	for _, p := range Presets() {
		if p.ID == id {
			return p.Label
		}
	}
	if id == "" || id == "default" {
		return "default"
	}
	return id
}
