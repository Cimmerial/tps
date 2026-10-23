package instance

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/git"
)

// ReadmeRelSource is the canonical doc path in a tps dev clone (for editors).
const ReadmeRelSource = "internal/instance/README.md"

// ReadmeDestName is written at the root of each linked app repo.
const ReadmeDestName = "TPS.md"

// SpecsRelSource is the canonical spec path in a tps dev clone (for editors).
const SpecsRelSource = "internal/instance/SPECS.md"

// SpecsDestName is written at the root of each linked app repo.
const SpecsDestName = "TPS_SPECS.md"

func IsDevRepo(projectRoot string) bool {
	if projectRoot == "" {
		return false
	}
	_, err := os.Stat(filepath.Join(projectRoot, "cmd", "tps"))
	return err == nil
}

func DevRepoRoot(cfg *config.Config) string {
	if cfg == nil {
		return ""
	}
	for _, p := range cfg.Projects {
		if IsDevRepo(p.Path) {
			return p.Path
		}
	}
	return ""
}

func readmeBytes(tpsRepo string) ([]byte, error) {
	if tpsRepo != "" {
		for _, rel := range []string{ReadmeRelSource, "instance/README.md"} {
			src := filepath.Join(tpsRepo, rel)
			data, err := os.ReadFile(src)
			if err == nil {
				return data, nil
			}
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("read instance readme: %w", err)
			}
		}
	}
	if len(EmbeddedReadme()) > 0 {
		return EmbeddedReadme(), nil
	}
	return nil, nil
}

func specsBytes(tpsRepo string) ([]byte, error) {
	if tpsRepo != "" {
		for _, rel := range []string{SpecsRelSource, "instance/SPECS.md"} {
			src := filepath.Join(tpsRepo, rel)
			data, err := os.ReadFile(src)
			if err == nil {
				return data, nil
			}
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("read instance specs: %w", err)
			}
		}
	}
	if len(EmbeddedSpecs()) > 0 {
		return EmbeddedSpecs(), nil
	}
	return nil, nil
}

func SyncProject(tpsRepo, projectRoot string) (bool, error) {
	if projectRoot == "" {
		return false, nil
	}
	if IsDevRepo(projectRoot) {
		return false, nil
	}
	readme, err := readmeBytes(tpsRepo)
	if err != nil {
		return false, err
	}
	specs, err := specsBytes(tpsRepo)
	if err != nil {
		return false, err
	}
	wrote := false
	if len(readme) > 0 {
		dst := filepath.Join(projectRoot, ReadmeDestName)
		if err := os.WriteFile(dst, readme, 0o644); err != nil {
			return false, fmt.Errorf("write %s: %w", ReadmeDestName, err)
		}
		wrote = true
	}
	if len(specs) > 0 {
		dst := filepath.Join(projectRoot, SpecsDestName)
		if err := os.WriteFile(dst, specs, 0o644); err != nil {
			return false, fmt.Errorf("write %s: %w", SpecsDestName, err)
		}
		wrote = true
	}
	if err := ensureInstanceDocsGitignored(projectRoot); err != nil {
		return wrote, err
	}
	return wrote, nil
}

func ensureInstanceDocsGitignored(projectRoot string) error {
	return git.AppendGitignoreLines(projectRoot, []string{ReadmeDestName, SpecsDestName})
}

func SyncAllLinked(tpsRepo string, cfg *config.Config) (int, error) {
	if cfg == nil {
		return 0, nil
	}
	n := 0
	for _, p := range cfg.Projects {
		wrote, err := SyncProject(tpsRepo, p.Path)
		if err != nil {
			return n, err
		}
		if wrote {
			n++
		}
	}
	return n, nil
}
