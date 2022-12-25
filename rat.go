/*
Package rat is Inspired by Bryan Ford's PEG packrat parser paper. This no-scanner-required parser is a happy medium between maintainability, simplicity, and performance. Developers can dial up performance by using only Checks, or convenience and simplicity with Grammars and Pack structures enabling simple code generation and dynamic PEGN expression parsing.
*/
package rat

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// FuncName is a utility function that makes a best guess at the
// function name of the function passed to it. Returns "func1" (and such)
// for anonymous functions and generally should not be depended upon for
// strict uses where the name must be known. In such cases, creation of
// a proper Rule that associates Text with the rule function is strongly
// preferred.
func FuncName(fn any) string {
	fp := reflect.ValueOf(fn).Pointer()
	long := runtime.FuncForPC(fp).Name()
	parts := strings.Split(long, `.`)
	return parts[len(parts)-1]
}

// Grammar is primarily a cache (embedded sync.Map) for Rules that are
// usually keyed to their Rule.Text. The special Grammar.Pack method
// caches multiple rules represented as different struct types from the
// x (expression) sub-package allowing pure Go representation of
// any PEG construct.
type Grammar struct {
	sync.Map
}

// NewGrammar takes zero or more Rules and returns a grammar
// instantiated from them. Rules are expected to have unique Text field
// values. If conflicting rule.Text values occur the last one wins.
func NewGrammar(rules ...Rule) *Grammar {
	g := new(Grammar)
	for _, rule := range rules {
		g.Store(rule.Text, rule)
	}
	return g
}

type ErrNotExist struct{ This any }

func (e ErrNotExist) Error() string {
	return fmt.Sprintf(``)
}

// Check calls Check on the rule specified by Text key. Sets result
// error (X) to ErrNotExist if rule cannot be found.
func (g *Grammar) Check(ruletext string, r []rune, i int) Result {
	it, found := g.Load(ruletext)
	if !found {
		return Result{B: i, E: i, X: ErrNotExist{ruletext}}
	}
	rule, is := it.(Rule)
	if !is {
		return Result{B: i, E: i, X: ErrNotExist{ruletext}}
	}
	return rule.Check(r, i)
}

// Rule combines one rule function (Check) with some identifying text
// providing rich possibilities for representing grammars textually.
// Rule functions are the fundamental building blocks of any functional
// PEG packrat parser.
//
// Note that two Rules may have the same Check with different Text
// strings and that some applications may change the Text field
// dynamically after instantiation (though this is generally uncommon).
//
// Rules have no external dependencies allowing them to be safely
// combined from multiple packages. For best performance, Rules
// should be created and used from a Grammar with proper caching.
type Rule struct {
	Text string
	Check
}

// Check evaluates the []rune buffer at a specific position for
// a specific grammar rule and should generally only be used from an
// encapsulating Rule so that it has a Text identifier associated with
// it. One or more Rules may, however, encapsulate the same Check
// function.
//
// Check MUST return a Result indicating success or failure by setting
// Err for failure. (Note that this is unlike many packrat designs that
// return nil to indicate rule failure.)
//
// Check MUST advance the End to the farthest possible position
// in the []rune slice before failure occurred providing an Err if
// unable to match the entire rule. This allows for better recovery and
// specific user-facing error messages and promotes succinct rule
// development.
type Check func(r []rune, i int) Result

// Result contains the result of an evaluated Rule function along with
// its own slice (Buf) referring to the same underlying array in memory
// as other rules.
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
	T int      // type (mapped to string name usually) 0=Unknown
	R []rune   // reference data (copy of slice only)
	B int      // beginning (inclusive)
	E int      // ending (non-inclusive)
	S []Result // equivalent to parens in regexp (S)
	X error    // error, eXpected something else (X)
}

// MarshalJSON fulfills the encoding.JSONMarshaler interface by
// translating Beg to B, End to E, Sub to S, and Err to X as a string.
// Buf is never included. An error is never returned.
func (m Result) MarshalJSON() ([]byte, error) {
	s := `{`
	if m.T > 0 {
		s += fmt.Sprintf(`"T":%v`, m.T)
	}
	s += fmt.Sprintf(`"B":%v,"E":%v`, m.B, m.E)
	if m.X != nil {
		s += fmt.Sprintf(`,"X":%q`, m.X)
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

// Print is shortcut for fmt.Println(String).
func (m Result) Print() { fmt.Println(m) }

// PEGNString returns a PEGN grammar converted from a Go string literal.
// PEGN "Strings" are composed of visible ASCII characters excluding all
// white space except space and single quote and are wrapped in single
// quotes. All other valid Go string runes must be represented other ways.
// Popular runes among these are included as their PEGN token names.
//
//     * TAB
//     * CR
//     * LF
//     * CRLF
//
// All others are represented in PEGN hexadecimal notation (ex: 😊 xe056)
// since it requires the least digits and will be used as part of
// a caching key.
//
// Panics if string passed has zero length.
//
func PEGNString(lit string) string {
	var s string
	for _, r := range lit {
		switch r {
		case '\n':
			s += " LF"
			// TODO
		}
	}
	return s[1:]
}

// Lit first checks for an existing rule for the given string in the
// Cache and returns if found. Otherwise, it creates a new Rule that
// matches the literal string as a []rune slice and sets the Rule.Text
// to the string passed.
func (c *Grammar) Lit(s string) Rule {

	rule := Rule{
		Text: PEGNString(s),
	}

	if cached, has := c.Load(rule.Text); has {
		if rule, ok := cached.(Rule); ok {
			return rule
		}
	}

	rule.Check = func(r []rune, i int) Result {
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
		return Result{R: r, B: i, E: e, X: err}
	}

	c.Store(rule.Text, rule)
	return rule
}

type Sequence []Rule

func (s Sequence) String() string {
	str := s[0].Text
	for _, v := range s[1:] {
		str += " " + v.Text
	}
	return str
}

type OneOf []Rule

func (s OneOf) String() string {
	str := s[0].Text
	for _, v := range s[1:] {
		str += " / " + v.Text
	}
	return str
}

// Seq returns a new rule that is the sequential aggregation of all
// rules passed to it stopping on the first to return an Err. Each
// result is added as a Sub match along with the last failed match (if
// any).
func (c *Grammar) Seq(rules ...Rule) Rule {

	rule := Rule{
		Text: Sequence(rules).String(),
	}

	if cached, has := c.Load(rule.Text); has {
		if rule, ok := cached.(Rule); ok {
			return rule
		}
	}

	rule.Check = func(r []rune, i int) Result {
		var err error
		sub := []Result{}
		start := i
		for _, rule := range rules {
			res := rule.Check(r, i)
			i = res.E
			if res.X != nil {
				err = res.X
				break
			}
			sub = append(sub, res)
		}
		if len(sub) == 0 {
			return Result{R: r, B: start, E: i, X: err}
		}
		return Result{R: r, B: start, E: i, S: sub, X: err}
	}

	c.Store(rule.Text, rule)
	return rule
}

// OneOf returns the results of the first rule to successfully match.
func (c *Grammar) OneOf(rules ...Rule) Rule {

	rule := Rule{
		Text: OneOf(rules).String(),
	}

	if cached, has := c.Load(rule.Text); has {
		if rule, ok := cached.(Rule); ok {
			return rule
		}
	}

	rule.Check = func(r []rune, i int) Result {
		for _, rule := range rules {
			res := rule.Check(r, i)
			if res.X == nil {
				return res
			}
		}
		return Result{R: r, B: i, E: i, X: ErrOneOf{rules}}
	}

	c.Store(rule.Text, rule)
	return rule
}
