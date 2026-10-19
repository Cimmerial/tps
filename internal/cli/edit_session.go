package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/runner"
	"github.com/Cimmerial/tps/internal/ui"
)

type commandSnap struct {
	Shell       string
	Abbrev      string
	Description string
	Precursor   string
	Capture     string
	Detach      bool
	Tags        []string
}

func snapCommand(c config.Command) commandSnap {
	return commandSnap{
		Shell:       c.Shell,
		Abbrev:      c.Abbrev,
		Description: c.Description,
		Precursor:   c.Precursor,
		Capture:     c.Capture,
		Detach:      c.Detach,
		Tags:        append([]string(nil), c.Tags...),
	}
}

func snapshotProjectCommands(cfg *config.Config, projectID string) (map[string]commandSnap, error) {
	p := cfg.FindByID(projectID)
	if p == nil {
		return nil, fmt.Errorf("project not found")
	}
	out := make(map[string]commandSnap, len(p.Commands))
	for _, c := range p.Commands {
		out[c.Alias] = snapCommand(c)
	}
	return out, nil
}

func sameTagSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func snapsEqual(a, b commandSnap) bool {
	return a.Shell == b.Shell &&
		a.Abbrev == b.Abbrev &&
		a.Description == b.Description &&
		a.Precursor == b.Precursor &&
		a.Capture == b.Capture &&
		a.Detach == b.Detach &&
		sameTagSlice(a.Tags, b.Tags)
}

func editCommandSnapDiffs(label string, b, a commandSnap) []string {
	var lines []string
	if a.Shell != b.Shell {
		lines = append(lines, formatPrefChange(label+" shell", b.Shell, a.Shell))
	}
	if a.Abbrev != b.Abbrev {
		lines = append(lines, formatPrefChange(label+" abbrev", emptyDash(b.Abbrev), emptyDash(a.Abbrev)))
	}
	if a.Description != b.Description {
		lines = append(lines, formatPrefChange(label+" description", emptyDash(b.Description), emptyDash(a.Description)))
	}
	if a.Precursor != b.Precursor {
		lines = append(lines, formatPrefChange(label+" precursor", emptyDash(b.Precursor), emptyDash(a.Precursor)))
	}
	if a.Capture != b.Capture {
		lines = append(lines, formatPrefChange(label+" capture", emptyDash(b.Capture), emptyDash(a.Capture)))
	}
	if a.Detach != b.Detach {
		lines = append(lines, formatPrefChange(label+" detach", boolLabel(b.Detach), boolLabel(a.Detach)))
	}
	if !sameTagSlice(a.Tags, b.Tags) {
		lines = append(lines, formatPrefChange(label+" tags", tagsLabel(b.Tags), tagsLabel(a.Tags)))
	}
	return lines
}

func editSessionChanges(before, after map[string]commandSnap) []string {
	var lines []string
	var added []string
	var removed []string

	for alias, a := range after {
		b, ok := before[alias]
		if !ok {
			added = append(added, alias)
			continue
		}
		lines = append(lines, editCommandSnapDiffs(alias, b, a)...)
	}
	for alias := range before {
		if _, ok := after[alias]; !ok {
			removed = append(removed, alias)
		}
	}

	usedAdded := make(map[string]bool)
	usedRemoved := make(map[string]bool)
	for _, oldAlias := range removed {
		oldSnap := before[oldAlias]
		for _, newAlias := range added {
			if usedAdded[newAlias] {
				continue
			}
			newSnap := after[newAlias]
			if snapsEqual(oldSnap, newSnap) {
				lines = append(lines, formatPrefChange("command", oldAlias, newAlias))
				usedAdded[newAlias] = true
				usedRemoved[oldAlias] = true
				break
			}
			// Renamed and edited in one session: same abbrev, shell unchanged enough to pair.
			if oldSnap.Abbrev != "" && oldSnap.Abbrev == newSnap.Abbrev && oldSnap.Shell == newSnap.Shell {
				lines = append(lines, formatPrefChange("command", oldAlias, newAlias))
				lines = append(lines, editCommandSnapDiffs(newAlias, oldSnap, newSnap)...)
				usedAdded[newAlias] = true
				usedRemoved[oldAlias] = true
				break
			}
		}
	}
	for _, oldAlias := range removed {
		if !usedRemoved[oldAlias] {
			lines = append(lines, ui.Red(fmt.Sprintf("removed %q", oldAlias)))
		}
	}
	for _, newAlias := range added {
		if !usedAdded[newAlias] {
			lines = append(lines, ui.Green(fmt.Sprintf("added %q", newAlias)))
		}
	}
	return lines
}

func emptyDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}

func boolLabel(v bool) string {
	if v {
		return "on"
	}
	return "off"
}

func tagsLabel(tags []string) string {
	if len(tags) == 0 {
		return "—"
	}
	return strings.Join(tags, ", ")
}

func finishEditSession(projectID string, statusTick int, start map[string]commandSnap) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	end, err := snapshotProjectCommands(cfg, projectID)
	if err != nil {
		return err
	}
	if err := refreshInteractiveSession(projectID, statusTick, ""); err != nil {
		return err
	}
	changes := editSessionChanges(start, end)
	if len(changes) == 0 {
		return nil
	}
	fmt.Println()
	fmt.Println(ui.Header("EDIT CHANGES:"))
	for _, line := range changes {
		fmt.Printf("  - %s\n", line)
	}
	return nil
}

func editCommandList(p *config.Project) []config.Command {
	var list []config.Command
	for _, g := range p.CommandsByTag() {
		list = append(list, g.Commands...)
	}
	return list
}

func runEditMenu(reader *bufio.Reader, cfg *config.Config, projectID string, statusTick int, jumpAlias string) error {
	start, err := snapshotProjectCommands(cfg, projectID)
	if err != nil {
		return err
	}

	if jumpAlias != "" {
		p := cfg.FindByID(projectID)
		if p == nil {
			return fmt.Errorf("project not found")
		}
		if resolveProjectCommand(p, jumpAlias) == nil {
			return fmt.Errorf("no command %q", jumpAlias)
		}
		if err := editCommandSubmenu(reader, projectID, jumpAlias); err != nil {
			return err
		}
		return finishEditSession(projectID, statusTick, start)
	}

	for {
		ui.ClearScreen()
		cfg, p, err := loadProject(projectID)
		if err != nil {
			return err
		}
		if len(p.Commands) == 0 {
			fmt.Println()
			fmt.Printf("%s\n", ui.Header("edit"))
			fmt.Println()
			fmt.Println(ui.Dim("no commands yet"))
			fmt.Println()
			printNumberedMenuRows([]numberedMenuRow{
				{Num: "1", Desc: ui.Dim("add command")},
				{Num: "0", Abbrev: "q", Desc: ui.Dim("done")},
			})
			fmt.Println()
			choice, back, err := readOptionsMenuChoice()
			if err != nil || back {
				return finishEditSession(projectID, statusTick, start)
			}
			if choice == 1 {
				if err := editAddCommand(reader, cfg, projectID); err != nil {
					fmt.Println(ui.Dim(fmt.Sprintf("add: %v", err)))
				}
				continue
			}
			fmt.Println(ui.Dim("invalid choice"))
			continue
		}

		cmds := editCommandList(p)
		fmt.Println()
		fmt.Printf("%s\n", ui.Header("edit"))
		fmt.Println()
		fmt.Printf("    %s\n", ui.Dim("numbered choices — 0 or q = done"))
		fmt.Println()
		rows := make([]numberedMenuRow, len(cmds))
		for i, c := range cmds {
			desc := editCommandRowDesc(c)
			rows[i] = numberedMenuRow{
				Num:    fmt.Sprintf("%d", i+1),
				Abbrev: c.Abbrev,
				Desc:   desc,
			}
		}
		printNumberedMenuRows(rows)
		fmt.Println()
		addNum := len(cmds) + 1
		tagNum := len(cmds) + 2
		printNumberedMenuRows([]numberedMenuRow{
			{Num: fmt.Sprintf("%d", addNum), Desc: ui.Dim("add command")},
			{Num: fmt.Sprintf("%d", tagNum), Desc: ui.Dim("manage tags")},
			{Num: "0", Abbrev: "q", Desc: ui.Dim("done")},
		})
		fmt.Println()

		choice, back, err := readOptionsMenuChoice()
		if err != nil {
			return err
		}
		if back {
			return finishEditSession(projectID, statusTick, start)
		}
		if choice >= 1 && choice <= len(cmds) {
			if err := editCommandSubmenu(reader, projectID, cmds[choice-1].Alias); err != nil {
				return err
			}
			continue
		}
		if choice == addNum {
			if err := editAddCommand(reader, cfg, projectID); err != nil {
				fmt.Println(ui.Dim(fmt.Sprintf("add: %v", err)))
			}
			continue
		}
		if choice == tagNum {
			if err := editManageTags(reader, projectID); err != nil {
				return err
			}
			continue
		}
		fmt.Println(ui.Dim("invalid choice"))
	}
}

