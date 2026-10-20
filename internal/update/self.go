package update

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/instance"
	"github.com/Cimmerial/tps/internal/ui"
	"github.com/Cimmerial/tps/internal/version"
)

type InstallResult struct {
	Updated bool
	Version string // release tag (vX.Y.Z) when installing from GitHub
}

func InstallFromProjectRoot(root string) (InstallResult, error) {
	return installFromProjectRoot(root, false)
}

func installFromProjectRoot(root string, quiet bool) (InstallResult, error) {
	mainDir := filepath.Join(root, "cmd", "tps")
	if _, err := os.Stat(mainDir); err != nil {
		return InstallResult{}, fmt.Errorf("update needs a Go CLI at cmd/tps (run from the tps dev repo)")
	}

	binPath, _ := exec.LookPath("tps")
	before, _ := fileFingerprint(binPath)

	cmd := exec.Command("go", "install", "./cmd/tps")
	cmd.Dir = root
	if quiet {
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			msg := strings.TrimSpace(stderr.String())
			if msg != "" {
				return InstallResult{}, fmt.Errorf("%w: %s", err, msg)
			}
			return InstallResult{}, err
		}
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return InstallResult{}, err
		}
	}

	afterPath, err := exec.LookPath("tps")
	if err != nil {
		afterPath = binPath
	}
	after, _ := fileFingerprint(afterPath)

	updated := before != after
	if before == "" && after != "" {
		updated = true
	}
	return InstallResult{Updated: updated}, nil
}

// MaybeAutoUpdate rebuilds tps from root when the installed binary changed. Restarts on success.
func MaybeAutoUpdate(root string) error {
	result, err := installFromProjectRoot(root, true)
	if err != nil {
		return err
	}
	if !result.Updated {
		return nil
	}
	syncInstanceReadmes(root)
	if err := RestartInteractive(true, version.Label()); err != nil {
		return fmt.Errorf("restart failed: %w (run tps again)", err)
	}
	return nil
}

func fileFingerprint(path string) (string, error) {
	if path == "" {
		return "", nil
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func InstallAndRestart(root string) error {
	var result InstallResult
	var err error
	if root == "" {
		result, err = InstallFromGitHubRelease()
		if err != nil {
			return err
		}
		source, err := exec.LookPath("tps")
		if err != nil {
			source, _ = os.Executable()
		}
		synced, err := syncAllInstances("", source)
		if err != nil {
			return err
		}
		if len(synced) > 0 {
			fmt.Println(ui.Dim(fmt.Sprintf("synced %d other install(s)", len(synced))))
		}
		updated := result.Updated || len(synced) > 0
		syncInstanceReadmes("")
		return finishInstallAndRestart(updated, result.Version)
	}

	if err := PullFastForward(root); err != nil {
		return err
	}
	result, err = InstallFromProjectRoot(root)
	if err != nil {
		return err
	}

	source, err := exec.LookPath("tps")
	if err != nil {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("locate installed tps: %w", err)
		}
		source = exe
	}

	synced, err := syncAllInstances(root, source)
	if err != nil {
		return err
	}
	if len(synced) > 0 {
		fmt.Println(ui.Dim(fmt.Sprintf("synced %d other install(s)", len(synced))))
	}

	updated := result.Updated || len(synced) > 0
	syncInstanceReadmes(root)
	ver := result.Version
	if ver == "" {
		ver = version.Label()
	}
	return finishInstallAndRestart(updated, ver)
}

func finishInstallAndRestart(updated bool, installedVersion string) error {
	if runtime.GOOS == "windows" {
		if !updated {
			fmt.Println(ui.Yellow("tps already up to date") + ui.Dim(" ("+version.Label()+")"))
			return nil
		}
		if err := RestartInteractive(true, installedVersion); err != nil {
			return err
		}
		return nil
	}
	if updated {
		fmt.Println(ui.Green("tps updated — restarting…"))
	}
	if err := RestartInteractive(updated, installedVersion); err != nil {
		return fmt.Errorf("restart failed: %w (run tps again)", err)
	}
	return nil
}

func syncInstanceReadmes(tpsRepo string) {
	cfg, err := config.Load()
	if err != nil {
		return
	}
	if _, err := instance.SyncAllLinked(tpsRepo, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("instance readme: "+err.Error()))
	}
}
