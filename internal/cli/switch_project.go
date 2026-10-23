package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/ui"
)

func switchProject(reader *bufio.Reader, cfg *config.Config, currentID string, name string) (string, error) {
	if len(cfg.Projects) == 0 {
		return "", fmt.Errorf("no projects linked")
	}
	if name != "" {
		id := findProjectID(cfg, name)
		if id == "" {
			return "", fmt.Errorf("no project matching %q", name)
		}
		return id, nil
	}

	ui.ClearScreen()
	fmt.Println()
	fmt.Printf("%s\n", ui.HeaderInfo("switch project"))
	fmt.Println()

	rows := make([]numberedMenuRow, len(cfg.Projects))
	for i, p := range cfg.Projects {
		row := numberedMenuRow{
			Num:      fmt.Sprintf("%d", i+1),
			Desc:     p.Name,
			HelpLine: p.Path,
		}
		if p.ID == currentID {
			row.Status = ui.Green("current")
		}
		rows[i] = row
	}
	printChoiceEditMenuPanel("", []menuPanelSection{{Rows: rows}}, nil)
	fmt.Println()
	fmt.Println(ui.Dim("number or name — Enter or q = back"))
	fmt.Println()

	line, err := ui.ReadMenuLine()
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if isProjectPickBack(line) {
		return "", nil
	}
	if id := findProjectID(cfg, line); id != "" {
		return id, nil
	}
	n, err := parsePickIndex(line, len(cfg.Projects))
	if err != nil {
		return "", err
	}
	return cfg.Projects[n-1].ID, nil
}

// isProjectPickBack reports whether a picker line backs out of the switch-project
// menu instead of switching: Enter on its own (or q), like the other tps submenus.
func isProjectPickBack(line string) bool {
	line = strings.TrimSpace(line)
	return line == "" || strings.EqualFold(line, "q")
}

func findProjectID(cfg *config.Config, query string) string {
	query = strings.TrimSpace(query)
	for i := range cfg.Projects {
		p := &cfg.Projects[i]
		if strings.EqualFold(p.Name, query) || strings.EqualFold(p.ID, query) {
			return p.ID
		}
	}
	lower := strings.ToLower(query)
	for i := range cfg.Projects {
		p := &cfg.Projects[i]
		if strings.Contains(strings.ToLower(p.Name), lower) {
			return p.ID
		}
	}
	return ""
}

// applyInteractiveProjectSwitch moves the session to newID when it differs from projectID.
// It returns the updated project and back-target ids.
func applyInteractiveProjectSwitch(selected *config.Project, projectID, prevProjectID, newID string, statusTick int) (string, string, error) {
	if newID == "" || newID == projectID {
		return projectID, prevProjectID, nil
	}
	prevName := selected.Name
	backTarget := projectID
	if err := config.SaveSessionProjectID(newID); err != nil {
		return projectID, prevProjectID, err
	}
	if err := refreshInteractiveSession(newID, statusTick, prevName, false); err != nil {
		return projectID, prevProjectID, err
	}
	return newID, backTarget, nil
}

func parsePickIndex(line string, max int) (int, error) {
	var n int
	_, err := fmt.Sscanf(line, "%d", &n)
	if err != nil || n < 1 || n > max {
		return 0, fmt.Errorf("invalid selection")
	}
	return n, nil
}
