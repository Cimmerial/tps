package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/project"
	"github.com/Cimmerial/tps/internal/ui"
	"github.com/Cimmerial/tps/internal/update"
	"github.com/spf13/cobra"
)

var setupName string
var setupPath string

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install/update tps and link a project (run from the app repo)",
	Long:  "Runs the same steps as the ih -copy install script: go install from your linked tps dev repo, ensure Go bin is on PATH in ~/.zshrc, then link the project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSetup(setupName, setupPath)
	},
}

func init() {
	setupCmd.Flags().StringVarP(&setupName, "name", "n", "", "project name (default: directory name)")
	setupCmd.Flags().StringVar(&setupPath, "path", "", "project directory to link (default: current directory)")
	rootCmd.AddCommand(setupCmd)
}

func runSetup(nameOverride, pathOverride string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	linkDir := cwd
	if pathOverride != "" {
		linkDir = pathOverride
	}

	tpsRepo, _ := installCopyRoots(cfg, nil, cwd)
	if tpsRepo != "" {
		spin := ui.NewSpinner("installing tps")
		result, err := update.InstallFromProjectRoot(tpsRepo)
		spin.Stop()
		if err != nil {
			return err
		}
		if result.Updated {
			fmt.Println(ui.Green("tps updated"))
		} else {
			fmt.Println(ui.Dim("tps already up to date"))
		}
	} else {
		fmt.Println(ui.Dim("skip install: link the tps dev repo in config first (tps project add from the tps repo)"))
	}

	if err := ensurePathInZshrc(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("PATH: "+err.Error()))
	} else {
		fmt.Println(ui.Dim("PATH ok (GOPATH/bin in ~/.zshrc)"))
	}

	if _, err := exec.LookPath("tps"); err != nil {
		fmt.Println(ui.Dim("reload your shell: source ~/.zshrc  (then: which tps)"))
	}

	linked, err := linkProjectAt(cfg, linkDir, nameOverride)
	if err != nil {
		return err
	}
	fmt.Printf("%s %q at %s\n", ui.Green("linked"), linked.Name, linked.Path)
	if linked.GitRemote != "" {
		fmt.Printf("%s %s\n", ui.Dim("remote:"), linked.GitRemote)
	}
	fmt.Println(ui.Dim("run: tps"))
	return nil
}

func linkProjectAt(cfg *config.Config, linkDir, name string) (*config.Project, error) {
	if name == "" {
		root, err := project.ProjectRootFromCwd(linkDir)
		if err != nil {
			return nil, err
		}
		name = filepath.Base(root)
	}
	linked, err := project.LinkFromCwd(cfg, linkDir, name)
	if err != nil {
		return nil, err
	}
	if err := config.Save(cfg); err != nil {
		return nil, err
	}
	return linked, nil
}

func ensurePathInZshrc() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	zshrc := filepath.Join(home, ".zshrc")
	var content []byte
	if data, err := os.ReadFile(zshrc); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		content = data
	}
	if strings.Contains(string(content), "GOPATH/bin") {
		return nil
	}
	line := "export PATH=\"$PATH:$(go env GOPATH)/bin\"\n"
	f, err := os.OpenFile(zshrc, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(line); err != nil {
		return err
	}
	return nil
}

func execLookPathTps() (string, error) {
	return os.Getenv("PATH"), nil // weak check - use exec.LookPath
}
