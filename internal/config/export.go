package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const exportVersion = 1

type ExportBundle struct {
	Version  int       `yaml:"version"`
	Project  *Project  `yaml:"project,omitempty"`
	Projects []Project `yaml:"projects,omitempty"`
}

func ExportProject(p *Project) ([]byte, error) {
	b := ExportBundle{Version: exportVersion, Project: p}
	return yaml.Marshal(b)
}

func ExportAllProjects(projects []Project) ([]byte, error) {
	b := ExportBundle{Version: exportVersion, Projects: projects}
	return yaml.Marshal(b)
}

func WriteExportFile(path string, data []byte) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("export path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil && filepath.Dir(path) != "." {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func ImportFromFile(path string) (ExportBundle, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ExportBundle{}, err
	}
	var b ExportBundle
	if err := yaml.Unmarshal(data, &b); err != nil {
		return ExportBundle{}, fmt.Errorf("parse export: %w", err)
	}
	if b.Version == 0 {
		b.Version = 1
	}
	return b, nil
}

func (c *Config) MergeImport(b ExportBundle) (int, error) {
	merged := 0
	if b.Project != nil {
		if err := c.upsertProject(*b.Project); err != nil {
			return merged, err
		}
		merged++
	}
	for _, p := range b.Projects {
		if err := c.upsertProject(p); err != nil {
			return merged, err
		}
		merged++
	}
	if merged == 0 {
		return 0, fmt.Errorf("import file has no project data")
	}
	return merged, nil
}

func (c *Config) upsertProject(p Project) error {
	if strings.TrimSpace(p.ID) == "" {
		return fmt.Errorf("project %q missing id", p.Name)
	}
	if i := c.FindByID(p.ID); i != nil {
		*i = p
		return nil
	}
	c.Projects = append(c.Projects, p)
	return nil
}

func SanitizeExportBasename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "project"
	}
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "project"
	}
	return out
}