func editCommandRowDesc(c config.Command) string {
	parts := []string{ui.Bold(c.Alias)}
	if c.Abbrev != "" {
		parts = append(parts, ui.Dim("("+c.Abbrev+")"))
	}
	shell := c.Shell
	if c.HasMultiShell() {
		shells := c.ResolvedShells()
		if len(shells) > 0 {
			shell = strings.Join(shells, " · ")
		}
		parts = append(parts, ui.Dim("— "+shell+" (multi-shell)"))
	} else if shell != "" {
		parts = append(parts, ui.Dim("— "+shell))
	}
	if tag := c.PrimaryTag(); tag != "" {
		parts = append(parts, ui.Dim("["+tag+"]"))
	}
	return strings.Join(parts, " ")
}

func editCommandSubmenu(reader *bufio.Reader, projectID, alias string) error {
	for {
		ui.ClearScreen()
		_, p, err := loadProject(projectID)
		if err != nil {
			return err
		}
		cmd := resolveProjectCommand(p, alias)
		if cmd == nil {
			return fmt.Errorf("no command %q", alias)
		}
		alias = cmd.Alias

		fmt.Println()
		fmt.Printf("%s\n", ui.Header("edit"))
		fmt.Println()
		fmt.Printf("%s\n", ui.SubHeader(alias))
		fmt.Println()
		fmt.Printf("    %s\n", ui.Dim("numbered choices — 0 or q = back"))
		fmt.Println()
		subRows := editCommandSubmenuRows(*cmd)
		subRows = append(subRows, numberedMenuRow{Num: "0", Desc: ui.Dim("back (q)")})
		printNumberedMenuRows(subRows)
		fmt.Println()

		choice, back, err := readOptionsMenuChoice()
		if err != nil {
			return err
		}
		if back {
			return nil
		}
		switch choice {
		case 1:
			if err := editShellSubmenu(reader, projectID, alias); err != nil {
				return err
			}
		case 2:
			newAlias, err := editPromptName(reader, projectID, alias, cmd)
			if err != nil {
				return err
			}
			if newAlias != "" {
				alias = newAlias
			}
		case 3:
			if err := editPromptAbbrev(reader, projectID, alias, cmd); err != nil {
				return err
			}
		case 4:
			if err := editPromptDescription(reader, projectID, alias, cmd); err != nil {
				return err
			}
		case 5:
			if err := editPromptTags(reader, projectID, alias, cmd); err != nil {
				return err
			}
		case 6:
			if err := editPromptPrecursor(reader, projectID, alias, cmd); err != nil {
				return err
			}
		case 7:
			if err := editToggleDetach(projectID, alias, cmd); err != nil {
				return err
			}
		case 8:
			if err := editPromptCapture(reader, projectID, alias, cmd); err != nil {
				return err
			}
		case 9:
			confirmed, err := editConfirmDelete(reader, alias)
			if err != nil {
				return err
			}
			if confirmed {
				if err := editDeleteCommand(projectID, alias); err != nil {
					return err
				}
				return nil
			}
		default:
			fmt.Println(ui.Dim("invalid choice"))
		}
	}
}

