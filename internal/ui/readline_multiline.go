package ui

import "strings"

// shiftEnterInsertsNewline reports whether an escape sequence is Shift+Enter (or a close variant).
func shiftEnterInsertsNewline(seq string, multiline bool) bool {
	if !multiline || seq == "" {
		return false
	}
	switch seq {
	case "[13;2~", "[13;2u", "[27;2;13~", "\x1b\r", "\x1b\n":
		return true
	}
	// xterm modifyOtherKeys / kitty keyboard protocol variants.
	if strings.HasPrefix(seq, "[13;") {
		if strings.HasSuffix(seq, ";2~") || strings.HasSuffix(seq, ";2u") {
			return true
		}
	}
	if strings.HasPrefix(seq, "[27;2;13") {
		return true
	}
	return false
}
