package cli

import "strings"

// splitInputFields splits a line like the shell: whitespace separators, double/single quotes keep words together.
func splitInputFields(line string) []string {
	s := strings.TrimSpace(line)
	if s == "" {
		return nil
	}
	var out []string
	var cur strings.Builder
	inQuote := false
	var quote byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if inQuote {
			if c == quote {
				inQuote = false
				continue
			}
			cur.WriteByte(c)
			continue
		}
		switch c {
		case '"', '\'':
			inQuote = true
			quote = c
		case ' ', '\t':
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteByte(c)
		}
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}
