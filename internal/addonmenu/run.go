package addonmenu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/pick"
	"github.com/Cimmerial/tps/internal/runner"
	"github.com/Cimmerial/tps/internal/ui"
	"gopkg.in/yaml.v3"
)

// RunCommandMenu runs a project command's menu tree until the user backs out.
func RunCommandMenu(p *config.Project, cmd *config.Command, readLine func() (string, error)) (int, error) {
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
	for {
		cur := stack[len(stack)-1]
		title := strings.TrimSpace(cur.Title)
		if title == "" {
			title = cmd.Alias
		}
		ui.ClearScreen()
		fmt.Println()
		fmt.Printf("%s\n", ui.Header(title))
		fmt.Println()
		fmt.Printf("    %s\n", ui.Dim("number to choose — 0 or q = back"))
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

		line, err := readLine()
		if err != nil {
			return 1, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		low := strings.ToLower(line)
		if low == "0" || low == "q" || low == "quit" || low == "back" {
			if len(stack) == 1 {
				return 0, nil
			}
			stack = stack[:len(stack)-1]
			continue
		}
		n, err := parseChoice(line, len(choices))
		if err != nil {
			fmt.Printf("%s\n", ui.Yellow(err.Error()))
			continue
		}
		item := choices[n-1]
		if item.Back {
			if len(stack) == 1 {
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
			code, err := runAction(p, item.Action)
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

func runAction(p *config.Project, action *MenuAction) (int, error) {
	picks, err := collectPicks(action.Picks)
	if err != nil {
		if err == pick.ErrCancelled {
			return 0, nil
		}
		return 1, err
	}
	shell := expandShell(p.Path, action.Shell, picks)
	var code int
	ui.BrailleWaitStderr("running…", func() {
		code, err = runner.RunScript(p, shell)
	})
	return code, err
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