func editTruncateDisplay(s string, max int) string {
	s = strings.TrimSpace(s)
	if max <= 1 || len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func editFieldSetMark(set bool) string {
	if set {
		return ui.Green("✓")
	}
	return ui.Red("✗")
}

func editDescriptionIsSet(cmd config.Command) bool {
	d := strings.TrimSpace(cmd.Description)
	return d != "" && d != cmd.Alias
}

func editShellHelpLine(cmd config.Command) string {
	fields := cmd.AllShellFields()
	if len(fields) == 0 {
		return ""
	}
	if len(fields) == 1 {
		return editTruncateDisplay(fields[0].Value, 72)
	}
	return editTruncateDisplay(fields[0].Value, 48) + fmt.Sprintf(" (+%d more)", len(fields)-1)
}

func editCommandSubmenuRows(cmd config.Command) []numberedMenuRow {
	shellSet := cmd.HasAnyShellLine()
	nameSet := strings.TrimSpace(cmd.Alias) != ""
	abbrevSet := cmd.Abbrev != ""
	descSet := editDescriptionIsSet(cmd)
	tagsSet := len(cmd.Tags) > 0
	precSet := cmd.Precursor != ""
	capSet := cmd.Capture != ""
	detached, detachHelp := runner.DetachInfo(cmd)
	return []numberedMenuRow{
		{Num: "1", Status: editFieldSetMark(shellSet), Desc: ui.Dim("shell"), HelpLine: editHelpIf(shellSet, editShellHelpLine(cmd))},
		{Num: "2", Status: editFieldSetMark(nameSet), Desc: ui.Dim("name"), HelpLine: editHelpIf(nameSet, cmd.Alias)},
		{Num: "3", Status: editFieldSetMark(abbrevSet), Desc: ui.Dim("abbreviation"), HelpLine: editHelpIf(abbrevSet, cmd.Abbrev)},
		{Num: "4", Status: editFieldSetMark(descSet), Desc: ui.Dim("description"), HelpLine: editHelpIf(descSet, cmd.Description)},
		{Num: "5", Status: editFieldSetMark(tagsSet), Desc: ui.Dim("tags"), HelpLine: editHelpIf(tagsSet, tagsLabel(cmd.Tags))},
		{Num: "6", Status: editFieldSetMark(precSet), Desc: ui.Dim("precursor"), HelpLine: editHelpIf(precSet, cmd.Precursor)},
		{Num: "7", Status: editFieldSetMark(detached), Desc: ui.Dim("detach / terminal"), HelpLine: detachHelp},
		{Num: "8", Status: editFieldSetMark(capSet), Desc: ui.Dim("capture"), HelpLine: editHelpIf(capSet, cmd.Capture)},
		{Num: "9", Desc: ui.Dim("delete command")},
	}
}

func editHelpIf(set bool, val string) string {
	if !set {
		return ""
	}
	return editTruncateDisplay(val, 56)
}

func editPromptErr(err error) error {
	if err == nil {
		return nil
	}
	if err.Error() == "interrupted" {
		fmt.Println(ui.Dim("cancelled"))
		return nil
	}
	return err
}

func editMutateCommand(projectID, alias string, mutate func(*config.Command)) error {
	cfg, p, err := loadProject(projectID)
	if err != nil {
		return err
	}
	cmd := resolveProjectCommand(p, alias)
	if cmd == nil {
		return fmt.Errorf("no command %q", alias)
	}
	updated := *cmd
	mutate(&updated)
	if err := p.AddCommand(updated); err != nil {
		return err
	}
	return config.Save(cfg)
}

func editPromptField(reader *bufio.Reader, field, current string, allowClear bool) (string, bool, error) {
	fmt.Println()
	cur := emptyDash(current)
	if cur != "—" {
		fmt.Printf("current %s: %s\n", field, ui.Faint(cur))
	} else {
		fmt.Printf("current %s: %s\n", field, ui.Dim("—"))
	}
	if allowClear {
		fmt.Println(ui.Dim("Enter to keep · - to clear"))
	} else {
		fmt.Println(ui.Dim("Enter to keep"))
	}
	line, err := ui.ReadEditedLine(fmt.Sprintf("new %s: ", field))
	if err != nil {
		return "", false, editPromptErr(err)
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return "", false, nil
	}
	if line == "-" {
		if !allowClear {
			fmt.Println(ui.Dim("unchanged"))
			return "", false, nil
		}
		return "", true, nil
	}
	return line, true, nil
}

func editPromptName(reader *bufio.Reader, projectID, alias string, cmd *config.Command) (string, error) {
	val, ok, err := editPromptField(reader, "name", cmd.Alias, false)
	if err != nil {
		return "", err
	}
	if !ok {
		fmt.Println(ui.Dim("unchanged"))
		return "", nil
	}
	val = strings.TrimSpace(val)
	if val == "" {
		fmt.Println(ui.Dim("unchanged"))
		return "", nil
	}
	cfg, p, err := loadProject(projectID)
	if err != nil {
		return "", err
	}
	if err := p.RenameCommand(alias, val); err != nil {
		fmt.Println(ui.Dim(err.Error()))
		return "", nil
	}
	if err := config.Save(cfg); err != nil {
		return "", err
	}
	fmt.Println(ui.Green("name updated"))
	return val, nil
}

func editPromptAbbrev(reader *bufio.Reader, projectID, alias string, cmd *config.Command) error {
	val, ok, err := editPromptField(reader, "abbreviation", cmd.Abbrev, true)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println(ui.Dim("unchanged"))
		return nil
	}
	return editSaveStringField(projectID, alias, func(c *config.Command) { c.Abbrev = val }, "abbrev")
}

