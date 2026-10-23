package cli

import (
	"fmt"

	gittool "github.com/Cimmerial/tps/internal/git"
	"github.com/Cimmerial/tps/internal/ui"
)

type shipFileDiff struct {
	Added   int
	Removed int
	Binary  bool
}

func buildShipFileDiffs(repoDir string, paths []string, fileStats map[string]shipFileStat) map[string]shipFileDiff {
	raw, err := gittool.DiffNumstatAgainstHEAD(repoDir)
	if err != nil {
		raw = map[string]gittool.PathLineDiff{}
	}
	out := make(map[string]shipFileDiff, len(paths))
	for _, p := range paths {
		if d, ok := raw[p]; ok {
			out[p] = shipFileDiff{Added: d.Added, Removed: d.Removed, Binary: d.Binary}
			continue
		}
		if gittool.IsUntrackedPath(repoDir, p) {
			lines := 0
			if st, ok := fileStats[p]; ok && st.Lines >= 0 {
				lines = st.Lines
			}
			out[p] = shipFileDiff{Added: lines, Removed: 0}
		}
	}
	return out
}

func shipFileDiffLabel(d shipFileDiff) string {
	if d.Binary {
		return ui.Gray("binary")
	}
	return ui.Green(fmt.Sprintf("+%d", d.Added)) + ui.Dim(" ") + ui.Red(fmt.Sprintf("-%d", d.Removed))
}

func sumShipFileDiffs(paths []string, diffs map[string]shipFileDiff) (added, removed int) {
	for _, p := range paths {
		d, ok := diffs[p]
		if !ok || d.Binary {
			continue
		}
		added += d.Added
		removed += d.Removed
	}
	return added, removed
}

func shipFileDiffTotalsLabel(added, removed int) string {
	return ui.Green(fmt.Sprintf("+%d", added)) + ui.Dim(" ") + ui.Red(fmt.Sprintf("-%d", removed)) + ui.Gray(" total")
}
