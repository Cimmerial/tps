package cli

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/config"
	"github.com/Cimmerial/tps/internal/platform"
	"github.com/Cimmerial/tps/internal/ui"
)

func editCaptureSubmenu(reader *bufio.Reader, projectID, alias string) error {
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

		fmt.Println()
		fmt.Printf("%s\n", ui.Header("edit"))
		fmt.Println()
		fmt.Printf("%s\n", ui.SubHeader(alias+" — capture"))
		fmt.Printf("    %s %s\n", ui.Dim("pattern:"), ui.Faint(emptyDash(cmd.Capture)))
		fmt.Println()
		fmt.Printf("    %s\n", ui.Dim("numbered choices — 0 or q = back"))
		fmt.Println()
		captureRows := []numberedMenuRow{
			{Num: "1", Desc: ui.Dim("capture pattern (preset or regex)")},
			{Num: "2", Status: editFieldSetMark(cmd.HasCaptureAction(config.CapturePrint)), Desc: ui.Dim("toggle print")},
			{Num: "3", Status: editFieldSetMark(cmd.HasCaptureAction(config.CaptureClipboard)), Desc: ui.Dim("toggle clipboard")},
		}
		if platform.SupportsCaptureSMS() {
			captureRows = append(captureRows, numberedMenuRow{
				Num: "4", Status: editFieldSetMark(cmd.HasCaptureAction(config.CaptureSMS)),
				Desc: ui.Dim("toggle sms (set phone in options → capture sms)"),
			})
		}
		captureRows = append(captureRows, numberedMenuRow{Num: "0", Desc: ui.Dim("back (q)")})
		printNumberedMenuRows(captureRows)
		fmt.Println()
		if len(cmd.Capture) == 0 {
			fmt.Println(ui.Dim("set a pattern first (e.g. url, trycloudflare)"))
			fmt.Println()
		}

		choice, back, err := readOptionsMenuChoice()
		if err != nil || back {
			return err
		}
		switch choice {
		case 1:
			if err := editPromptCapturePattern(reader, projectID, alias, cmd); err != nil {
				return err
			}
		case 2:
			if err := editToggleCaptureAction(projectID, alias, config.CapturePrint); err != nil {
				return err
			}
		case 3:
			if err := editToggleCaptureAction(projectID, alias, config.CaptureClipboard); err != nil {
				return err
			}
		case 4:
			if !platform.SupportsCaptureSMS() {
				fmt.Println(ui.Dim("invalid choice"))
				continue
			}
			if err := editToggleCaptureAction(projectID, alias, config.CaptureSMS); err != nil {
				return err
			}
		default:
			fmt.Println(ui.Dim("invalid choice"))
		}
	}
}

func editPromptCapturePattern(reader *bufio.Reader, projectID, alias string, cmd *config.Command) error {
	val, ok, err := editPromptField(reader, "capture pattern", cmd.Capture, true)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println(ui.Dim("unchanged"))
		return nil
	}
	return editMutateCommand(projectID, alias, func(c *config.Command) {
		c.Capture = val
		if val == "" {
			c.CaptureActions = nil
		} else {
			c.Detach = true
			if len(c.CaptureActions) == 0 {
				c.CaptureActions = config.DefaultCaptureActions()
			}
		}
	})
}

func editToggleCaptureAction(projectID, alias, action string) error {
	return editMutateCommand(projectID, alias, func(c *config.Command) {
		if strings.TrimSpace(c.Capture) == "" {
			return
		}
		acts := c.CaptureActions
		if len(acts) == 0 {
			acts = config.DefaultCaptureActions()
		}
		on := false
		var next []string
		for _, a := range acts {
			if a == action {
				on = true
				continue
			}
			next = append(next, a)
		}
		if !on {
			next = append(next, action)
		}
		c.CaptureActions = config.ParseCaptureActions(next)
	})
}
