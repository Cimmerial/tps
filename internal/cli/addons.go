package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps/internal/addon"
	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/pick"
	"github.com/Cimmerial/tps/internal/ui"
)

func syncAddonManifest(cfg *config.Config, p *config.Project) error {
	if cfg == nil || p == nil {
		return nil
	}
	if err := addon.SyncManifestCommands(p); err != nil {
		return err
	}
	return config.Save(cfg)
}

func printAddonsHelpSection() {
	printInfoMenuPanel("TPS ADDONS", []menuPanelSection{{
		Label: "install",
		Rows: []numberedMenuRow{{
			Desc:     "install utilities",
			Abbrev:   "iu",
			HelpLine: "clone and link tps-utilities (folder picker on macOS / Windows)",
		}},
	}})
}

func linkedUtilitiesProjects(cfg *config.Config) []*config.Project {
	if cfg == nil {
		return nil
	}
	var out []*config.Project
	for i := range cfg.Projects {
		p := &cfg.Projects[i]
		if addon.IsUtilitiesProject(p) {
			out = append(out, p)
		}
	}
	return out
}

func confirmAnotherUtilitiesInstall(cfg *config.Config) (bool, error) {
	linked := linkedUtilitiesProjects(cfg)
	if len(linked) == 0 {
		return true, nil
	}
	fmt.Println()
	fmt.Println(ui.Yellow("tps-utilities is already installed:"))
	for _, p := range linked {
		fmt.Printf("    %s  %s\n", ui.Bold(p.Name), ui.Dim(p.Path))
	}
	fmt.Println()
	fmt.Println(ui.Dim("Install another instance? [y/N]"))
	ui.SetPromptChrome(ui.ChromeChoice)
	line, err := ui.ReadEditedLine(ui.Prompt())
	if err != nil {
		return false, err
	}
	if !parseConfirmYes(line) {
		fmt.Println(ui.Dim("cancelled"))
		return false, nil
	}
	return true, nil
}

func parseConfirmYes(line string) bool {
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes"
}

func runInstallUtilities(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("no config")
	}
	ok, err := confirmAnotherUtilitiesInstall(cfg)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	parent, err := pick.ChooseFolder("Choose where to install tps-utilities")
	if err != nil {
		if err == pick.ErrCancelled {
			fmt.Println(ui.Dim("cancelled"))
			return nil
		}
		return err
	}
	dest := filepath.Join(parent, addon.UtilitiesDirName)
	if _, err := os.Stat(dest); err == nil {
		if _, err := addon.LoadManifest(dest); err != nil {
			return fmt.Errorf("%s exists but is not a tps addon", dest)
		}
		fmt.Println(ui.Dim("already cloned at " + dest))
	} else {
		spin := ui.NewSpinner("cloning tps-utilities")
		cmd := exec.Command("git", "clone", addon.UtilitiesCloneURL, dest)
		out, err := cmd.CombinedOutput()
		spin.Stop()
		if err != nil {
			return fmt.Errorf("git clone: %v: %s", err, strings.TrimSpace(string(out)))
		}
		fmt.Println(ui.Green("cloned to " + dest))
	}

	if err := buildUtilitiesBinary(dest); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", ui.Yellow("build tpsu: "+err.Error()))
	}

	linked, err := linkProjectAt(cfg, dest, addon.UtilitiesDirName)
	if err != nil {
		return err
	}
	if err := addon.MergeManifestCommands(linked); err != nil {
		return err
	}
	if m, err := addon.LoadManifest(dest); err == nil {
		addon.PrintContractWarnings(m)
	}
	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Printf("%s %q — switch with %s then run %s\n",
		ui.Green("utilities ready"),
		linked.Name,
		ui.Abbrev("p"),
		ui.Abbrev("uic"),
	)
	return nil
}

func buildUtilitiesBinary(utilitiesRoot string) error {
	binDir := filepath.Join(utilitiesRoot, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return err
	}
	out := filepath.Join(binDir, "tpsu")
	cmd := exec.Command("go", "build", "-o", out, "./cmd/tpsu")
	cmd.Dir = utilitiesRoot
	combined, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, strings.TrimSpace(string(combined)))
	}
	return nil
}

func parseInstallUtilitiesInput(fields []string) (MenuInput, bool, error) {
	if len(fields) == 0 {
		return MenuInput{}, false, nil
	}
	if len(fields) == 1 && strings.EqualFold(fields[0], "iu") {
		return MenuInput{Action: ActionInstallUtilities}, true, nil
	}
	if len(fields) >= 2 && strings.EqualFold(fields[0], "install") && strings.EqualFold(fields[1], "utilities") {
		return MenuInput{Action: ActionInstallUtilities}, true, nil
	}
	return MenuInput{}, false, nil
}
