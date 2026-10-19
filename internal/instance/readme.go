package instance

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cimmerial/tps/internal/config"
)

// ReadmeRelSource is the canonical doc in the tps dev repo. Update it when behavior or limits change.
const ReadmeRelSource = "instance/README.md"

// ReadmeDestName is written at the root of each linked app repo.
const ReadmeDestName = "TPS.md"

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

func SyncProject(tpsRepo, projectRoot string) (bool, error) {
	if tpsRepo == "" || projectRoot == "" {
		return false, nil
	}
	if IsDevRepo(projectRoot) {
		return false, nil
	}
	src := filepath.Join(tpsRepo, ReadmeRelSource)
	data, err := os.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("read instance readme: %w", err)
	}
	dst := filepath.Join(projectRoot, ReadmeDestName)
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		return false, fmt.Errorf("write %s: %w", ReadmeDestName, err)
	}
	return true, nil
}

func SyncAllLinked(tpsRepo string, cfg *config.Config) (int, error) {
	if cfg == nil || tpsRepo == "" {
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
