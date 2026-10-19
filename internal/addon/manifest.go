package addon

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"gopkg.in/yaml.v3"
)

const (
	UtilitiesID       = "tps-utilities"
	UtilitiesDirName  = "tps-utilities"
	UtilitiesCloneURL = "https://github.com/Cimmerial/tps-utilities.git"
	ManifestName      = "tps-addon.yaml"
)

type Manifest struct {
	ID                 string           `yaml:"id"`
	Name               string           `yaml:"name"`
	CloneURL           string           `yaml:"clone_url"`
	ContractIteration  int              `yaml:"contract_iteration"`
	MinTPS             string           `yaml:"min_tps"`
	Commands           []config.Command `yaml:"commands"`
}

func ManifestPath(projectRoot string) string {
	return filepath.Join(projectRoot, ManifestName)
}

func LoadManifest(projectRoot string) (*Manifest, error) {
	path := ManifestPath(projectRoot)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func IsUtilitiesProject(p *config.Project) bool {
	if p == nil {
		return false
	}
	m, err := LoadManifest(p.Path)
	if err != nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(m.ID), UtilitiesID)
}

// MergeManifestCommands adds or updates commands from the project manifest.
func MergeManifestCommands(p *config.Project) error {
	if p == nil {
		return fmt.Errorf("no project")
	}
	m, err := LoadManifest(p.Path)
	if err != nil {
		return err
	}
	for _, cmd := range m.Commands {
		if err := p.AddCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

func DefaultUtilitiesManifest() Manifest {
	return Manifest{
		ID:                UtilitiesID,
		Name:              UtilitiesDirName,
		CloneURL:          UtilitiesCloneURL,
		ContractIteration: CurrentContractIteration,
		Commands: []config.Command{
			{
				Alias:       "universal image conversion",
				Abbrev:      "uic",
				Description: "convert images via PNG hub (JPEG, PNG, GIF, …)",
				Menu:        "menus/uic.yaml",
			},
		},
	}
}
