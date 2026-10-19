package project

import (
	"errors"
	"fmt"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/google/uuid"
)

type ResolveResult struct {
	Project *config.Project
	Healed  bool
}

func ResolveFromCwd(cfg *config.Config, cwd string) (*ResolveResult, error) {
	root := cwd
	if gitRoot, err := GitRoot(cwd); err == nil && gitRoot != "" {
		root = gitRoot
	}

	cleanRoot, err := config.CleanPath(root)
	if err != nil {
		return nil, err
	}

	if p := cfg.FindByPath(cleanRoot); p != nil {
		return &ResolveResult{Project: p, Healed: false}, nil
	}

	remote, _ := GitRemote(cleanRoot)
	if remote != "" {
		if p := cfg.FindByGitRemote(remote); p != nil {
			stored, err := config.CleanPath(p.Path)
			if err != nil {
				return nil, err
			}
			healed := stored != cleanRoot
			if healed {
				p.Path = cleanRoot
			}
			return &ResolveResult{Project: p, Healed: healed}, nil
		}
	}

	return nil, fmt.Errorf("no linked project for %s (run: tps project add)", cleanRoot)
}

func ProjectRootFromCwd(cwd string) (string, error) {
	if gitRoot, err := GitRoot(cwd); err == nil && gitRoot != "" {
		return config.CleanPath(gitRoot)
	}
	return config.CleanPath(cwd)
}

func LinkFromCwd(cfg *config.Config, cwd string, name string) (*config.Project, error) {
	root, err := ProjectRootFromCwd(cwd)
	if err != nil {
		return nil, err
	}
	if name == "" {
		return nil, errors.New("project name is required")
	}

	remote, _ := GitRemote(root)
	p := config.Project{
		Name:      name,
		Path:      root,
		GitRemote: remote,
	}

	if remote != "" {
		if existing := cfg.FindByGitRemote(remote); existing != nil {
			p.ID = existing.ID
			p.Commands = existing.Commands
			p.CommandUsage = existing.CommandUsage
			p.MetaUsage = existing.MetaUsage
		}
	}
	if p.ID == "" {
		if existing := cfg.FindByPath(root); existing != nil {
			p.ID = existing.ID
			p.Commands = existing.Commands
			p.CommandUsage = existing.CommandUsage
			p.MetaUsage = existing.MetaUsage
		}
	}
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	cfg.UpsertProject(p)
	linked := cfg.FindByPath(root)
	if linked == nil && p.ID != "" {
		linked = cfg.FindByID(p.ID)
	}
	if linked == nil {
		return nil, errors.New("failed to link project")
	}
	return linked, nil
}
