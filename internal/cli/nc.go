package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
)

func isNewCommandToken(s string) bool {
	return strings.EqualFold(s, "nc")
}

func isNewCommandPhrase(a, b string) bool {
	return strings.EqualFold(a, "new") && strings.EqualFold(b, "command")
}

func isNCFlagToken(s string) bool {
	return classifyNCFlag(s) != ""
}

func ncUsageError() error {
	return fmt.Errorf("%s", strings.TrimSpace(`
nc — add a project command
  nc -n name -a abbrev -c "shell"
  shells:  -c / -shell (repeat, or -c2) · -mac / -darwin · -linux · -win / -windows
  options: -desc text · -tags a,b · -precursor abbrev
           -t detach · -close · -H handoff · -fg foreground
           -m pattern · -actions print,clipboard,sms
           -handler id · -menu path
  On each OS, tps runs that OS's shells when any are set, and the -c shells otherwise.
  nc -help`))
}

func ncFlagBody(token string) (string, bool) {
	switch {
	case strings.HasPrefix(token, "--"):
		if len(token) == 2 {
			return "", false
		}
		return token[2:], true
	case strings.HasPrefix(token, "-") && len(token) > 1:
		return token[1:], true
	default:
		return "", false
	}
}

// parseNCShellFlag reports platform ("", "windows", "linux", "darwin") and an
// explicit index (0 means "next free"). -c, -shell, -mac, -linux, -win, and
// numbered forms such as -c2 / -windows_shell2 match.
func parseNCShellFlag(token string) (platform string, idx int, ok bool) {
	body, isFlag := ncFlagBody(token)
	if !isFlag {
		return "", 0, false
	}
	body = strings.ToLower(body)
	prefixes := []struct {
		name string
		plat string
	}{
		{"windows_shell", "windows"},
		{"windows", "windows"},
		{"win", "windows"},
		{"darwin_shell", "darwin"},
		{"darwin", "darwin"},
		{"mac_shell", "darwin"},
		{"mac", "darwin"},
		{"linux_shell", "linux"},
		{"linux", "linux"},
		{"shell", ""},
		{"c", ""},
	}
	for _, p := range prefixes {
		if !strings.HasPrefix(body, p.name) {
			continue
		}
		rest := body[len(p.name):]
		if rest == "" {
			return p.plat, 0, true
		}
		n, err := strconv.Atoi(rest)
		if err != nil || n < 1 || strconv.Itoa(n) != rest {
			continue
		}
		return p.plat, n, true
	}
	return "", 0, false
}

func classifyNCFlag(token string) string {
	if token == "-H" {
		return "handoff"
	}
	if _, _, ok := parseNCShellFlag(token); ok {
		return "shell"
	}
	body, ok := ncFlagBody(token)
	if !ok {
		return ""
	}
	switch strings.ToLower(body) {
	case "n", "name":
		return "name"
	case "a", "abbrev", "abbreviation":
		return "abbrev"
	case "t", "terminal":
		return "terminal"
	case "handoff":
		return "handoff"
	case "precursor":
		return "precursor"
	case "tags", "tag":
		return "tags"
	case "capture", "m":
		return "capture"
	case "desc", "description":
		return "description"
	case "close", "close-on-complete", "close_on_complete":
		return "close"
	case "fg", "foreground":
		return "foreground"
	case "actions", "action", "capture-actions", "capture_actions":
		return "actions"
	case "handler":
		return "handler"
	case "menu":
		return "menu"
	case "help":
		return "help"
	default:
		return ""
	}
}

func readNCValue(fields []string, i int, flag string) (string, int, error) {
	i++
	if i >= len(fields) {
		return "", i, fmt.Errorf("nc: %s requires a value", flag)
	}
	start := i
	for i < len(fields) && !isNCFlagToken(fields[i]) {
		i++
	}
	if start == i {
		return "", i, fmt.Errorf("nc: %s requires a value", flag)
	}
	return strings.TrimSpace(strings.Join(fields[start:i], " ")), i - 1, nil
}

type ncShellSpec struct {
	platform string
	idx      int
	value    string
}

