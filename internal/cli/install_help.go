package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func installCopyRoots(cfg *config.Config, selected *config.Project, cwd string) (tpsRepo string, targets []linkTarget) {
	if selected != nil && canSelfUpdate(selected.Path) {
		tpsRepo = selected.Path
	}
	if cfg != nil {
		for _, p := range cfg.Projects {
			if canSelfUpdate(p.Path) {
				if tpsRepo == "" {
					tpsRepo = p.Path
				}
				break
			}
		}
	}
	targets = linkTargetsForCopy(cfg, selected, cwd)
	return tpsRepo, targets
}

func devRepoRoot(cfg *config.Config, selected *config.Project) string {
	tpsRepo, _ := installCopyRoots(cfg, selected, "")
	return tpsRepo
}

func installCommandsText(tpsRepo string, targets []linkTarget) string {
	cfgPath, err := config.Path()
	if err != nil {
		cfgPath = "~/Library/Application Support/tps/config.yaml"
	}

	cdTPS := "cd /path/to/tps"
	if tpsRepo != "" {
		cdTPS = "cd " + shellQuote(tpsRepo)
	}
	lines := []string{
		"echo '--- install / update tps ---'",
		cdTPS,
		"go install ./cmd/tps",
		"grep -q 'GOPATH/bin' ~/.zshrc 2>/dev/null || echo 'export PATH=\"$PATH:$(go env GOPATH)/bin\"' >> ~/.zshrc",
		"source ~/.zshrc 2>/dev/null || true",
		"which tps",
		"echo ''",
		"echo '--- link your app repo ---'",
	}
	if len(targets) == 0 {
		lines = append(lines, "echo 'Run ih -copy from your app directory, then paste or: tps setup -n name'")
	} else {
		for _, t := range targets {
			lines = append(lines, setupShellLine(t))
		}
	}
	lines = append(lines,
		"echo ''",
		"echo 'optional config path:'",
		"echo "+shellQuote(cfgPath),
	)
	return strings.Join(lines, "\n")
}

func setupShellLine(t linkTarget) string {
	base := filepath.Base(t.Path)
	line := "tps setup --path " + shellQuote(t.Path)
	if t.Name != "" && !strings.EqualFold(t.Name, base) {
		line = "tps setup -n " + shellQuote(t.Name) + " --path " + shellQuote(t.Path)
	}
	return line
}

func shellQuote(path string) string {
	if path == "" {
		return "''"
	}
	if !strings.ContainsAny(path, " \t\"'$\\") {
		return path
	}
	return "'" + strings.ReplaceAll(path, "'", "'\"'\"'") + "'"
}

func copyInstallCommands(cfg *config.Config, selected *config.Project) error {
	cwd, _ := os.Getwd()
	tpsRepo, targets := installCopyRoots(cfg, selected, cwd)
	text := installCommandsText(tpsRepo, targets)
	if err := ui.CopyToClipboard(text); err != nil {
		return err
	}
	fmt.Println(ui.Green("copied install commands to clipboard (⌘V in the other terminal)"))
	if len(targets) == 0 {
		fmt.Println(ui.Dim("tip: run ih -copy from inside your app repo (e.g. Clastone) to include project add"))
	} else {
		for _, t := range targets {
			fmt.Println(ui.Dim("  will link: "+t.Name+" at "+t.Path))
		}
	}
	return nil
}

func printInstallHelp() {
	cfgPath, err := config.Path()
	if err != nil {
		cfgPath = "~/Library/Application Support/tps/config.yaml (macOS)"
	}
	goBin := "$(go env GOPATH)/bin"
	if out, err := exec.Command("go", "env", "GOPATH").Output(); err == nil {
		goBin = strings.TrimSpace(string(out)) + "/bin"
	}

	fmt.Println(ui.Bold("Install tps on this Mac"))
	fmt.Println(ui.Dim("  clone or cd to the tps repo, then:"))
	fmt.Printf("    %s\n", ui.Cyan("go install ./cmd/tps"))
	fmt.Println(ui.Dim("  add Go’s bin dir to PATH (once), in ~/.zshrc:"))
	fmt.Printf("    %s\n", ui.Cyan("export PATH=\"$PATH:$(go env GOPATH)/bin\""))
	fmt.Println(ui.Dim("  reload shell and check:"))
	fmt.Printf("    %s  (expect %s)\n", ui.Cyan("which tps"), ui.Dim(goBin+"/tps"))

	fmt.Println()
	fmt.Println(ui.Bold("Use tps from any terminal"))
	fmt.Println(ui.Dim("  after PATH is set, run from anywhere:"))
	fmt.Printf("    %s\n", ui.Cyan("tps"))
	fmt.Println(ui.Dim("  link a repo (from inside the project, or with --path):"))
	fmt.Printf("    %s\n", ui.Cyan("tps setup"))
	fmt.Printf("    %s\n", ui.Cyan("tps setup -n myname"))
	fmt.Println(ui.Dim("  or:"))
	fmt.Printf("    %s\n", ui.Cyan("tps project add"))

	fmt.Println()
	fmt.Println(ui.Bold("Update the tps binary"))
	fmt.Println(ui.Dim("  inside the tps dev repo (interactive):"))
	fmt.Printf("    %s\n", ui.Cyan("u"))
	fmt.Println(ui.Dim("  sync every tps binary on PATH plus ./tps in this repo:"))
	fmt.Printf("    %s\n", ui.Cyan("u -a"))
	fmt.Println(ui.Dim("  copy paste-safe install commands (run ih -copy from your app repo):"))
	fmt.Printf("    %s\n", ui.Cyan("ih -copy"))
	fmt.Println(ui.Dim("  or from the repo root:"))
	fmt.Printf("    %s\n", ui.Cyan("go install ./cmd/tps"))
	fmt.Println(ui.Dim("  then start tps again if you did not use u"))

	fmt.Println()
	fmt.Println(ui.Bold("Another machine or fresh user account"))
	fmt.Println(ui.Dim("  1. Install Go, clone tps, go install, fix PATH (same as above)"))
	fmt.Println(ui.Dim("  2. Copy your config or re-link projects:"))
	fmt.Printf("    %s\n", ui.Dim(cfgPath))
	fmt.Println(ui.Dim("  3. Or copy that file from your other Mac, then run tps"))
	fmt.Println()
}
