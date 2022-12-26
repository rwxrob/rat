/*
Package rat is Inspired by Bryan Ford's PEG packrat parser paper and is
PEGN (version 2023-01) aware (without dependencies). The packrat
approach does not require a scanner. It is a happy medium between
maintainability and performance. Developers can use rat in
different ways depending on the performance requirements:

* Generate compilable Go code from PEGN or rat structs with rat.Pack(in).Gen(w)
* Dynamically create new packrat grammars from PEGN or rat structs with rat.Pack(in)
* Dynamically parse from ad-hoc grammar with rat.Parse(pegn,in)
* Code grammars directly using shareable rat.Rule libraries

*/
package rat

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// ------------------------------- Pack -------------------------------

// Pack returns a new Grammar from a variety of different input types.
// All arguments are evaluated by type with the most common types first
// in the type switch. For fastest grammar creation pass only predefined
// Rules. For the slowest creation (albeit more convenient) pass
// PEGN(string) types. Aggregate Grammars can be created from others as
// well by importing other Grammar types. All types are parsed and
// cached using the highly concurrent and optimized sync.Map type (see
// Grammar). Once Packed, the returned Grammar can be used to
// dynamically generate compilable Go code (see Gen).
//
// The order of type evaluation and caching is prioritized as
// follows:
//
//     Rule     <as is>
//     string   <inferred Lit>
//     Lit      "foo\n" -> ('foo' LF)
//     Seq      (rule1 rule2)
//     One      (rule1 / rule2)
//     Any      rune{n}
//     Opt      rule?
//     Min0     rule*
//     Min1     rule+
//     Rep      rule{n}
//     Look     (&rule)
//     Not      (!rule)
//     MMax     rule{min,max}
//     Min      rule{min,}
//     Max      rule{0,max}
//     Rng      [a-z]
//     Set      ('a' / 'b' / 'c')
//     NSet     !('a' / 'b' / 'c')
//     To       (rune* &rule)
//     Toi      (rune* rule)
//     Cap      <capture as sub>
//     Tag      <=Foo 'foo' > or Foo <- 'foo'
//     Grammar  <rules imported, overwrite existing>
//     PEGN     <dynamically parsed into rules and cached>
//
// Pointers to these types are not accepted (but the pointer to the underlying Check
// function always is).
func Pack(in ...any) *Grammar {
	g := new(Grammar)
	for _, it := range in {
		switch v := it.(type) {
		case Rule:
			g.Store(v.Text, v)
		case string, Lit:
			g.Lit(string(v))
		case Seq:
			g.Seq(v.Rules...)
		case One:
			g.One(v.Rules...)
		case Any:
			g.Any(v)
		case Opt:
			g.Opt(v.Rules...)
		case Min0:
			g.Min0(v.Rules...)
		case Min1:
			g.Min1(v.Rules...)
		case Rep:
			g.Rep(v.N, v.Rules...)
		case Look:
			g.Look(v.Rules...)
		case Not:
			g.Not(v.Rules...)
		case MMax:
			g.MMax(v.Min, v.Max, v.Rules...)
		case Min:
			g.Min(v.Min, v.Rules...)
		case Max:
			g.Max(v.Max, v.Rules...)
		case Rng:
			g.Rng(v.Beg, v.End)
		case Set:
			g.Set(v)
		case NSet:
			g.NSet(v)
		case To:
			g.To(v.Rules...)
		case Toi:
			g.Toi(v.Rules...)
		case Cap:
			g.Cap(v.Rules...)
		case Tag:
			g.Tag(v.Tag, v.Rules...)
		case *Grammar:
			g.Import(v)
		case Grammar:
			g.Import(&v)
		case PEGN:
			g.PEGN(v)
		}
	}
	return g
}

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
	Text  string
	Check CheckFunc
}

func (r Rule) String() string {
	if r.Text == "" {
		return "Rule"
	}
	return r.Text
}

// ----------------------------- CheckFunc ----------------------------

// CheckFunc evaluates the []rune buffer at a specific position for
// a specific grammar rule and should generally only be used from an
// encapsulating Rule so that it has a Text identifier associated with
// it. One or more Rules may, however, encapsulate the same CheckFunc
// function.
//
// CheckFunc MUST return a Result indicating success or failure by setting
// Err for failure. (Note that this is unlike many packrat designs that
// return nil to indicate rule failure.)
//
// CheckFunc MUST advance the End to the farthest possible position
// in the []rune slice before failure occurred providing an Err if
// unable to match the entire rule. This allows for better recovery and
// specific user-facing error messages and promotes succinct rule
// development.
type CheckFunc func(r []rune, i int) Result

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

// Tagged does a depth first descent into the child sub-results (S)
// adding anything with any of the matching tags to the returned slice.
func (m Result) Tagged(tags ...string) []Results {
	res := []Results{}
	// TODO
	return res
}

