package pegn

import (
	"fmt"
	"strings"
)

// FromString returns a PEGN grammar converted from a Go string literal.
// PEGN "Strings" are composed of visible ASCII characters excluding all
// white space except space and single quote and are wrapped in single
// quotes. All other valid Go string runes must be represented other ways.
// Popular runes among these are included as their PEGN token names.
//
//     * TAB
//     * CR
//     * LF
//
// All others are represented in PEGN hexadecimal notation (ex: 😊 xe056)
// since it requires the least digits and will be used as part of
// a caching key.
//
// Panics if string passed has zero length.
//
func FromString(lit string) string {
	var s string
	var instr bool
	for _, r := range lit {

		if 'a' <= r && r <= 'z' {
			if !instr {
				s += " '" + string(r)
				instr = true
				continue
			}
			s += string(r)
			continue
		}

		if instr {
			s += "'"
			instr = false
		}

		// common tokens
		switch r {
		case '\r':
			s += " CR"
			continue
		case '\n':
			s += " LF"
			continue
		case '\t':
			s += " TAB"
			continue
		case '\'':
			s += " SQ"
			continue
		}

		// escaped
		s += " x" + fmt.Sprintf("%x", r)

	}

	if instr {
		s += "'"
	}

	if strings.Index(s[1:], " ") > 0 {
		return "(" + s[1:] + ")"
	}
	return s[1:]
}