func ncShellFlagLabel(platform string, idx int) string {
	base := "-c"
	switch platform {
	case "windows":
		base = "-win"
	case "linux":
		base = "-linux"
	case "darwin":
		base = "-mac"
	}
	if idx > 1 {
		return fmt.Sprintf("%s%d", base, idx)
	}
	return base
}

func ncYAMLShellKey(platform string, idx int) string {
	if idx < 1 {
		idx = 1
	}
	prefix := ""
	switch platform {
	case "windows":
		prefix = "windows_"
	case "linux":
		prefix = "linux_"
	case "darwin":
		prefix = "darwin_"
	}
	if idx == 1 {
		return prefix + "shell"
	}
	return fmt.Sprintf("%sshell%d", prefix, idx)
}

func applyNCShells(cmd *config.Command, specs []ncShellSpec) error {
	used := map[string]map[int]string{}
	type assigned struct {
		platform string
		idx      int
		value    string
	}
	order := make([]assigned, 0, len(specs))
	for _, sp := range specs {
		if used[sp.platform] == nil {
			used[sp.platform] = map[int]string{}
		}
		idx := sp.idx
		if idx == 0 {
			for i := 1; ; i++ {
				if _, ok := used[sp.platform][i]; !ok {
					idx = i
					break
				}
			}
		} else if _, ok := used[sp.platform][idx]; ok {
			return fmt.Errorf("nc: %s set twice", ncShellFlagLabel(sp.platform, idx))
		}
		used[sp.platform][idx] = sp.value
		order = append(order, assigned{platform: sp.platform, idx: idx, value: sp.value})
	}
	for _, a := range order {
		key := ncYAMLShellKey(a.platform, a.idx)
		if err := cmd.SetShellByKey(key, a.value); err != nil {
			return fmt.Errorf("nc: %s: %w", ncShellFlagLabel(a.platform, a.idx), err)
		}
	}
	return nil
}

func parseNCActions(raw string) ([]string, error) {
	raw = strings.ReplaceAll(raw, " ", ",")
	acts := config.ParseCaptureActions(raw)
	if len(acts) == 0 {
		return nil, fmt.Errorf("nc: -actions expects print, clipboard, sms")
	}
	return acts, nil
}

