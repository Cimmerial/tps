package git

import (
	"fmt"
	"strings"

	"github.com/Cimmerial/tps/internal/ui"
)

// SuspiciousAction is the user's choice when ship would commit sensitive paths.
type SuspiciousAction int

const (
	SuspiciousProceed          SuspiciousAction = 1
	SuspiciousCancel           SuspiciousAction = 2
	SuspiciousExcludeGitignore SuspiciousAction = 3
)

// ParseShipSuspiciousChoice maps interactive input to an action. ok is false when unrecognized.
func ParseShipSuspiciousChoice(line string) (SuspiciousAction, bool) {
	line = strings.TrimSpace(strings.ToLower(line))
	switch line {
	case "1", "y", "yes":
		return SuspiciousProceed, true
	case "2", "n", "no":
		return SuspiciousCancel, true
	case "3":
		return SuspiciousExcludeGitignore, true
	default:
		return 0, false
	}
}

// PromptSuspiciousFiles lists paths and asks how to proceed (1/y, 2/n, 3).
func PromptSuspiciousFiles(paths []string) (SuspiciousAction, error) {
	fmt.Println()
	fmt.Println(ui.Yellow("ship would commit suspicious files:"))
	for _, p := range paths {
		fmt.Printf("    %s\n", ui.Bold(p))
	}
	fmt.Println()
	fmt.Printf("    %s | %s | %s\n", ui.Bold("1"), ui.Abbrev("y"), "confirm and ship anyway")
	fmt.Printf("    %s | %s | %s\n", ui.Bold("2"), ui.Abbrev("n"), "cancel")
	fmt.Printf("    %s | %s | %s\n", ui.Bold("3"), ui.Abbrev(" "), "commit without these files and add them to .gitignore")
	fmt.Println()
	for {
		line, err := ui.ReadEditedLine(ui.Prompt())
		if err != nil {
			return 0, err
		}
		if act, ok := ParseShipSuspiciousChoice(line); ok {
			return act, nil
		}
		fmt.Println(ui.Dim("enter 1, 2, 3, y, or n"))
	}
}
