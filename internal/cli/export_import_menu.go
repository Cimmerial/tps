package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

const (
	exportChoiceOptions     = "options"
	exportChoiceEverything  = "everything"
	exportChoiceEverythingK = "all"
	importChoiceOptions     = "options"
)

type exportImportItem struct {
	Num  int
	Name string
	Key  string
	Desc string
}

type exportImportSection struct {
	title string
	items []exportImportItem
}

func printExportImportMenu(title string, sections []exportImportSection) {
	fmt.Println()
	fmt.Printf("%s\n", ui.Header(title))
	maxNum, maxName := 0, 0
	for _, sec := range sections {
		for _, it := range sec.items {
			numLen := len(strconv.Itoa(it.Num))
			if numLen > maxNum {
				maxNum = numLen
			}
			if len(it.Name) > maxName {
				maxName = len(it.Name)
			}
		}
	}
	for i, sec := range sections {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s\n", ui.SubHeader(sec.title))
		for _, it := range sec.items {
			numCol := ui.Bold(padPlain(strconv.Itoa(it.Num), maxNum))
			nameCol := ui.Bold(padPlain(it.Name, maxName))
			fmt.Printf("    %s | %s | %s\n", numCol, nameCol, ui.Dim(it.Desc))
		}
	}
	fmt.Println()
	fmt.Printf("    %s\n", ui.Dim("number or name — Enter or 0 or q = back"))
	fmt.Println()
}

func isExportImportBack(line string) bool {
	line = strings.TrimSpace(line)
	if line == "" {
		return true
	}
	switch strings.ToLower(line) {
	case "0", "q", "quit", "exit":
		return true
	default:
		return false
	}
}

func readExportImportChoice(sections []exportImportSection) (key string, back bool, err error) {
	var items []exportImportItem
	for _, sec := range sections {
		items = append(items, sec.items...)
	}
	for {
		line, err := ui.ReadEditedLine(ui.Prompt())
		if err != nil {
			return "", false, err
		}
		if isExportImportBack(line) {
			return "", true, nil
		}
		line = strings.TrimSpace(line)
		low := strings.ToLower(line)
		if n, err := strconv.Atoi(line); err == nil {
			for _, it := range items {
				if it.Num == n {
					return it.Key, false, nil
				}
			}
			fmt.Println(ui.Dim("unknown choice"))
			continue
		}
		for _, it := range items {
			if strings.EqualFold(line, it.Name) || strings.EqualFold(line, it.Key) {
				return it.Key, false, nil
			}
		}
		if low == exportChoiceEverything {
			return exportChoiceEverythingK, false, nil
		}
		fmt.Println(ui.Dim("unknown choice"))
	}
}

func exportMenuSections(cfg *config.Config) []exportImportSection {
	n := 1
	var sections []exportImportSection
	sections = append(sections, exportImportSection{
		title: "preferences",
		items: []exportImportItem{{
			Num:  n,
			Name: "options",
			Key:  exportChoiceOptions,
			Desc: "export options (logo, colors, username, …)",
		}},
	})
	n++
	var projectItems []exportImportItem
	for _, p := range cfg.Projects {
		projectItems = append(projectItems, exportImportItem{
			Num:  n,
			Name: p.Name,
			Key:  p.ID,
			Desc: "export this project's commands",
		})
		n++
	}
	if len(projectItems) > 0 {
		sections = append(sections, exportImportSection{title: "projects", items: projectItems})
	}
	sections = append(sections, exportImportSection{
		title: "all",
		items: []exportImportItem{{
			Num:  n,
			Name: "everything",
			Key:  exportChoiceEverythingK,
			Desc: "export options and every project's commands (one file)",
		}},
	})
	return sections
}

func importMenuSections(cfg *config.Config) []exportImportSection {
	n := 1
	var sections []exportImportSection
	sections = append(sections, exportImportSection{
		title: "preferences",
		items: []exportImportItem{{
			Num:  n,
			Name: "options",
			Key:  importChoiceOptions,
			Desc: "import options from file",
		}},
	})
	n++
	var projectItems []exportImportItem
	for _, p := range cfg.Projects {
		projectItems = append(projectItems, exportImportItem{
			Num:  n,
			Name: p.Name,
			Key:  p.ID,
			Desc: "merge file commands into this project",
		})
		n++
	}
	if len(projectItems) > 0 {
		sections = append(sections, exportImportSection{title: "projects", items: projectItems})
	}
	return sections
}

func projectIDFromImportKey(cfg *config.Config, key string) *config.Project {
	if cfg == nil {
		return nil
	}
	for i := range cfg.Projects {
		p := &cfg.Projects[i]
		if strings.EqualFold(p.ID, key) || strings.EqualFold(p.Name, key) {
			return p
		}
	}
	return nil
}

func bundleProjectMenuSections(entries []config.Project) []exportImportSection {
	n := 1
	var items []exportImportItem
	for _, p := range entries {
		key := bundleProjectExportKey(p)
		items = append(items, exportImportItem{
			Num:  n,
			Name: p.Name,
			Key:  key,
			Desc: fmt.Sprintf("%d command(s)", len(p.Commands)),
		})
		n++
	}
	return []exportImportSection{{title: "projects in file", items: items}}
}

func bundleProjectExportKey(p config.Project) string {
	if strings.TrimSpace(p.ID) != "" {
		return p.ID
	}
	return strings.ToLower(config.SanitizeExportBasename(p.Name))
}

func pickBundleProjectKey(entries []config.Project) (string, bool, error) {
	sections := bundleProjectMenuSections(entries)
	for {
		ui.ClearScreen()
		printExportImportMenu("import — pick project from file", sections)
		key, back, err := readExportImportChoice(sections)
		if err != nil || back {
			return "", back, err
		}
		for _, p := range entries {
			if strings.EqualFold(bundleProjectExportKey(p), key) || strings.EqualFold(p.Name, key) {
				return p.Name, false, nil
			}
		}
		fmt.Println(ui.Dim("unknown choice"))
	}
}