func editPromptDescription(reader *bufio.Reader, projectID, alias string, cmd *config.Command) error {
	val, ok, err := editPromptField(reader, "description", cmd.Description, true)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println(ui.Dim("unchanged"))
		return nil
	}
	if val == "" {
		return editMutateCommand(projectID, alias, func(c *config.Command) {
			c.Description = c.Alias
		})
	}
	return editSaveStringField(projectID, alias, func(c *config.Command) { c.Description = val }, "description")
}

func editPromptTags(reader *bufio.Reader, projectID, alias string, cmd *config.Command) error {
	val, ok, err := editPromptField(reader, "tags", tagsLabel(cmd.Tags), true)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println(ui.Dim("unchanged"))
		return nil
	}
	var tags []string
	if val != "" {
		tags = config.ParseTagList(val)
	}
	return editMutateCommand(projectID, alias, func(c *config.Command) {
		c.Tags = tags
	})
}

func editPromptPrecursor(reader *bufio.Reader, projectID, alias string, cmd *config.Command) error {
	val, ok, err := editPromptField(reader, "precursor", cmd.Precursor, true)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println(ui.Dim("unchanged"))
		return nil
	}
	return editMutateCommand(projectID, alias, func(c *config.Command) {
		c.Precursor = strings.TrimSpace(val)
	})
}

func editPromptCapture(reader *bufio.Reader, projectID, alias string, cmd *config.Command) error {
	val, ok, err := editPromptField(reader, "capture", cmd.Capture, true)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println(ui.Dim("unchanged"))
		return nil
	}
	return editMutateCommand(projectID, alias, func(c *config.Command) {
		c.Capture = strings.TrimSpace(val)
		if c.Capture != "" {
			c.Detach = true
		}
	})
}

func editSaveStringField(projectID, alias string, set func(*config.Command), label string) error {
	if err := editMutateCommand(projectID, alias, set); err != nil {
		return err
	}
	fmt.Println(ui.Green(label + " updated"))
	return nil
}

func editToggleDetach(projectID, alias string, cmd *config.Command) error {
	nextDetached := !runner.ShouldDetach(*cmd)
	if err := editMutateCommand(projectID, alias, func(c *config.Command) {
		if nextDetached {
			c.ForceForeground = false
			c.Detach = true
		} else {
			c.ForceForeground = true
			c.Detach = false
		}
	}); err != nil {
		return err
	}
	after := *cmd
	if nextDetached {
		after.ForceForeground = false
		after.Detach = true
	} else {
		after.ForceForeground = true
		after.Detach = false
	}
	_, hint := runner.DetachInfo(after)
	fmt.Println(ui.Dim(hint))
	return nil
}

