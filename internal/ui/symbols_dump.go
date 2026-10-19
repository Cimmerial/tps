package ui

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/Cimmerial/tps/internal/ui/logos"
)

// SymbolSection is one block in the sym info dump.
type SymbolSection struct {
	Title   string
	Entries []SymbolEntry
}

// ANSICodeEntry documents terminal styling sequences used by ui.paint.
type ANSICodeEntry struct {
	ID    string
	Seq   string
	Usage string
}

// ANSICodes lists escape sequences (not Unicode glyphs).
func ANSICodes() []ANSICodeEntry {
	return []ANSICodeEntry{
		{ID: "reset", Seq: "\\033[0m", Usage: "end styling"},
		{ID: "bold", Seq: "\\033[1m", Usage: "Bold(), headers"},
		{ID: "dim", Seq: "\\033[2m", Usage: "Dim(), hints, secondary text"},
		{ID: "red", Seq: "\\033[31m", Usage: "Red(), cross mark"},
		{ID: "green", Seq: "\\033[32m", Usage: "Green(), check, success"},
		{ID: "yellow", Seq: "\\033[33m", Usage: "Yellow(), warnings, project name"},
		{ID: "blue", Seq: "\\033[34m", Usage: "Blue()"},
		{ID: "magenta", Seq: "\\033[35m", Usage: "Magenta()"},
		{ID: "cyan", Seq: "\\033[36m", Usage: "Cyan(); default accent / title card"},
		{ID: "faint", Seq: "\\033[37m", Usage: "Faint(), current-value previews"},
		{ID: "clear_screen", Seq: "\\033[2J\\033[H", Usage: "ClearScreen()"},
		{ID: "clear_scrollback", Seq: "\\033[2J\\033[3J\\033[H", Usage: "ClearTerminal() (Cmd+K style)"},
		{ID: "erase_line", Seq: "\\r\\033[K", Usage: "spinner clear line on stderr"},
	}
}

// SymbolDump builds the full sym catalog: named glyphs, spinners, ANSI, logos, and repo scan.
func SymbolDump(repoRoot string, filter string) []SymbolSection {
	filter = strings.TrimSpace(strings.ToLower(filter))
	sections := []SymbolSection{
		{Title: "named ui glyphs", Entries: filterNamedCatalog(filter)},
		{Title: "braille spinner frames", Entries: filterEntries(brailleFrameEntries(), filter)},
		{Title: "ansi terminal codes (not unicode)", Entries: filterANSISection(filter)},
	}
	if repoRoot != "" {
		if logo := logoGlyphSection(repoRoot, filter); len(logo.Entries) > 0 {
			sections = append(sections, logo)
		}
		if scanned := scannedGlyphSection(repoRoot, filter); len(scanned.Entries) > 0 {
			sections = append(sections, scanned)
		}
	}
	out := make([]SymbolSection, 0, len(sections))
	for _, s := range sections {
		if len(s.Entries) > 0 {
			out = append(out, s)
		}
	}
	return out
}

func brailleFrameEntries() []SymbolEntry {
	out := make([]SymbolEntry, 0, len(BrailleSpinnerFrames))
	for i, g := range BrailleSpinnerFrames {
		out = append(out, SymbolEntry{
			ID:    fmt.Sprintf("braille-%d", i),
			Glyph: g,
			Usage: fmt.Sprintf("spinner / busy prompt frame %d of %d", i+1, len(BrailleSpinnerFrames)),
		})
	}
	return out
}

func filterANSISection(filter string) []SymbolEntry {
	var out []SymbolEntry
	for _, a := range ANSICodes() {
		e := SymbolEntry{ID: a.ID, Glyph: a.Seq, Usage: a.Usage}
		if filter == "" || entryMatchesFilter(e, filter) || strings.Contains(strings.ToLower(a.Seq), filter) {
			out = append(out, e)
		}
	}
	return out
}

func filterNamedCatalog(filter string) []SymbolEntry {
	all := SymbolCatalog()
	if filter == "" {
		return all
	}
	return filterEntries(all, filter)
}

func filterEntries(in []SymbolEntry, filter string) []SymbolEntry {
	if filter == "" {
		return in
	}
	var out []SymbolEntry
	for _, e := range in {
		if entryMatchesFilter(e, filter) {
			out = append(out, e)
		}
	}
	return out
}

func entryMatchesFilter(e SymbolEntry, filter string) bool {
	if strings.Contains(strings.ToLower(e.ID), filter) {
		return true
	}
	if strings.Contains(strings.ToLower(e.Usage), filter) {
		return true
	}
	if strings.Contains(strings.ToLower(e.Glyph), filter) {
		return true
	}
	if strings.Contains(strings.ToLower(glyphCategoryLabel(e.Glyph)), filter) {
		return true
	}
	return false
}

