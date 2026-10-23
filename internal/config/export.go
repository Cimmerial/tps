package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const exportVersion = 1
const exportVersion2 = 2

type ExportBundle struct {
	Version     int          `yaml:"version"`
	Preferences *Preferences `yaml:"preferences,omitempty"`
	Project     *Project     `yaml:"project,omitempty"`
	Projects    []Project    `yaml:"projects,omitempty"`
}

// ExportKind selects default export filename shape.
type ExportKind int

const (
	ExportKindProject ExportKind = iota
	ExportKindOptions
	ExportKindEverything
)

func ExportProject(p *Project) ([]byte, error) {
	b := ExportBundle{Version: exportVersion, Project: p}
	return yaml.Marshal(b)
}

func ExportAllProjects(projects []Project) ([]byte, error) {
	b := ExportBundle{Version: exportVersion, Projects: projects}
	return yaml.Marshal(b)
}

func ExportProjectCommands(p *Project) ([]byte, error) {
	slim := slimProjectForExport(*p)
	b := ExportBundle{Version: exportVersion2, Project: &slim}
	return yaml.Marshal(b)
}

func ExportPreferencesBundle(prefs *Preferences, includeUsername bool) ([]byte, error) {
	b := ExportBundle{Version: exportVersion2, Preferences: CopyPortablePreferences(prefs, includeUsername)}
	return yaml.Marshal(b)
}

func ExportEverything(prefs *Preferences, projects []Project, includeUsername bool) ([]byte, error) {
	slim := make([]Project, 0, len(projects))
	for _, p := range projects {
		slim = append(slim, slimProjectForExport(p))
	}
	b := ExportBundle{
		Version:     exportVersion2,
		Preferences: CopyPortablePreferences(prefs, includeUsername),
		Projects:    slim,
	}
	return yaml.Marshal(b)
}

func slimProjectForExport(p Project) Project {
	return Project{
		Name:     p.Name,
		Commands: p.Commands,
	}
}

// CopyPortablePreferences clones user-facing prefs for export (no session project).
func CopyPortablePreferences(p *Preferences, includeUsername bool) *Preferences {
	if p == nil {
		return &Preferences{}
	}
	cp := *p
	cp.SessionProjectID = ""
	if !includeUsername {
		cp.Username = ""
	}
	return &cp
}

// ApplyPortablePreferences overwrites portable fields; session_project_id is unchanged.
// When includeUsername is false, dst.Username is left unchanged.
func ApplyPortablePreferences(dst, src *Preferences, includeUsername bool) {
	if dst == nil || src == nil {
		return
	}
	dst.Logo = src.Logo
	dst.TitleCard = src.TitleCard
	dst.TitleCardColor = src.TitleCardColor
	dst.Accent = src.Accent
	dst.ColorEnabled = src.ColorEnabled
	dst.SMSTo = src.SMSTo
	dst.SMSShortcut = src.SMSShortcut
	dst.PreferredIDE = src.PreferredIDE
	dst.PreferredBrowser = src.PreferredBrowser
	dst.SimplifiedOutput = src.SimplifiedOutput
	if includeUsername {
		dst.Username = src.Username
	}
}

// MergeCommandsByAlias updates matching aliases and appends new commands.
func MergeCommandsByAlias(dst *Project, imported []Command) {
	if dst == nil {
		return
	}
	for _, c := range imported {
		if existing := dst.FindCommand(c.Alias); existing != nil {
			*existing = c
			continue
		}
		dst.Commands = append(dst.Commands, c)
	}
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

// BundleProjectEntries returns command sets from a bundle (v1 or v2).
func (b ExportBundle) BundleProjectEntries() []Project {
	var out []Project
	if b.Project != nil && len(b.Project.Commands) > 0 {
		out = append(out, *b.Project)
	}
	for _, p := range b.Projects {
		if len(p.Commands) > 0 {
			out = append(out, p)
		}
	}
	return out
}

func (b ExportBundle) HasPreferences() bool {
	return b.Preferences != nil && b.Version >= exportVersion2
}

func (b ExportBundle) HasProjectCommands() bool {
	return len(b.BundleProjectEntries()) > 0
}

// CommandsForImport extracts commands from v1 full project or v2 name+commands.
func (b ExportBundle) CommandsForImport(projectName string) ([]Command, bool) {
	entries := b.BundleProjectEntries()
	if len(entries) == 0 {
		if b.Version < exportVersion2 && b.Project != nil && len(b.Project.Commands) > 0 {
			return b.Project.Commands, true
		}
		return nil, false
	}
	if len(entries) == 1 {
		return entries[0].Commands, true
	}
	projectName = strings.TrimSpace(projectName)
	for _, p := range entries {
		if strings.EqualFold(strings.TrimSpace(p.Name), projectName) {
			return p.Commands, true
		}
	}
	return nil, false
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

// ExportDefaultFilename builds tps-…-username.tps.yaml default names.
func ExportDefaultFilename(kind ExportKind, projectName, username string) string {
	user := SanitizeExportBasename(username)
	if user == "" {
		user = "user"
	}
	switch kind {
	case ExportKindOptions:
		return fmt.Sprintf("tps-options-%s.tps.yaml", user)
	case ExportKindEverything:
		return fmt.Sprintf("tps-all-%s.tps.yaml", user)
	default:
		proj := SanitizeExportBasename(projectName)
		return fmt.Sprintf("tps-%s-%s.tps.yaml", proj, user)
	}
}
