package rat

import (
	"fmt"
	"strings"
)

// Result contains the result of an evaluated Rule function along with
// its own []rune slice (R) (which refers to the same underlying array
// in memory as other rules).
//
// T (for "type") is an integer mapped to the rule that was used to
// produce this result, usually associated with a longer string name.
//
// B (for "beginning") is the inclusive beginning position index of the
// first rune in Buf that matches.
//
// E (for "ending") is the exclusive ending position index of the end of
// the match (and Beg of next match). End must be advanced to the
// farthest possible point even if the rule fails (and an Err set). This
// allows recovery attempts from that position.
//
// C (for "children") contains results within this result, sub-matches,
// equivalent to parenthesized patterns of a regular expression.
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
	B int      // beginning (inclusive)
	E int      // ending (non-inclusive)
	X error    // error, eXpected something else
	C []Result // children, results within this result
	R []rune   // reference data (underlying slice array shared)
}

// MarshalJSON fulfills the encoding.JSONMarshaler interface. The begin
// (B), end (E) are always included. The type (T), buffer (R), error (X)
// and child sub-matches (C) are only included if not empty. Child
// sub-matches omit the buffer (R). The order of fields is guaranteed
// not to change.  Output is always a single line. There is no
// dependency on the reflect package. The buffer (R) is rendered as
// a quoted string (%q) with no further escaping (unlike built-in Go
// JSON marshaling which escapes things unnecessarily producing
// unreadable output). An error is never returned.
func (m Result) MarshalJSON() ([]byte, error) {

	s := "{"

	if m.T > 0 {
		s += fmt.Sprintf(`"T":%v,`, m.T)
	}

	s += fmt.Sprintf(`"B":%v,"E":%v`, m.B, m.E)

	if m.X != nil {
		s += fmt.Sprintf(`,"X":%q`, m.X)
	}

	if len(m.C) > 0 {
		results := []string{}
		for _, c := range m.C {
			results = append(results, Result{c.T, c.B, c.E, c.X, c.C, nil}.String())
		}
		s += `,"C":[` + strings.Join(results, ",") + `]`
	}

	if m.R != nil {
		s += fmt.Sprintf(`,"R":%q`, string(m.R))
	}

	s += "}"

	return []byte(s), nil
}

// String fulfills the fmt.Stringer interface as JSON by calling
// MarshalJSON. If JSON marshaling fails for any reason a "null" string
// is returned.
func (m Result) String() string {
	buf, err := m.MarshalJSON()
	if err != nil {
		return "null"
	}
	return string(buf)
}

// Print is shortcut for fmt.Println(String).
func (m Result) Print() { fmt.Println(m) }

// Text returns the text between beginning (B) and ending (E)
// (non-inclusively) It is a shortcut for res.R[res.B:res.E].
func (m Result) Text() { return m.R[m.B:m.E] }

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
