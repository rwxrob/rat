package rat

import (
	"fmt"
	"slices"
	"strings"
)

// Result contains the result of an evaluated Rule function along with
// its own shared []rune slice (R).
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
// Note that B == E does not indicate failure. E is usually greater than
// B, but not necessarily, for example, for lookahead rules are
// successful without advancing the position at all. E is also greater
// than B if a partial match was made that still resulted in an error.
// Only checking X can absolutely confirm a rule failure.
//
// X (for "expected") contains any error encountered while parsing.
//
// C (for "children") contains results within this result, sub-matches,
// equivalent to parenthesized patterns of a regular expression.
//
// R (for "result" or "runes" or "buffer") contains the slice with all
// the data in it. Since Go uses the same underlying array for any slice,
// no matter how many times it is referenced, there is no loss of memory
// efficiency even though marshaling a [Result] would produce duplicate
// output. For this reason [MarshalJSON] omits this field from any
// children ([C]) when marshaling.
//
// # Avoid taking reference to Result
//
// A Result is already made up of references so no further dereferencing
// is required. The buffer (R) is a slice and therefore all slices point
// to the same underlying array in memory. And no actual string data is
// saved in any Result. Rather, the beginning and ending positions
// within the buffer data are stored and retrieved when needed with
// methods such as Text().
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
		s += fmt.Sprintf(`"N":%q,`, m.N)
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

// PrintText is short for fmt.Println(m.Text()).
func (m Result) PrintText() { fmt.Println(m.Text()) }

// PrintError is short for fmt.Println(m.X) but adds position.
func (m Result) PrintError() { fmt.Println(m.X) }

// Text returns the text between beginning (B) and ending (E)
// (non-inclusively) It is a shortcut for
// string(res.R[res.B:res.E]).
func (m Result) Text() string {
	return string(m.R[m.B:m.E])
}

// FlatFunc is function that returns a flattened rooted-node tree.
type FlatFunc func(root Result) []Result

// VisitFunc is a first-class function passed one result. Typically
// these functions will enclose variables, contexts, or a channel
// outside of its own scope to be updated for each visit.  Functional
// recursion is usually used, which may present some limitations
// depending on the depth required.
type VisitFunc func(a Result)

// DefaultFlatFunc is the default FlatFunc to use when filtering for results
// using With* methods.
var DefaultFlatFunc = ByDepth

// ByDepth flattens a rooted node tree of Result structs by
// traversing in a synchronous, depth-first, preorder way.
func ByDepth(root Result) []Result {
	results := []Result{root}
	for _, child := range root.C {
		results = append(results, []Result(ByDepth(child))...)
	}
	return results
}

// Walk calls WalkBy(DefaultFlatFunc, root, do).  Use this when the
// order of processing matters more than speed (ASTs, etc.). Also see
// WalkAsync.
func Walk(root Result, do VisitFunc) { WalkBy(DefaultFlatFunc, root, do) }

// WalkBy takes a function to flatten a rooted node tree (FlatFunc),
// creates a flattened slice of Results starting from root Result, and
// then passes each synchronously to the VisitFunc waiting for it to
// complete before doing the next.
// Walk calls WalkBy(DefaultFlatFunc, root, do).  Use this when the
// order of processing matters more than speed (ASTs, etc.). Also see
func WalkBy(flatten FlatFunc, root Result, do VisitFunc) {
	for _, result := range flatten(root) {
		do(result)
	}
}

// MaxGoroutines set the maximum number of goroutines by any method or
// function in this package (WalkByAsync, for example). By default,
// there is no limit (0).
var MaxGoroutines int

// WithName returns all results with any of the passed names. Returns
// zero length slice if no results. As a convenience, multiple names may
// be passed and all matches for each will be grouped together in the
// order provided. See WalkDefault for details on the algorithm used.
func (m Result) WithName(names ...string) []Result {
	results := []Result{}
	Walk(m, func(r Result) {
		if slices.Contains(names, r.N) {
			results = append(results, r)
		}
	})
	return results
}