type glyphRecord struct {
	rune    rune
	sources map[string]bool
}

func logoGlyphSection(repoRoot string, filter string) SymbolSection {
	records := map[rune]*glyphRecord{}
	logos.WalkArtRunes(func(r rune, src string) {
		if !interestingRune(r) {
			return
		}
		rec := records[r]
		if rec == nil {
			rec = &glyphRecord{rune: r, sources: map[string]bool{}}
			records[r] = rec
		}
		rec.sources[src] = true
	})
	return sectionFromRecords("title-card logo art (embedded presets)", records, filter)
}

func scannedGlyphSection(repoRoot string, filter string) SymbolSection {
	records, err := scanRepoGlyphs(repoRoot)
	if err != nil {
		return SymbolSection{
			Title: "repo scan (error)",
			Entries: []SymbolEntry{{
				ID:    "scan",
				Glyph: "!",
				Usage: err.Error(),
			}},
		}
	}
	return sectionFromRecords("all non-ascii in tps source (internal/, cmd/, instance/, *.md)", records, filter)
}

func sectionFromRecords(title string, records map[rune]*glyphRecord, filter string) SymbolSection {
	keys := make([]rune, 0, len(records))
	for r := range records {
		keys = append(keys, r)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	var entries []SymbolEntry
	for _, r := range keys {
		rec := records[r]
		g := string(r)
		srcList := sortedKeys(rec.sources)
		usage := fmt.Sprintf("%s · seen in %s", glyphCategoryLabel(g), strings.Join(srcList, ", "))
		if len(srcList) > 4 {
			usage = fmt.Sprintf("%s · seen in %s (+%d more)", glyphCategoryLabel(g), strings.Join(srcList[:4], ", "), len(srcList)-4)
		}
		e := SymbolEntry{ID: glyphID(r), Glyph: g, Usage: usage}
		if filter != "" && !entryMatchesFilter(e, filter) {
			continue
		}
		entries = append(entries, e)
	}
	return SymbolSection{Title: title, Entries: entries}
}

func glyphID(r rune) string {
	return fmt.Sprintf("U+%04X", r)
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func interestingRune(r rune) bool {
	if r == '\n' || r == '\r' || r == '\t' {
		return false
	}
	if r < 128 {
		return false
	}
	return true
}

func glyphCategoryLabel(glyph string) string {
	runes := []rune(glyph)
	if len(runes) == 0 {
		return "unknown"
	}
	return runeCategory(runes[0])
}

func runeCategory(r rune) string {
	switch {
	case r >= 0x2800 && r <= 0x28FF:
		return "Braille"
	case r >= 0x2580 && r <= 0x259F:
		return "Block elements (shade bars)"
	case r >= 0x2500 && r <= 0x257F:
		return "Box drawing"
	case unicode.Is(unicode.Dash, r):
		return "Dash punctuation"
	case r == '·':
		return "Middle dot (separator)"
	case r == '…':
		return "Ellipsis"
	case r == '→':
		return "Arrow"
	case r == '⌘':
		return "Command key"
	case r == '✓' || r == '✗':
		return "Status mark"
	case r >= 0x2550:
		return "Box drawing / symbols"
	default:
		return "Other (" + unicodeCategory(r) + ")"
	}
}

func unicodeCategory(r rune) string {
	if unicode.IsLetter(r) {
		return "letter"
	}
	if unicode.IsDigit(r) {
		return "digit"
	}
	if unicode.IsPunct(r) {
		return "punctuation"
	}
	if unicode.IsSymbol(r) {
		return "symbol"
	}
	return "other"
}

func scanRepoGlyphs(repoRoot string) (map[rune]*glyphRecord, error) {
	records := map[rune]*glyphRecord{}
	add := func(r rune, src string) {
		if !interestingRune(r) {
			return
		}
		rec := records[r]
		if rec == nil {
			rec = &glyphRecord{rune: r, sources: map[string]bool{}}
			records[r] = rec
		}
		rec.sources[src] = true
	}
	for _, rel := range []string{"internal", "cmd", "instance"} {
		root := filepath.Join(repoRoot, rel)
		if err := walkTreeGlyphs(repoRoot, root, add); err != nil {
			return nil, err
		}
	}
	for _, name := range []string{"README.md", "DESIGN_DOC.md"} {
		path := filepath.Join(repoRoot, name)
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		for _, r := range string(data) {
			add(r, name)
		}
	}
	return records, nil
}

func walkTreeGlyphs(repoRoot, root string, add func(rune, string)) error {
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !info.IsDir() {
		return nil
	}
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == "vendor" || name == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".go" && ext != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		src, err := filepath.Rel(repoRoot, path)
		if err != nil {
			src = path
		}
		for _, r := range string(data) {
			add(r, src)
		}
		return nil
	})
}
