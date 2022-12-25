package pegn

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
	for _, r := range lit {
		switch r {
		case '\r':
			s += " CR"
		case '\n':
			s += " LF"
		case '\t':
			s += " TAB"
		}
		if is.String() {
		}
	}
	return s[1:]
}
