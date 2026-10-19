package update

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func discoverTPSBinaries(projectRoot string) []string {
	seen := map[string]bool{}
	var paths []string

	add := func(p string) {
		if p == "" {
			return
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			return
		}
		abs, err = filepath.EvalSymlinks(abs)
		if err != nil {
			abs = filepath.Clean(p)
		}
		if seen[abs] {
			return
		}
		if st, err := os.Stat(abs); err != nil || st.IsDir() {
			return
		}
		seen[abs] = true
		paths = append(paths, abs)
	}

	if out, err := exec.Command("which", "-a", "tps").Output(); err == nil {
		for _, line := range strings.Split(string(out), "\n") {
			add(strings.TrimSpace(line))
		}
	}

	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		add(filepath.Join(dir, "tps"))
	}

	add(filepath.Join(projectRoot, "tps"))

	return paths
}

func syncAllInstances(projectRoot, source string) ([]string, error) {
	source, err := filepath.Abs(source)
	if err != nil {
		return nil, err
	}
	source, err = filepath.EvalSymlinks(source)
	if err != nil {
		return nil, err
	}

	targets := discoverTPSBinaries(projectRoot)
	updated := make([]string, 0, len(targets))
	for _, dst := range targets {
		if sameFile(source, dst) {
			continue
		}
		if err := copyBinary(source, dst); err != nil {
			return updated, fmt.Errorf("copy to %s: %w", dst, err)
		}
		updated = append(updated, dst)
	}
	return updated, nil
}

func sameFile(a, b string) bool {
	sa, err1 := os.Stat(a)
	sb, err2 := os.Stat(b)
	if err1 != nil || err2 != nil {
		return filepath.Clean(a) == filepath.Clean(b)
	}
	return os.SameFile(sa, sb)
}

func copyBinary(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	mode := 0o755
	if st, err := os.Stat(src); err == nil {
		mode = int(st.Mode().Perm())
	}

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(mode))
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}