func editConfirmDelete(reader *bufio.Reader, alias string) (bool, error) {
	ui.ClearScreen()
	fmt.Println()
	fmt.Printf("%s\n", ui.Header("edit"))
	fmt.Println()
	fmt.Printf("    %s\n", ui.Red(fmt.Sprintf("delete command %q?", alias)))
	fmt.Println()
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "1", Desc: ui.Dim("yes, delete")},
		{Num: "0", Abbrev: "q", Desc: ui.Dim("cancel")},
	})
	fmt.Println()
	choice, back, err := readOptionsMenuChoice()
	if err != nil || back || choice != 1 {
		return false, err
	}
	return true, nil
}

func editDeleteCommand(projectID, alias string) error {
	cfg, p, err := loadProject(projectID)
	if err != nil {
		return err
	}
	if !p.RemoveCommand(alias) {
		return fmt.Errorf("no command %q", alias)
	}
	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Println(ui.Red(fmt.Sprintf("removed %q", alias)))
	return nil
}

func editAddCommand(reader *bufio.Reader, cfg *config.Config, projectID string) error {
	name, err := ui.ReadEditedLine("name (-n): ")
	if err != nil {
		return editPromptErr(err)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("name is required")
	}
	abbrev, err := ui.ReadEditedLine("abbrev (-a) [optional]: ")
	if err != nil {
		return editPromptErr(err)
	}
	shell, err := ui.ReadEditedLine("shell (-c): ")
	if err != nil {
		return editPromptErr(err)
	}
	shell = strings.TrimSpace(shell)
	if shell == "" {
		return fmt.Errorf("shell is required")
	}
	cmd := config.Command{
		Alias:       name,
		Abbrev:      strings.TrimSpace(abbrev),
		Shell:       shell,
		Description: name,
	}
	p := cfg.FindByID(projectID)
	if p == nil {
		return fmt.Errorf("project not found")
	}
	if resolveProjectCommand(p, name) != nil {
		return fmt.Errorf("command %q already exists", name)
	}
	if err := p.AddCommand(cmd); err != nil {
		return err
	}
	if err := config.Save(cfg); err != nil {
		return err
	}
	fmt.Println(ui.Green(fmt.Sprintf("added %q", name)))
	return nil
}

func editManageTags(reader *bufio.Reader, projectID string) error {
	for {
		ui.ClearScreen()
		_, p, err := loadProject(projectID)
		if err != nil {
			return err
		}
		tags := p.ProjectTags()
		fmt.Println()
		fmt.Printf("%s\n", ui.Header("edit"))
		fmt.Println()
		fmt.Printf("%s\n", ui.SubHeader("manage tags"))
		fmt.Println()
		if len(tags) == 0 {
			fmt.Println(ui.Dim("no tags in use (set tags on a command first)"))
			fmt.Println()
			printNumberedMenuRows([]numberedMenuRow{
				{Num: "0", Abbrev: "q", Desc: ui.Dim("back")},
			})
			fmt.Println()
			_, _, err := readOptionsMenuChoice()
			return err
		}
		fmt.Printf("    %s\n", ui.Dim("pick a tag to rename — 0 or q = back"))
		fmt.Println()
		rows := make([]numberedMenuRow, len(tags))
		for i, t := range tags {
			rows[i] = numberedMenuRow{Num: fmt.Sprintf("%d", i+1), Desc: t}
		}
		printNumberedMenuRows(rows)
		fmt.Println()
		printNumberedMenuRows([]numberedMenuRow{
			{Num: "0", Abbrev: "q", Desc: ui.Dim("back")},
		})
		fmt.Println()

		choice, back, err := readOptionsMenuChoice()
		if err != nil || back {
			return err
		}
		if choice < 1 || choice > len(tags) {
			fmt.Println(ui.Dim("invalid choice"))
			continue
		}
		old := tags[choice-1]
		line, err := ui.ReadEditedLine(fmt.Sprintf("rename %q to: ", old))
		if err != nil {
			return editPromptErr(err)
		}
		newName := strings.TrimSpace(line)
		if newName == "" {
			fmt.Println(ui.Dim("unchanged"))
			continue
		}
		cfg, p, err := loadProject(projectID)
		if err != nil {
			return err
		}
		n := p.RenameTag(old, newName)
		if n == 0 {
			fmt.Println(ui.Dim("no changes"))
			continue
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Println(ui.Green(fmt.Sprintf("renamed tag on %d command(s)", n)))
	}
}
