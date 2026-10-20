package addonmenu

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/pick"
	"github.com/Cimmerial/tps/internal/runner"
	"github.com/Cimmerial/tps/internal/ui"
	"gopkg.in/yaml.v3"
)

// RunCommandMenu runs a project command's menu tree until the user backs out.
func RunCommandMenu(p *config.Project, cmd *config.Command) (int, error) {
	if p == nil || cmd == nil {
		return 1, fmt.Errorf("no command")
	}
	menuRel := strings.TrimSpace(cmd.Menu)
	if menuRel == "" {
		return 1, fmt.Errorf("command %q has no menu", cmd.Alias)
	}
	path := filepath.Join(p.Path, menuRel)
	data, err := os.ReadFile(path)
	if err != nil {
		return 1, err
	}
	var root MenuFile
	if err := yaml.Unmarshal(data, &root); err != nil {
		return 1, err
	}
	stack := []*MenuFile{&root}
	menuFresh := true
	for {
		cur := stack[len(stack)-1]
		title := strings.TrimSpace(cur.Title)
		if title == "" {
			title = cmd.Alias
		}
		if menuFresh {
			ui.ClearScreen()
			menuFresh = false
		} else {
			fmt.Println()
		}
		fmt.Println()
		fmt.Printf("%s\n", ui.Header(title))
		fmt.Println()
		fmt.Printf("    %s\n", ui.Dim("number to choose — Enter, 0, or q = back"))
		fmt.Println()

		choices := visibleItems(cur.Items)
		for i, item := range choices {
			num := fmt.Sprintf("%d", i+1)
			label := item.Label
			if IsDimmed(item) {
				label = ui.Dim(label)
			} else {
				label = ui.Bold(label)
			}
			fmt.Printf("    %s  %s\n", ui.Abbrev(num), label)
		}
		fmt.Println()
		fmt.Printf("    %s  %s\n", ui.Abbrev("0"), ui.Dim("back"))
		fmt.Println()

		choice, back, err := readAddonMenuChoice()
		if err != nil {
			return 1, err
		}
		if back {
			if len(stack) == 1 {
				ui.ClearScreen()
				return 0, nil
			}
			stack = stack[:len(stack)-1]
			continue
		}
		if choice < 1 || choice > len(choices) {
			fmt.Printf("%s\n", ui.Yellow(fmt.Sprintf("enter a number 1–%d or 0 to go back", len(choices))))
			continue
		}
		item := choices[choice-1]
		if item.Back {
			if len(stack) == 1 {
				ui.ClearScreen()
				return 0, nil
			}
			stack = stack[:len(stack)-1]
			continue
		}
		if reason := DimReason(item); reason != "" {
			fmt.Printf("%s\n", ui.Yellow(reason))
			continue
		}
		if len(item.Children) > 0 {
			child := &MenuFile{Title: item.Label, Items: item.Children}
			stack = append(stack, child)
			continue
		}
		if item.Action != nil {
			code, _, err := runAction(p, item.Action)
			if err != nil {
				return code, err
			}
			if code != 0 {
				fmt.Fprintf(os.Stderr, "command exited with code %d\n", code)
			}
			continue
		}
		fmt.Printf("%s\n", ui.Yellow("menu item has no action"))
	}
}

// readAddonMenuChoice matches edit/options menus: Enter or 0 = back.
func readAddonMenuChoice() (choice int, back bool, err error) {
	line, err := ui.ReadEditedLine(ui.Prompt())
	if err != nil {
		return 0, false, err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, true, nil
	}
	low := strings.ToLower(line)
	if low == "0" || low == "q" || low == "quit" || low == "back" {
		return 0, true, nil
	}
	n, err := strconv.Atoi(line)
	if err != nil {
		fmt.Println(ui.Dim("enter a number"))
		return 0, true, nil
	}
	return n, false, nil
}

func visibleItems(items []MenuItem) []MenuItem {
	return items
}

func parseChoice(line string, max int) (int, error) {
	var choice int
	if _, err := fmt.Sscanf(line, "%d", &choice); err != nil {
		return 0, fmt.Errorf("enter a number 1–%d or 0 to go back", max)
	}
	if choice < 1 || choice > max {
		return 0, fmt.Errorf("enter a number 1–%d or 0 to go back", max)
	}
	return choice, nil
}

func printActionSummary(lines []string) {
	if len(lines) == 0 {
		return
	}
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println()
}

func runAction(p *config.Project, action *MenuAction) (int, []string, error) {
	picks, err := collectPicks(action.Picks)
	if err != nil {
		if err == pick.ErrCancelled {
			fmt.Println(ui.Dim("picker cancelled"))
			return 0, nil, nil
		}
		fmt.Printf("%s\n", ui.Yellow(err.Error()))
		return 0, nil, nil
	}
	shell := expandShell(p.Path, action.Shell, picks)
	if strings.Contains(action.Shell, "bin/tpsu") {
		bin := filepath.Join(p.Path, "bin", "tpsu")
		if _, err := os.Stat(bin); err != nil {
			fmt.Printf("%s\n", ui.Yellow("missing bin/tpsu — in the addon repo run: go build -o bin/tpsu ./cmd/tpsu"))
			return 1, nil, nil
		}
	}
	code, err := runner.RunScript(p, shell)
	return code, nil, err
}

func collectPicks(specs []PickSpec) ([]string, error) {
	out := make([]string, 0, len(specs))
	for _, spec := range specs {
		prompt := strings.TrimSpace(spec.Prompt)
		if prompt == "" {
			prompt = "Choose"
		}
		switch strings.ToLower(strings.TrimSpace(spec.Kind)) {
		case "file":
			path, err := pick.ChooseFile(prompt)
			if err != nil {
				return nil, err
			}
			out = append(out, path)
		case "files":
			paths, err := pick.ChooseFiles(prompt)
			if err != nil {
				return nil, err
			}
			out = append(out, strings.Join(paths, "\n"))
		case "folder":
			path, err := pick.ChooseFolder(prompt)
			if err != nil {
				if spec.Optional {
					out = append(out, "")
					continue
				}
				return nil, err
			}
			out = append(out, path)
		default:
			return nil, fmt.Errorf("unknown pick kind %q", spec.Kind)
		}
	}
	return out, nil
}

func expandShell(projectRoot, template string, picks []string) string {
	s := strings.ReplaceAll(template, "{root}", shellQuote(projectRoot))
	for i, v := range picks {
		placeholder := fmt.Sprintf("[PICK%d]", i+1)
		if v == "" {
			s = strings.ReplaceAll(s, placeholder, "")
			continue
		}
		if strings.Contains(v, "\n") {
			var quoted []string
			for _, part := range strings.Split(v, "\n") {
				part = strings.TrimSpace(part)
				if part != "" {
					quoted = append(quoted, shellQuote(part))
				}
			}
			s = strings.ReplaceAll(s, placeholder, strings.Join(quoted, " "))
			continue
		}
		s = strings.ReplaceAll(s, placeholder, shellQuote(v))
	}
	return s
}

func shellQuote(s string) string {
	if s == "" {
		return "''"
	}
	if !strings.ContainsAny(s, " \t'\"\\$") {
		return s
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
