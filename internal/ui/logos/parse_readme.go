package logos

import (
	"strings"
	"unicode"
)

func parseREADMELogos(text string) []Preset {
	lines := strings.Split(text, "\n")
	start := -1
	for i, ln := range lines {
		if strings.TrimSpace(ln) == "default" && i > 0 && strings.TrimSpace(lines[i-1]) == "" {
			start = i
			break
		}
	}
	if start < 0 {
		return nil
	}
	var out []Preset
	i := start
	for i < len(lines) {
		name := strings.TrimSpace(lines[i])
		if name == "" {
			i++
			continue
		}
		i++
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
		var art []string
		for i < len(lines) {
			if strings.TrimSpace(lines[i]) == "" {
				j := i + 1
				for j < len(lines) && strings.TrimSpace(lines[j]) == "" {
					j++
				}
				if j >= len(lines) {
					break
				}
				if isREADMELogoName(lines[j]) {
					break
				}
				i++
				continue
			}
			art = append(art, lines[i])
			i++
		}
		out = append(out, Preset{
			ID:    slugLogoName(name),
			Label: name,
			Art:   strings.TrimRight(strings.Join(art, "\n"), "\n"),
		})
		for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
			i++
		}
	}
	return out
}

func isREADMELogoName(line string) bool {
	if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
		return false
	}
	t := strings.TrimSpace(line)
	if t == "" || strings.HasPrefix(t, "#") {
		return false
	}
	if len(t) > 50 {
		return false
	}
	sym := 0
	for _, r := range t {
		if strings.ContainsRune("|\\/[]{}$░▒▓█▄▖▗▘▝▞▚`~^", r) {
			sym++
		}
	}
	if sym >= 3 {
		return false
	}
	letters := 0
	for _, r := range t {
		if unicode.IsLetter(r) || r == ' ' || r == '\'' || r == '.' {
			letters++
		}
	}
	return letters >= len(t)/2
}

func slugLogoName(name string) string {
	var b strings.Builder
	lastUnderscore := false
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastUnderscore = false
			continue
		}
		if !lastUnderscore {
			b.WriteByte('_')
			lastUnderscore = true
		}
	}
	return strings.Trim(b.String(), "_")
}
