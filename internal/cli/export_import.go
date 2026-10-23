package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/pick"
	"github.com/Cimmerial/tps/internal/ui"
)

func parseExportInput(fields []string) (MenuInput, error) {
	if len(fields) == 0 || !strings.EqualFold(fields[0], "export") {
		return MenuInput{}, fmt.Errorf("export: expected export")
	}
	in := MenuInput{Action: ActionExport}
	for _, f := range fields[1:] {
		switch f {
		case "-a", "--all":
			in.ExportAll = true
		default:
			in.Text = f
		}
	}
	return in, nil
}

func parseImportInput(fields []string) (MenuInput, error) {
	if len(fields) < 1 || !strings.EqualFold(fields[0], "import") {
		return MenuInput{}, fmt.Errorf("import: expected import")
	}
	if len(fields) < 2 {
		return MenuInput{Action: ActionImport}, nil
	}
	path := strings.TrimSpace(strings.Join(fields[1:], " "))
	if path == "" {
		return MenuInput{Action: ActionImport}, nil
	}
	return MenuInput{Action: ActionImport, Text: path}, nil
}

func runExport(cfg *config.Config, p *config.Project, in MenuInput) (backedOut bool, err error) {
	explicitPath := strings.TrimSpace(in.Text) != ""

	if in.ExportAll {
		err = runExportTarget(cfg, p, config.ExportKindEverything, "", in.Text)
		return false, err
	}
	if explicitPath {
		err = runExportTarget(cfg, p, config.ExportKindProject, p.Name, in.Text)
		return false, err
	}

	sections := exportMenuSections(cfg)
	for {
		ui.ClearScreen()
		printExportImportMenu("export", sections)
		key, back, err := readExportImportChoice(sections)
		if err != nil {
			return false, err
		}
		if back {
			return true, nil
		}
		switch key {
		case exportChoiceOptions:
			err = runExportTarget(cfg, p, config.ExportKindOptions, "", "")
			return false, err
		case exportChoiceEverythingK, exportChoiceEverything:
			err = runExportTarget(cfg, p, config.ExportKindEverything, "", "")
			return false, err
		default:
			proj := projectIDFromImportKey(cfg, key)
			if proj == nil {
				fmt.Println(ui.Dim("unknown project"))
				continue
			}
			err = runExportTarget(cfg, p, config.ExportKindProject, proj.Name, "")
			return false, err
		}
	}
}

func runExportTarget(cfg *config.Config, current *config.Project, kind config.ExportKind, projectName string, explicitPath string) error {
	prefs, err := config.LoadPreferences()
	if err != nil {
		return err
	}
	prefs, err = ensureExportUsername(prefs)
	if err != nil {
		return err
	}

	includeUsername := false
	if kind == config.ExportKindOptions || kind == config.ExportKindEverything {
		includeUsername, err = exportIncludeUsername(prefs.Username)
		if err != nil {
			return err
		}
	}

	var data []byte
	var label string
	switch kind {
	case config.ExportKindOptions:
		data, err = config.ExportPreferencesBundle(prefs, includeUsername)
		label = "options"
	case config.ExportKindEverything:
		data, err = config.ExportEverything(prefs, cfg.Projects, includeUsername)
		label = fmt.Sprintf("options + %d project(s)", len(cfg.Projects))
	default:
		proj := current
		if projectName != "" {
			for i := range cfg.Projects {
				if strings.EqualFold(cfg.Projects[i].Name, projectName) {
					proj = &cfg.Projects[i]
					break
				}
			}
		}
		if proj == nil {
			return fmt.Errorf("no project selected")
		}
		data, err = config.ExportProjectCommands(proj)
		projectName = proj.Name
		label = fmt.Sprintf("commands for %q", proj.Name)
	}
	if err != nil {
		return err
	}

	defaultName := config.ExportDefaultFilename(kind, projectName, prefs.Username)
	path := strings.TrimSpace(explicitPath)
	if path == "" {
		path, err = pick.ChooseSaveFile("Save tps export", defaultName)
		if err == pick.ErrCancelled {
			return nil
		}
		if err != nil {
			return err
		}
		path = pick.ResolveSavePath(path, defaultName)
	}
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = filepath.Join(cwd, defaultName)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		path = expandImportPath(path, cwd)
		path = pick.ResolveSavePath(path, defaultName)
	}
	if err := config.WriteExportFile(path, data); err != nil {
		return err
	}
	fmt.Println(ui.Green(fmt.Sprintf("exported %s to %s", label, path)))
	return nil
}

