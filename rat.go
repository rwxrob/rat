package rat

import (
	"fmt"
)

// Rule functions are the fundamental building block of any functional
// PEG packrat parser. They can be combined together to form other rules
// and full grammar definitions. They can be returned as closures (see
// Lit, etc.). Since they have no external dependencies they can safely
// be intermixed with Rule functions from other packages and passed as
// first-class anonymous functions.
//
// Note that it is impossible to return nil (as some packrat parser
// designs regularly do). Instead, a failed result with an Err set must
// be returned instead.
//
// Rule functions must advance the End to the farthest possible position
// in the []rune slice before failure occurred providing an Err if
// unable to match the entire rule.  This allows for better recovery and
// specific user-facing error messages. Rune-level granularity is therefore
// preferred when implementing Rule functions.
type Rule func(r []rune, i int) Result

// Result encapsulates the result of a Rule function and contains only
// reference information to the underlying, in-memory buffer.
//
// Buf is but a slice abstraction pointing to the same underlying array
// the same as all other matches.
//
// Beg is the inclusive index of the first rune in Buf that matches.
//
// End is the exclusive index of the end of the match (and Beg of next
// match). End must be advanced to the farthest possible point even if
// the rule fails (and an Err set). This allows recovery attempts from
// that position.
//
// Sub contains sub-matches equivalent to parenthesized patterns of
// a regular expression, the successful child matches under this match.
//
// Err contains any error encountered while parsing.
//
// Note that Beg == End does not indicate failure. End is usually
// greater than Beg, but not necessarily, for example, for look-ahead
// and look-behind rules. End is greater than Beg if a partial match was
// made. Only checking Err can absolutely confirm a rule failure.
type Result struct {
	Buf []rune   // back reference to data (copy of slice only)
	Beg int      // inclusive (B)
	End int      // non-inclusive (E)
	Sub []Result // equivalent to parens in regexp (S)
	Err error    // error, eXpected something else (X)
}

// MarshalJSON fulfills the encoding.JSONMarshaler interface by
// translating Beg to B, End to E, Sub to S, and Err to X as a string.
// Buf is never included. An error is never returned.
func (m Result) MarshalJSON() ([]byte, error) {
	s := `{`
	s += fmt.Sprintf(`"B":%v,"E":%v`, m.Beg, m.End)
	if m.Err != nil {
		s += fmt.Sprintf(`,"X":%q`, m.Err)
	}
	s += `}`
	return []byte(s), nil
}

// String fulfills the fmt.Stringer interface as JSON with "null" for
// any errors.
func (m Result) String() string {
	buf, err := m.MarshalJSON()
	if err != nil {
		return `null`
	}
	return string(buf)
}

// Print is shortcut for fmt.Println(String) but will print "null" if
// the receiver does not exist meaning it can safely be called chained
// to the end of any ParseFunc.
func (m Result) Print() { fmt.Println(m) }

// Lit returns a new Rule function that matches the literal string as
// a []rune slice.
func Lit(s string) Rule {
	return func(r []rune, i int) Result {
		var err error
		var n int
		e := i
		runes := []rune(s)
		for e < len(r) && n < len(runes) {
			if r[e] != runes[n] {
				err = ErrLit{s}
				break
			}
			e++
			n++
		}
		if n < len(runes) {
			err = ErrLit{s}
		}
		return Result{r, i, e, nil, err}
	}
}
