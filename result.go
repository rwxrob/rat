package rat

import (
	"fmt"
	"strings"
)

// Result contains the result of an evaluated Rule function along with
// its own []rune slice (R).
//
// N (for "Name") is a string name for a result mapped to the x.Name
// rule. This makes for easy reading and walking of results trees since
// the string name is included in the output JSON (see String).
// Normally, mixing x.ID and x.Name rules is avoided.
//
// I (for "ID") is an integer mapped to the x.ID rule. Integer IDs are
// preferable to names (N) in cases where the use of names would
// increase the parse tree output JSON (see String) beyond acceptable
// levels since integer identifiers rarely take more than 2 runes each.
// Normally, mixing x.ID and x.Name rules is avoided.
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
//
// Avoid taking reference to Result
//
// A Result is already made up of references so no further dereferencing
// is required. The buffer (R) is a slice and therefore all slices point
// to the same underlying array in memory. And no actual string data is
// saved in any Result. Rather, the beginning and ending positions
// within the buffer data are stored and retrieved when needed with
// methods such as Text().
//
type Result struct {
	N string   // string name (x.Name)
	I int      // integer identifier (x.ID)
	B int      // beginning (inclusive)
	E int      // ending (non-inclusive)
	X error    // error, eXpected something else
	C []Result // children, results within this result
	R []rune   // reference data (underlying slice array shared)
}

// MarshalJSON fulfills the encoding.JSONMarshaler interface. The begin
// (B), end (E) are always included. The name (N), id (I), buffer (R),
// error (X) and child sub-matches (C) are only included if not empty.
// Child sub-matches omit the buffer (R). The order of fields is
// guaranteed not to change.  Output is always a single line. There is
// no dependency on the reflect package. The buffer (R) is rendered as
// a quoted string (%q) with no further escaping (unlike built-in Go
// JSON marshaling which escapes things unnecessarily producing
// unreadable output). The buffer (R) is never included for children
// (which is the same). An error is never returned.
func (m Result) MarshalJSON() ([]byte, error) {

	s := "{"

	if m.N != "" {
		s += fmt.Sprintf(`"N":%v,`, m.N)
	}

	if m.I > 0 {
		s += fmt.Sprintf(`"I":%v,`, m.I)
	}

	s += fmt.Sprintf(`"B":%v,"E":%v`, m.B, m.E)

	if m.X != nil {
		s += fmt.Sprintf(`,"X":%q`, m.X)
	}

	if len(m.C) > 0 {
		results := []string{}
		for _, c := range m.C {
			results = append(results, Result{c.N, c.I, c.B, c.E, c.X, c.C, nil}.String())
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
// (non-inclusively) It is a shortcut for
// string(res.R[res.B:res.E]).
func (m Result) Text() string { return string(m.R[m.B:m.E]) }

/*
// WalkLevels will pass a tree of results (starting with itself) to the
// given function traversing in a synchronous, breadth-first, leveler
// way. The function passed may be a closure containing variables,
// contexts, or a channel outside of its own scope to be updated for
// each visit. This method uses functional recursion which may have some
// limitations depending on the depth of node trees required.
func (m Result) WalkLevels(do func(n Result)) {
		list := qstack.New[*Node[T]]()
		list.Unshift(n)
		for list.Len > 0 {
			cur := list.Shift()
			list.Push(cur.Nodes()...)
			do(cur)
		}
}
*/

// Named returns all results with any of the passed names. Returns zero length slice if no results.
// TODO set tree walking algorithm
func (m Result) Named(names ...string) []Result {
	res := []Result{}

	// look up the current name

	// iterate through children, depth first

	return res
}

// First returns first hit from Named or nil if no matches.
func (m Result) First(named ...string) *Result {
	results := m.Named(named...)
	if len(results) == 0 {
		return nil
	}
	return &results[0]
}
