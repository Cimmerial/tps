package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

type shipFileStat struct {
	Lines int
	Bytes int64
}

func (s shipFileStat) label() string {
	if s.Lines < 0 {
		return "—"
	}
	return fmt.Sprintf("%d lines · %s", s.Lines, formatShipFileSize(s.Bytes))
}

func formatShipFileSize(n int64) string {
	if n < 0 {
		return "—"
	}
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	units := []string{"KB", "MB", "GB", "TB"}
	rem := float64(n)
	for i, u := range units {
		rem /= float64(unit)
		if rem < float64(unit) || i == len(units)-1 {
			if rem >= 10 || rem == float64(int64(rem)) {
				return fmt.Sprintf("%.0f %s", rem, u)
			}
			return fmt.Sprintf("%.1f %s", rem, u)
		}
	}
	return fmt.Sprintf("%d B", n)
}

func readShipFileStat(repoDir, relPath string) shipFileStat {
	abs := filepath.Join(repoDir, filepath.FromSlash(relPath))
	info, err := os.Stat(abs)
	if err != nil {
		return shipFileStat{Lines: -1, Bytes: -1}
	}
	if info.IsDir() {
		return shipFileStat{Lines: -1, Bytes: -1}
	}
	data, err := os.ReadFile(abs)
	if err != nil {
		return shipFileStat{Lines: -1, Bytes: info.Size()}
	}
	return shipFileStat{Lines: countFileLines(data), Bytes: info.Size()}
}

func countFileLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	n := bytes.Count(data, []byte{'\n'})
	if data[len(data)-1] != '\n' {
		n++
	}
	return n
}

func sumShipFileStats(paths []string, stats map[string]shipFileStat) (lines int, bytes int64) {
	for _, p := range paths {
		st := stats[p]
		if st.Lines >= 0 {
			lines += st.Lines
		}
		if st.Bytes >= 0 {
			bytes += st.Bytes
		}
	}
	return lines, bytes
}

func shipFileTotalsLabel(lines int, bytes int64) string {
	return fmt.Sprintf("%d lines total · %s total", lines, formatShipFileSize(bytes))
}

func buildShipFileStats(repoDir string, paths []string) map[string]shipFileStat {
	out := make(map[string]shipFileStat, len(paths))
	for _, p := range paths {
		out[p] = readShipFileStat(repoDir, p)
	}
	return out
}
