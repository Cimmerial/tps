package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const releaseRepo = "Cimmerial/tps"

type ghRelease struct {
	TagName string `json:"tag_name"`
}

// InstallFromGitHubRelease downloads the latest (or current) release binary for this OS/arch.
func InstallFromGitHubRelease() (InstallResult, error) {
	tag, err := latestReleaseTag()
	if err != nil {
		return InstallResult{}, err
	}
	ver := strings.TrimPrefix(tag, "v")
	asset := releaseAssetName(ver)
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", releaseRepo, tag, asset)

	tmp, err := os.MkdirTemp("", "tps-update-*")
	if err != nil {
		return InstallResult{}, err
	}
	defer os.RemoveAll(tmp)

	archive := filepath.Join(tmp, asset)
	if err := downloadFile(url, archive); err != nil {
		return InstallResult{}, err
	}

	bin := filepath.Join(tmp, "tps")
	if runtime.GOOS == "windows" {
		bin = filepath.Join(tmp, "tps.exe")
	}
	if err := extractRelease(archive, tmp, runtime.GOOS); err != nil {
		return InstallResult{}, err
	}

	clearWindowsReplaceQueue()
	dest, err := exec.LookPath("tps")
	if err != nil {
		dest, err = os.Executable()
		if err != nil {
			return InstallResult{}, fmt.Errorf("locate tps binary: %w", err)
		}
	}
	if runtime.GOOS == "windows" {
		if exe, err := os.Executable(); err == nil {
			dest = exe
		}
	}
	dest, err = filepath.EvalSymlinks(dest)
	if err != nil {
		return InstallResult{}, err
	}

	before, _ := fileFingerprint(dest)
	staged, err := installBinary(bin, dest)
	if err != nil {
		return InstallResult{}, err
	}
	after, _ := fileFingerprint(dest)
	updated := staged || before != after
	if before == "" && after != "" {
		updated = true
	}
	return InstallResult{Updated: updated}, nil
}

func latestReleaseTag() (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", releaseRepo))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github releases: HTTP %d", resp.StatusCode)
	}
	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", err
	}
	if rel.TagName == "" {
		return "", fmt.Errorf("empty release tag")
	}
	return rel.TagName, nil
}

func releaseAssetName(version string) string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	if goarch == "386" {
		goarch = "amd64"
	}
	base := fmt.Sprintf("tps_%s_%s_%s", version, goos, goarch)
	if goos == "windows" {
		return base + ".zip"
	}
	return base + ".tar.gz"
}

// ReleaseAssetName is exported for tests.
func ReleaseAssetName(version string) string {
	return releaseAssetName(version)
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download %s: HTTP %d", url, resp.StatusCode)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func extractRelease(archive, destDir, goos string) error {
	if goos == "windows" {
		return extractZip(archive, destDir)
	}
	return extractTarGz(archive, destDir)
}

func extractTarGz(archive, destDir string) error {
	f, err := os.Open(archive)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}
		out := filepath.Join(destDir, filepath.Base(h.Name))
		w, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, tr); err != nil {
			w.Close()
			return err
		}
		w.Close()
	}
	return nil
}

func extractZip(archive, destDir string) error {
	r, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out := filepath.Join(destDir, filepath.Base(f.Name))
		w, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(w, rc)
		w.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
