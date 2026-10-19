package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Cimmerial/tps/internal/addon"
	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/install"
	"github.com/Cimmerial/tps/internal/instance"
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
	Long:  "Install or update tps from your linked dev repo when present, ensure Go bin is on PATH in ~/.zshrc, link the project, and refresh TPS.md.",
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

	root, err := project.ProjectRootFromCwd(linkDir)
	if err != nil {
		return err
	}
	if existing := cfg.FindByPath(root); existing != nil && nameOverride == "" && pathOverride == "" {
		fmt.Println(ui.Yellow(fmt.Sprintf("%q is already set up at %s", existing.Name, existing.Path)))
		fmt.Println(ui.Dim("add commands in config.yaml, nc, or export/import · run ih for install help"))
		if tpsRepo := instance.DevRepoRoot(cfg); tpsRepo != "" {
			if _, err := instance.SyncProject(tpsRepo, existing.Path); err != nil {
				return err
			}
		}
		return nil
	}

	tpsRepo := devRepoRoot(cfg, nil)
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

	if err := install.EnsureOnPath(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", ui.Dim("PATH: "+err.Error()))
	} else {
		fmt.Println(ui.Dim("PATH ok"))
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
	if err := ensureDevProjectCommands(cfg, linked); err != nil {
		return linked, err
	}
	tpsRepo := instance.DevRepoRoot(cfg)
	if _, err := instance.SyncProject(tpsRepo, linked.Path); err != nil {
		return linked, err
	}
	if m, err := addon.LoadManifest(linked.Path); err == nil {
		if err := addon.MergeManifestCommands(linked); err != nil {
			return linked, err
		}
		addon.PrintContractWarnings(m)
		_ = config.Save(cfg)
		_ = buildUtilitiesBinary(linked.Path)
	}
	return linked, nil
}