func parseNewCommandFlags(fields []string) (config.Command, error) {
	if len(fields) == 0 {
		return config.Command{}, ncUsageError()
	}
	var (
		name, abbrev, desc, capture, tagsRaw, precursor, actionsRaw, handler, menu string
		detach, handoff, closeOnComplete, foreground                               bool
		shells                                                                     []ncShellSpec
	)
	for i := 0; i < len(fields); i++ {
		flag := fields[i]
		switch classifyNCFlag(flag) {
		case "help":
			return config.Command{}, ncUsageError()
		case "shell":
			plat, idx, ok := parseNCShellFlag(flag)
			if !ok {
				return config.Command{}, fmt.Errorf("nc: unknown flag %q", flag)
			}
			val, next, err := readNCValue(fields, i, flag)
			if err != nil {
				return config.Command{}, err
			}
			shells = append(shells, ncShellSpec{platform: plat, idx: idx, value: val})
			i = next
		case "name":
			val, next, err := readNCValue(fields, i, flag)
			if err != nil {
				return config.Command{}, err
			}
			if name != "" {
				return config.Command{}, fmt.Errorf("nc: -n set twice")
			}
			name = val
			i = next
		case "abbrev":
			val, next, err := readNCValue(fields, i, flag)
			if err != nil {
				return config.Command{}, err
			}
			if abbrev != "" {
				return config.Command{}, fmt.Errorf("nc: -a set twice")
			}
			abbrev = val
			i = next
		case "description":
			val, next, err := readNCValue(fields, i, flag)
			if err != nil {
				return config.Command{}, err
			}
			if desc != "" {
				return config.Command{}, fmt.Errorf("nc: -desc set twice")
			}
			desc = val
			i = next
		case "precursor":
			val, next, err := readNCValue(fields, i, flag)
			if err != nil {
				return config.Command{}, err
			}
			if precursor != "" {
				return config.Command{}, fmt.Errorf("nc: -precursor set twice")
			}
			precursor = val
			i = next
		case "tags":
			val, next, err := readNCValue(fields, i, flag)
			if err != nil {
				return config.Command{}, err
			}
			if tagsRaw != "" {
				return config.Command{}, fmt.Errorf("nc: -tags set twice")
			}
			tagsRaw = val
			i = next
		case "capture":
			val, next, err := readNCValue(fields, i, flag)
			if err != nil {
				return config.Command{}, err
			}
			if capture != "" {
				return config.Command{}, fmt.Errorf("nc: -m set twice")
			}
			capture = val
			i = next
		case "actions":
			val, next, err := readNCValue(fields, i, flag)
			if err != nil {
				return config.Command{}, err
			}
			if actionsRaw != "" {
				return config.Command{}, fmt.Errorf("nc: -actions set twice")
			}
			actionsRaw = val
			i = next
		case "handler":
			val, next, err := readNCValue(fields, i, flag)
			if err != nil {
				return config.Command{}, err
			}
			if handler != "" {
				return config.Command{}, fmt.Errorf("nc: -handler set twice")
			}
			handler = val
			i = next
		case "menu":
			val, next, err := readNCValue(fields, i, flag)
			if err != nil {
				return config.Command{}, err
			}
			if menu != "" {
				return config.Command{}, fmt.Errorf("nc: -menu set twice")
			}
			menu = val
			i = next
		case "terminal":
			detach = true
		case "handoff":
			handoff = true
		case "close":
			closeOnComplete = true
		case "foreground":
			foreground = true
		default:
			return config.Command{}, fmt.Errorf("nc: unknown flag %q", flag)
		}
	}
	if name == "" {
		return config.Command{}, fmt.Errorf("nc: -n name is required")
	}
	if len(shells) == 0 && handler == "" && menu == "" {
		return config.Command{}, fmt.Errorf("nc: set a shell with -c, -mac, -linux, or -win (or -handler / -menu)")
	}
	if actionsRaw != "" && capture == "" {
		return config.Command{}, fmt.Errorf("nc: -actions needs -m or -capture")
	}
	var actions []string
	if actionsRaw != "" {
		parsed, err := parseNCActions(actionsRaw)
		if err != nil {
			return config.Command{}, err
		}
		actions = parsed
	}
	if capture != "" {
		detach = true
	}
	if closeOnComplete {
		detach = true
	}
	if foreground {
		switch {
		case capture != "":
			return config.Command{}, fmt.Errorf("nc: -fg cannot be used with -m")
		case closeOnComplete:
			return config.Command{}, fmt.Errorf("nc: -fg cannot be used with -close")
		case detach:
			return config.Command{}, fmt.Errorf("nc: -fg cannot be used with -t")
		case handoff:
			return config.Command{}, fmt.Errorf("nc: -fg cannot be used with -H")
		}
	}
	if closeOnComplete && handoff {
		return config.Command{}, fmt.Errorf("nc: -close cannot be used with -H")
	}
	if handoff {
		detach = false
	}
	if desc == "" {
		desc = name
	}
	cmd := config.Command{
		Alias:           name,
		Abbrev:          abbrev,
		Description:     desc,
		Detach:          detach,
		ForceForeground: foreground,
		Handoff:         handoff,
		CloseOnComplete: closeOnComplete,
		Capture:         capture,
		CaptureActions:  actions,
		Precursor:       precursor,
		Handler:         handler,
		Menu:            menu,
	}
	if tagsRaw != "" {
		cmd.Tags = config.ParseTagList(tagsRaw)
	}
	if err := applyNCShells(&cmd, shells); err != nil {
		return config.Command{}, err
	}
	return cmd, nil
}

func addProjectCommand(cfg *config.Config, projectID string, cmd config.Command) error {
	p := cfg.FindByID(projectID)
	if p == nil {
		return fmt.Errorf("project not found")
	}
	if err := p.AddCommand(cmd); err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return err
	}
	return nil
}
