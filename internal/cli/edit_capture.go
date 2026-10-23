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

		active := cmd.CaptureActive()
		fmt.Println()
		fmt.Printf("%s\n", ui.HeaderChoice("edit"))
		fmt.Printf("%s\n", ui.SubHeaderChoice(alias+" — capture"))
		fmt.Printf("    %s %s\n", ui.Dim("pattern:"), ui.BrightBlue(emptyDash(cmd.Capture)))
		fmt.Println()

		captureRows := []numberedMenuRow{
			{Num: "1", Status: editFieldSetMark(active), Desc: "toggle capture"},
			{Num: "2", Desc: "capture pattern (preset or regex)"},
		}
		if active {
			captureRows = append(captureRows,
				numberedMenuRow{Num: "3", Status: editFieldSetMark(cmd.HasCaptureAction(config.CapturePrint)), Desc: "toggle print"},
				numberedMenuRow{Num: "4", Status: editFieldSetMark(cmd.HasCaptureAction(config.CaptureClipboard)), Desc: "toggle clipboard"},
			)
			if platform.SupportsCaptureSMS() {
				captureRows = append(captureRows, numberedMenuRow{
					Num: "5", Status: editFieldSetMark(cmd.HasCaptureAction(config.CaptureSMS)),
					Desc: "toggle sms (set phone in options → capture sms)",
				})
			}
		}
		printChoiceEditMenuPanel("", []menuPanelSection{{Rows: captureRows}}, nil)
		fmt.Println()
		printChoiceEditMenuPanel("", nil, []numberedMenuRow{
			{Num: "0", Desc: ui.Dim("back (q)")},
		})
		fmt.Println()
		if !active && strings.TrimSpace(cmd.Capture) == "" {
			fmt.Println(ui.Dim("set a pattern first (e.g. url, trycloudflare)"))
			fmt.Println()
		}

		choice, back, err := readOptionsMenuChoice()
		if err != nil || back {
			return err
		}
		switch choice {
		case 1:
			if err := editToggleCaptureEnabled(projectID, alias, cmd); err != nil {
				return err
			}
		case 2:
			if err := editPromptCapturePattern(reader, projectID, alias, cmd); err != nil {
				return err
			}
		case 3:
			if !active {
				fmt.Println(ui.Dim("invalid choice"))
				continue
			}
			if err := editToggleCaptureAction(projectID, alias, config.CapturePrint); err != nil {
				return err
			}
		case 4:
			if !active {
				fmt.Println(ui.Dim("invalid choice"))
				continue
			}
			if err := editToggleCaptureAction(projectID, alias, config.CaptureClipboard); err != nil {
				return err
			}
		case 5:
			if !active {
				fmt.Println(ui.Dim("invalid choice"))
				continue
			}
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

func editToggleCaptureEnabled(projectID, alias string, cmd *config.Command) error {
	next := !cmd.CaptureActive()
	return editMutateCommand(projectID, alias, func(c *config.Command) {
		if strings.TrimSpace(c.Capture) == "" {
			return
		}
		if next {
			c.CaptureDisabled = false
			c.Detach = true
			if len(c.CaptureActions) == 0 {
				c.CaptureActions = config.DefaultCaptureActions()
			}
		} else {
			c.CaptureDisabled = true
		}
	})
}

func editPromptCapturePattern(reader *bufio.Reader, projectID, alias string, cmd *config.Command) error {
	val, ok, err := editPromptField(reader, "capture pattern", cmd.Capture, true, false)
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
			c.CaptureDisabled = false
			return
		}
		c.CaptureDisabled = false
		c.Detach = true
		if len(c.CaptureActions) == 0 {
			c.CaptureActions = config.DefaultCaptureActions()
		}
	})
}

func editToggleCaptureAction(projectID, alias, action string) error {
	return editMutateCommand(projectID, alias, func(c *config.Command) {
		if !c.CaptureActive() {
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
