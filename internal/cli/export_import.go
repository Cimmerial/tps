package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
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
	if len(fields) < 2 || !strings.EqualFold(fields[0], "import") {
		return MenuInput{}, fmt.Errorf("import: expected import <file>")
	}
	path := strings.TrimSpace(strings.Join(fields[1:], " "))
	if path == "" {
		return MenuInput{}, fmt.Errorf("import: expected import <file>")
	}
	return MenuInput{Action: ActionImport, Text: path}, nil
}

func runExport(cfg *config.Config, p *config.Project, in MenuInput) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if in.ExportAll {
		data, err := config.ExportAllProjects(cfg.Projects)
		if err != nil {
			return err
		}
		path := filepath.Join(cwd, "tps-export.yaml")
		if strings.TrimSpace(in.Text) != "" {
			path = expandImportPath(in.Text, cwd)
		}
		if err := config.WriteExportFile(path, data); err != nil {
			return err
		}
		fmt.Println(ui.Green(fmt.Sprintf("exported %d project(s) to %s", len(cfg.Projects), path)))
		return nil
	}
	data, err := config.ExportProject(p)
	if err != nil {
		return err
	}
	path := filepath.Join(cwd, config.SanitizeExportBasename(p.Name)+".tps.yaml")
	if strings.TrimSpace(in.Text) != "" {
		path = expandImportPath(in.Text, cwd)
	}
	if err := config.WriteExportFile(path, data); err != nil {
		return err
	}
	fmt.Println(ui.Green(fmt.Sprintf("exported %q to %s", p.Name, path)))
	return nil
}

func runImport(cfg *config.Config, path string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	path = expandImportPath(path, cwd)
	bundle, err := config.ImportFromFile(path)
	if err != nil {
		return err
	}
	n, err := cfg.MergeImport(bundle)
	if err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Println(ui.Green(fmt.Sprintf("imported %d project(s) from %s", n, path)))
	return nil
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