// TaggedExp does a depth first descent into the child sub-results (S)
// returning all flattened results with a tag matching the passed
// regular expression.
func (m Result) TaggedExp(exp *regexp.Regexp) []Results {
	res := []Results{}
	// TODO
	return res
}

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
func (rules One) Rule() Rule {

	rule := Rule{
		Text: One(rules).String(),
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

	return rule
}

func (c *Grammar) One(rules ...Rule) Rule {
	rule := One(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// -------------------------------- Seq -------------------------------

// Sequence as in "(foo bar)".
type Seq []Rule

func (rules Seq) String() string {
	str := "(" + rules[0].Text
	for _, v := range rules[1:] {
		str += " " + v.Text
	}
	return str + ")"
}

func (rules Seq) Rule() Rule {
	return Rule{
		Text: Seq(rules).String(),
		Check: func(r []rune, i int) Result {
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
		},
	}
}

func (c *Grammar) Seq(rules ...Rule) Rule {
	rule := Seq(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// -------------------------------- Opt -------------------------------

// Optional as in "(foo bar)?".
type Opt []Rule

// -------------------------------- Rep -------------------------------

// Repeat n times as in "(foo bar){n}".
type Rep struct {
	N     int
	Rules []Rule
}

// ------------------------------ MMax --------------------------------

// Minimum (m) and maximum (n) as in "(foo bar){m,n}".
type MMax struct {
	Min   int
	Max   int
	Rules []Rule
}

// -------------------------------- Min -------------------------------

// Minimum (m) as in "(foo bar){m,}".
type Min struct {
	Min   int
	Rules []Rule
}

// -------------------------------- Max -------------------------------

// Maximum (m) as in "(foo bar){0,m}".
type Max struct {
	Max   int
	Rules []Rule
}

// ------------------------------- Min0 -------------------------------

// Zero minimum as in "(foo bar)*".
type Min0 []Rule

// ------------------------------- Min1 -------------------------------

// One minimum as in "(foo bar)+"
type Min1 []Rule

// -------------------------------- Rng -------------------------------

// Range between runes as in "[a-z]", "[x20-x34]".
type Rng struct {
	Beg rune
	End rune
}

// -------------------------------- Set -------------------------------

// Set from string as in "('a' / 'e' / 'i' / 'o' / 'u')".
type Set string

// ------------------------------- NSet -------------------------------

// Inverted set from script as in "!('a' / 'e' / 'i' / 'o' / 'u')". Same
// as "Not{Set(`aeiou`)}".
type NSet string

// -------------------------------- Cap -------------------------------

// Capture as sub as in "< foo bar >".
type Cap []Rule

// -------------------------------- Tag -------------------------------

// Tagged capture as sub as in "Foo <- 'foo'" or "<=Foo 'foo' >".
type Tag struct {
	Tag   string
	Rules []Rule
}

// ------------------------------- Look -------------------------------

// Positive lookahead as in "&(foo bar)".
type Look []Rule

// -------------------------------- Not -------------------------------

// Negative lookahead as in "!(foo bar)".
type Not []Rule

// -------------------------------- To --------------------------------

// Up to as in "((!(foo bar) rune)*)".
type To []Rule

// -------------------------------- ToI -------------------------------

// Up to inclusive as in "((!(foo bar) rune)*(foo bar))".
type Toi []Rule

// ------------------------------- PEGN -------------------------------

// PEGN communicates to Pack that a string has specific meaning and PEGN
// syntax, otherwise Pack assumes all strings as Lit types.
type PEGN string

func (p PEGN) String() string {
	// TODO

	return ""
}

func (p PEGN) Rule() Rule {
	rule := Rule{
		Text:  p.String(),
		Check: Parse(p),
	}
	return rule
}

// PEGN parses the input (string, []byte, []rune, or PEGN) and returns
// a single top level Rule covering the mini-grammar passed as an
// argument.
func (g *Grammar) PEGN(in any) Rule {
	rule := PEGN(pegn).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ------------------------------- Check ------------------------------

type textual interface {
	PEGN | string | []rune | []byte
}

// Check takes a PEGN rule definition, dynamically compiles it, and
// checks the input returning the Result. This function is slower than
// other alternatives (like a regular expression without compilation),
// but often more convenient. Tagging can be a convenient way to
// retrieve the specific desired result quickly when dealing with
// complex rule definitions. See Result.Tagged.
func Check[T textual, I textual](pegnrule T, in M) Result {
	rule := PEGN(pegnrule).Rule()
	runes := []rune(string(in))
	return rule.Check(runes, 0)
}

// ------------------------------ Import ------------------------------

// Import imports the passed Grammar into the current Grammar
// overwriting any conflicting keys.
func (g Grammar) Import(a *Grammar) {
	a.Range(func(k, v any) bool { g.Store(k, v); return true })
}

// -------------------------- Grammar.String --------------------------

// String fulfills the fmt.Stringer interface by rendering a PEGN
// grammar with incremental Rule names and keys as the Rule definitions.
func (g Grammar) String() string {
	// TODO eventually this should parse each of the rule keys breaking
	// them down into each rule subcomponent and lookup that item from the
	// grammar replacing that location with the invented Rule name. For
	// example, given the following:
	//
	// Rule1 <- 'Foo' LF 'Bar' LF{2}
	// Rule2 <- LF{2}
	//
	// Then as each rule in the Rule definition is examined, LF{2} should
	// be identified.
	//
	// Rule1 <- 'Foo' LF 'Bar' Rule2
	// Rule2 <- LF{2}
	//
	// Rules that have Text that matches RuleName, ClassName, or
	// LitName should be used as is.
	//
	// Rule1    <- 'Foo' LF 'Bar' EndBlock
	// EndBlock <- LF{2}
	return ""
}

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