func ensureExportUsername(prefs *config.Preferences) (*config.Preferences, error) {
	if strings.TrimSpace(prefs.Username) != "" {
		return prefs, nil
	}
	fmt.Println()
	fmt.Printf("    %s\n", ui.Dim("username is required for export — saved in options"))
	line, err := ui.ReadEditedLine("username: ")
	if err != nil {
		return nil, err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("export cancelled: username required")
	}
	prefs.Username = line
	if err := config.SavePreferences(prefs); err != nil {
		return nil, err
	}
	fmt.Println(ui.Green("saved username in options"))
	return prefs, nil
}

func runImport(cfg *config.Config, path string, current *config.Project) (backedOut bool, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	if strings.TrimSpace(path) == "" {
		path, err = pick.ChooseFile("Choose tps export to import", pick.FileFilter{Extensions: []string{"yaml", "yml"}})
		if err == pick.ErrCancelled {
			return true, nil
		}
		if err != nil {
			return false, err
		}
	} else {
		path = expandImportPath(path, cwd)
	}

	bundle, err := config.ImportFromFile(path)
	if err != nil {
		return false, err
	}

	destKey, back, err := chooseImportDestination(cfg)
	if err != nil {
		return false, err
	}
	if back {
		return true, nil
	}

	if destKey == importChoiceOptions {
		if !bundle.HasPreferences() {
			fmt.Println(ui.Yellow("file has no options to import"))
			return false, nil
		}
		prefs, err := config.LoadPreferences()
		if err != nil {
			return false, err
		}
		includeUsername, err := importIncludeUsername(bundle.Preferences.Username)
		if err != nil {
			return false, err
		}
		session := prefs.SessionProjectID
		config.ApplyPortablePreferences(prefs, bundle.Preferences, includeUsername)
		prefs.SessionProjectID = session
		if err := config.SavePreferences(prefs); err != nil {
			return false, err
		}
		fmt.Println(ui.Green(fmt.Sprintf("imported options from %s", path)))
		return false, nil
	}

	dest := projectIDFromImportKey(cfg, destKey)
	if dest == nil {
		return false, fmt.Errorf("unknown import destination")
	}
	if !bundle.HasProjectCommands() {
		fmt.Println(ui.Yellow("file has no project commands to import"))
		return false, nil
	}

	cmds, ok := bundle.CommandsForImport(dest.Name)
	if !ok {
		entries := bundle.BundleProjectEntries()
		name, back, err := pickBundleProjectKey(entries)
		if err != nil {
			return false, err
		}
		if back {
			return true, nil
		}
		for _, e := range entries {
			if strings.EqualFold(e.Name, name) {
				cmds = e.Commands
				ok = true
				break
			}
		}
	}
	if !ok || len(cmds) == 0 {
		return false, fmt.Errorf("no commands found in file for this project")
	}

	for i := range cfg.Projects {
		if cfg.Projects[i].ID == dest.ID {
			config.MergeCommandsByAlias(&cfg.Projects[i], cmds)
			if err := config.Save(cfg); err != nil {
				return false, err
			}
			fmt.Println(ui.Green(fmt.Sprintf("merged %d command(s) into %q from %s", len(cmds), dest.Name, path)))
			return false, nil
		}
	}
	return false, fmt.Errorf("project %q not found", dest.Name)
}

func chooseImportDestination(cfg *config.Config) (string, bool, error) {
	sections := importMenuSections(cfg)
	title := "import"
	for {
		ui.ClearScreen()
		printExportImportMenu(title, sections)
		key, back, err := readExportImportChoice(sections)
		if err != nil || back {
			return "", back, err
		}
		if key == importChoiceOptions || projectIDFromImportKey(cfg, key) != nil {
			return key, false, nil
		}
		fmt.Println(ui.Dim("unknown choice"))
	}
}

func expandImportPath(path, cwd string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return path
	}
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
		}
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(cwd, path)
	}
	return filepath.Clean(path)
}
