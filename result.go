package rat

import "fmt"

// Result contains the result of an evaluated Rule function along with
// its own []rune slice (R) (which refers to the same underlying array
// in memory as other rules).
//
// N (for "name") is an integer mapped to an enumeration of string
// names.
//
// B (for "beginning") is the inclusive beginning position index of the
// first rune in Buf that matches.
//
// E (for "ending") is the exclusive ending position index of the end of
// the match (and Beg of next match). End must be advanced to the
// farthest possible point even if the rule fails (and an Err set). This
// allows recovery attempts from that position.
//
// S (for "sub-match") contains sub-matches equivalent to parenthesized
// patterns of a regular expression, the successful child matches under
// this match.
//
// X contains any error encountered while parsing.
//
// Note that B == E does not indicate failure. E is usually greater than
// B, but not necessarily, for example, for lookahead rules are
// successful without advancing the position at all. E is also greater
// than B if a partial match was made that still resulted in an error.
// Only checking X can absolutely confirm a rule failure.
type Result struct {
	T int      // integer rule type corresponding to rule.ID
	R []rune   // reference data (underlying slice array shared)
	B int      // beginning (inclusive)
	E int      // ending (non-inclusive)
	S []Result // sub-match children, equivalent to parens in regexp
	X error    // error, eXpected something else
}

// MarshalJSON fulfills the encoding.JSONMarshaler interface by
// translating Beg to B, End to E, Sub to S, and Err to X as a string.
// Buf is never included. An error is never returned.
func (m Result) MarshalJSON() ([]byte, error) {
	s := "{"

	if m.T > 0 {
		s += fmt.Sprintf(`"T":%v,`, m.T)
	}

	s += fmt.Sprintf(`"B":%v,"E":%v`, m.B, m.E)

	if m.X != nil {
		s += fmt.Sprintf(`,"X":%q`, m.X)
	}

	s += "}"

	return []byte(s), nil
}

// String fulfills the fmt.Stringer interface as JSON with "null" for
// any errors.
func (m Result) String() string {
	buf, err := m.MarshalJSON()
	if err != nil {
		return "null"
	}
	return string(buf)
}

// Print is shortcut for fmt.Println(String).
func (m Result) Print() { fmt.Println(m) }

// Named does a depth first descent into the sub-match child results (S)
// adding any result with a name (Result.N>0) that matches.
func (m Result) Named(names ...int) []Result {
	res := []Result{}
	// TODO
	return res
}

// NamedString is the same as Named but uses the full string name.
func (m Result) NamedString(names ...string) []Result {
	res := []Result{}
	// TODO
	return res
}
