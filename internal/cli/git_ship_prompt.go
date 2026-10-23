package cli

import (
	"errors"
	"fmt"

	gittool "github.com/Cimmerial/tps/internal/git"
	"github.com/Cimmerial/tps/internal/ui"
)

func promptShipSuspiciousFiles(paths []string) (gittool.SuspiciousAction, error) {
	fmt.Println()
	fmt.Println(ui.Yellow("ship would commit suspicious files:"))
	for _, p := range paths {
		fmt.Printf("    %s\n", ui.Bold(p))
	}
	fmt.Println()
	printNumberedMenuRows([]numberedMenuRow{
		{Num: "1", Abbrev: "y", Desc: "confirm and ship anyway"},
		{Num: "2", Abbrev: "n", Desc: "cancel"},
		{Num: "3", Abbrev: "", Desc: "commit without these files and add them to .gitignore"},
	})
	fmt.Println()
	for {
		line, err := readMenuLine()
		if err != nil {
			return 0, err
		}
		if act, ok := gittool.ParseShipSuspiciousChoice(line); ok {
			return act, nil
		}
		fmt.Println(ui.Dim("enter 1, 2, 3, y, or n"))
	}
}

func promptShipFileSelect(dir string, paths []string) ([]string, error) {
	sel := gittool.NewShipFileSelection(paths)
	stats := buildShipFileStats(dir, paths)
	diffs := buildShipFileDiffs(dir, paths, stats)
	for {
		ui.ClearScreen()
		printShipFileBoxes(sel, stats, diffs)
		fmt.Println()
		fmt.Println(ui.Dim("Press c to commit, Enter or 0 to quit, or enter a number to toggle it."))
		fmt.Println()
		line, err := readMenuLine()
		if err != nil {
			if errors.Is(err, ui.ErrPromptCancelled) {
				return nil, gittool.ErrShipSelectCancel
			}
			return nil, err
		}
		commit, quit, num, ok := gittool.ParseShipFileSelectInput(line)
		if !ok {
			fmt.Println(ui.Dim("enter c to commit, 0 or Enter to quit, or a number to toggle"))
			continue
		}
		if quit {
			return nil, gittool.ErrShipSelectCancel
		}
		if commit {
			return sel.IncludedPaths(), nil
		}
		path, okPath := shipFilePathForDisplayNumber(sel.Paths, num)
		if !okPath || !sel.TogglePath(path) {
			fmt.Println(ui.Dim("invalid number"))
			continue
		}
	}
}

func printShipFileBoxes(sel *gittool.ShipFileSelection, stats map[string]shipFileStat, diffs map[string]shipFileDiff) {
	displayNums := shipFileDisplayNumbers(sel.Paths)
	maxNum := shipFileMaxNumWidth(len(sel.Paths))
	fmt.Println()
	printShipFileBox(ui.ChromeCommit, ui.HeaderCommit("BEING COMMITTED:"), sel, true, maxNum, stats, diffs, displayNums)
	fmt.Println()
	printShipFileBox(ui.ChromeChoice, ui.HeaderChoice("NOT BEING COMMITTED:"), sel, false, maxNum, stats, diffs, displayNums)
}

func shipFileMaxNumWidth(fileCount int) int {
	if fileCount < 1 {
		return 0
	}
	return len(fmt.Sprintf("%d", fileCount))
}

func pathsFromEntries(entries []shipFileTreeEntry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Path
	}
	return out
}

func printShipFileBox(mode ui.ChromeMode, title string, sel *gittool.ShipFileSelection, included bool, maxNum int, stats map[string]shipFileStat, diffs map[string]shipFileDiff, displayNums map[string]int) {
	box := ui.NewMenuBox(mode)
	box.Title(title)
	var entries []shipFileTreeEntry
	for i, p := range sel.Paths {
		if sel.Included[i] != included {
			continue
		}
		entries = append(entries, shipFileTreeEntry{Path: p})
	}
	if len(entries) == 0 {
		box.Plain(ui.SectionTrail{LastSection: true, LastRow: true}, ui.Dim("(none)"))
	} else {
		if included {
			added, removed := sumShipFileDiffs(pathsFromEntries(entries), diffs)
			box.Plain(ui.SectionTrail{LastSection: true, LastRow: false},
				shipFileDiffTotalsLabel(added, removed))
		}
		treeLines := shipFileTreeLines(buildShipFileTree(entries), func(name string) string {
			return ui.Blue(name)
		}, func(_ int, fileName, relPath string) string {
			num := padPlain(fmt.Sprintf("%d", displayNums[relPath]), maxNum)
			line := ui.Bold(ui.Red(num)) + ui.Dim("  |  ") + fileName
			if included {
				if d, ok := diffs[relPath]; ok {
					line += ui.Dim(" · ") + shipFileDiffLabel(d)
				}
			} else if st, ok := stats[relPath]; ok {
				line += ui.Dim(" · ") + ui.Gray(st.label())
			}
			return line
		})
		for i, line := range treeLines {
			trail := ui.SectionTrail{LastSection: true, LastRow: i == len(treeLines)-1}
			box.Plain(trail, line)
		}
	}
	box.Flush()
}

func runGitShip(dir, message string) (int, error) {
	opts := &gittool.ShipOptions{
		OnSuspicious: promptShipSuspiciousFiles,
	}
	if ui.StdinIsTTY() {
		opts.OnSelect = func(paths []string) ([]string, error) {
			return promptShipFileSelect(dir, paths)
		}
	}
	return gittool.ShipWithOptions(dir, message, opts)
}
