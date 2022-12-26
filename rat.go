/*
Package rat is Inspired by Bryan Ford's PEG packrat parser paper and is
PEGN aware (without depending on any external PEGN module). This
no-scanner-required parser is a happy medium between maintainability,
simplicity, and performance. Developers can dial up performance by using
only Checks, or convenience and simplicity with Grammars and Pack
structures enabling simple code generation and dynamic PEGN expression
parsing.

Consider github.com/rwxrob/pegn-go when a full pegn.Grammar is required
(which is an instance of rat.Grammar).

*/
package rat

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// ------------------------------ Grammar -----------------------------

// Grammar embeds sync.Map to cache Rules by their PEGN notation
// Rule.Text identifiers. Most methods add one or more Rule instances to
// this cache. The special Pack method caches multiple rules represented
// by type including string literals, is functions, and struct types
// from the x sub-package that match many of the methods (ex: In).
// This allows for pure Go representation of any PEG grammar.
//
// Adding rules to a Grammar is functionally equivalent to compiling
// a regular expression. That Grammar can then be retrieved by its
// cached definition (top-level Rule.Text) instantly to create parse
// trees of Results against any buffered text data.
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

// CheckString takes a string instead of []rune slice.
func (g *Grammar) CheckString(ruletext string, rstr string, i int) Result {
	return g.Check(ruletext, []rune(rstr), i)
}

// ------------------------------- Rule -------------------------------

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

func (r Rule) String() string {
	if r.Text == "" {
		return "Rule"
	}
	return r.Text
}

// ------------------------------- Check ------------------------------

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

// ------------------------------ Result ------------------------------

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

// ------------------------------ ToPEGN ------------------------------

// ToPEGN returns a PEGN grammar string converted from a Go string
// literal consisting of printable ASCII characters (graph + SP) except
// single quote and wrapped in single quotes. All other runes are PEGN
// hex escaped (ex: 😊 xe056) except for the following popular literals:
//
//     * TAB
//     * CR
//     * LF
//
// ToPEGN panics if string passed has zero length.
func ToPEGN(lit string) string {
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

// -------------------------------- Any -------------------------------

// Specific number (n) of any rune as in "rune{n}".
type Any int

// String fulfills the fmt.Stringer as PEGN "rune{n}".
func (n Any) String() string { return `rune{` + strconv.Itoa(int(n)) + `}` }

func (n Any) Rule() Rule {

	rule := Rule{
		Text: n.String(),
	}

	rule.Check = func(r []rune, i int) Result {
		remaining := len(r[i:])
		if remaining >= int(n) {
			return Result{B: i, E: i + int(n)}
		}
		return Result{B: i, E: i + remaining - int(n), X: ErrExpected{rule}}
	}

	return rule
}

func (c *Grammar) Any(n int) Rule {
	rule := Any(n).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// -------------------------------- Lit -------------------------------

// Lit is a shortcut method of considering a string as a efficient
// method of storing what would otherwise be a Seq of PEGN rules.
type Lit string

// String fulfills the fmt.Stringer interface by delegating ToPEGN.
func (s Lit) String() string { return ToPEGN(string(s)) }

// Lit dynamically creates a Rule for the given string and sets the
// Rule.Text to PEGN notation from ToPEGN for the given input
// argument.
func (s Lit) Rule() Rule {
	rule := Rule{
		Text: Lit(s).String(),
	}
	rule.Check = func(r []rune, i int) Result {
		var err error
		var n int
		e := i
		runes := []rune(s)
		for e < len(r) && n < len(runes) {
			if r[e] != runes[n] {
				err = ErrExpected{r[e]}
				break
			}
			e++
			n++
		}
		if n < len(runes) {
			err = ErrExpected{runes[n]}
		}
		return Result{R: r, B: i, E: e, X: err}
	}
	return rule
}

func (c *Grammar) Lit(s string) Rule {
	rule := Lit(s).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// -------------------------------- One -------------------------------

// One of rules matches as in "(foo / bar)".
type One []Rule

func (s One) String() string {
	str := "(" + s[0].Text
	for _, v := range s[1:] {
		str += " / " + v.Text
	}
	return str + ")"
}

// One returns the results of the first rule to successfully match.
// Equivalent to PEGN "(foo / bar)".
func (c *Grammar) One(rules ...Rule) Rule {

	rule := Rule{
		Text: One(rules).String(),
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

		return Result{R: r, B: i, E: i, X: ErrExpected{rule}}
	}

	c.Store(rule.Text, rule)
	return rule
}

// -------------------------------- Seq -------------------------------

// Sequence as in "(foo bar)".
type Seq []Rule

func (s Seq) String() string {
	str := "(" + s[0].Text
	for _, v := range s[1:] {
		str += " " + v.Text
	}
	return str + ")"
}

// Seq returns a new rule that is the sequential aggregation of all
// rules passed to it stopping on the first to return an Err. Each
// result is added as a Sub match along with the last failed match (if
// any). Equivalent to PEGN "(foo bar)".
func (c *Grammar) Seq(rules ...Rule) Rule {

	rule := Rule{
		Text: Seq(rules).String(),
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

// -------------------------------- Opt -------------------------------

// Optional as in "(foo bar)?".
type Opt []Rule

// Finite number (n) as in "(foo bar){n}".
type N struct {
	N     int
	Rules []Rule
}

// Minimum (m) and maximum (n) as in "(foo bar){m,n}".
type MinMax struct {
	Min   int
	Max   int
	Rules []Rule
}

// Minimum (m) as in "(foo bar){m,}".
type Min struct {
	Min   int
	Rules []Rule
}

// Maximum (m) as in "(foo bar){0,m}".
type Max struct {
	Max   int
	Rules []Rule
}

// Zero minimum as in "(foo bar)*".
type Min0 []Rule

// One minimum as in "(foo bar)+".
type Min1 []Rule

// Range between runes as in "[a-z]".
type Rng struct {
	Beg rune
	End rune
}

// Capture as sub as in "< foo bar >".
type Cap []Rule

// Tagged capture as sub as in "<=tag foo bar >".
type Tag struct {
	Tag   string
	Rules []Rule
}

// Positive lookahead as in "&(foo bar)".
type Pos []Rule

// Negative lookahead as in "!(foo bar)".
type Neg []Rule

// Up to as in "((!(foo bar) rune)*)".
type To []Rule

// Up to inclusive as in "((!(foo bar) rune)*(foo bar))".
type ToI []Rule

// ------------------------------ Errors ------------------------------

type ErrExpected struct {
	This any
}

func (e ErrExpected) Error() string {
	switch v := e.This.(type) {
	case rune:
		e.This = ToPEGN(string(v))
	}
	return fmt.Sprintf(_ErrExpected, e.This)
}

type ErrNotExist struct{ This any }

func (e ErrNotExist) Error() string {
	return fmt.Sprintf(_ErrNotExist, e.This)
}
