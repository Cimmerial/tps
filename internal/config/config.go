package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "tps", "config.yaml"), nil
}

func Load() (*Config, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		normalized := normalizeYAMLBytes(data)
		if err2 := yaml.Unmarshal(normalized, &cfg); err2 != nil {
			return nil, fmt.Errorf("parse config: %w (use spaces not tabs; run: tps config open)", err)
		}
		_ = os.WriteFile(p, normalized, 0o600)
	}
	return &cfg, nil
}

func normalizeYAMLBytes(data []byte) []byte {
	s := strings.ReplaceAll(string(data), "\t", "  ")
	return []byte(s)
}

func Save(c *Config) error {
	p, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o600)
}

func CleanPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

func (c *Config) FindByPath(path string) *Project {
	clean, err := CleanPath(path)
	if err != nil {
		return nil
	}
	for i := range c.Projects {
		pClean, err := CleanPath(c.Projects[i].Path)
		if err != nil {
			continue
		}
		if pClean == clean {
			return &c.Projects[i]
		}
	}
	return nil
}

func (c *Config) FindByGitRemote(remote string) *Project {
	if remote == "" {
		return nil
	}
	for i := range c.Projects {
		if c.Projects[i].GitRemote == remote {
			return &c.Projects[i]
		}
	}
	return nil
}

func (c *Config) FindByID(id string) *Project {
	for i := range c.Projects {
		if c.Projects[i].ID == id {
			return &c.Projects[i]
		}
	}
	return nil
}

func mergeProjectMeta(dst *Project, src Project) {
	dst.ID = src.ID
	dst.Name = src.Name
	dst.Path = src.Path
	dst.GitRemote = src.GitRemote
	dst.Commands = src.Commands
	if len(src.CommandUsage) > 0 {
		dst.CommandUsage = src.CommandUsage
	}
	if len(src.MetaUsage) > 0 {
		dst.MetaUsage = src.MetaUsage
	}
}

func (c *Config) UpsertProject(p Project) {
	if existing := c.FindByID(p.ID); existing != nil {
		if len(p.CommandUsage) == 0 {
			p.CommandUsage = existing.CommandUsage
		}
		if len(p.MetaUsage) == 0 {
			p.MetaUsage = existing.MetaUsage
		}
		mergeProjectMeta(existing, p)
		return
	}
	if p.GitRemote != "" {
		if existing := c.FindByGitRemote(p.GitRemote); existing != nil {
			p.ID = existing.ID
			if len(p.CommandUsage) == 0 {
				p.CommandUsage = existing.CommandUsage
			}
			if len(p.MetaUsage) == 0 {
				p.MetaUsage = existing.MetaUsage
			}
			mergeProjectMeta(existing, p)
			return
		}
	}
	if existing := c.FindByPath(p.Path); existing != nil {
		p.ID = existing.ID
		if len(p.CommandUsage) == 0 {
			p.CommandUsage = existing.CommandUsage
		}
		if len(p.MetaUsage) == 0 {
			p.MetaUsage = existing.MetaUsage
		}
		mergeProjectMeta(existing, p)
		return
	}
	c.Projects = append(c.Projects, p)
}

func (p *Project) FindCommand(alias string) *Command {
	for i := range p.Commands {
		if p.Commands[i].Alias == alias {
			return &p.Commands[i]
		}
	}
	return nil
}

func (p *Project) FindCommandByAbbrev(abbrev string) *Command {
	if abbrev == "" {
		return nil
	}
	for i := range p.Commands {
		if p.Commands[i].Abbrev != "" && strings.EqualFold(p.Commands[i].Abbrev, abbrev) {
			return &p.Commands[i]
		}
	}
	return nil
}

func (p *Project) AddCommand(cmd Command) error {
	if cmd.Alias == "" {
		return errors.New("command alias is required")
	}
	if cmd.Shell == "" {
		return errors.New("command shell is required")
	}
	for i := range p.Commands {
		if p.Commands[i].Alias == cmd.Alias {
			p.Commands[i] = cmd
			return nil
		}
	}
	p.Commands = append(p.Commands, cmd)
	return nil
}

func (p *Project) RemoveCommand(alias string) bool {
	for i := range p.Commands {
		if p.Commands[i].Alias == alias {
			p.Commands = append(p.Commands[:i], p.Commands[i+1:]...)
			if p.CommandUsage != nil {
				delete(p.CommandUsage, alias)
			}
			return true
		}
	}
	return false
}
