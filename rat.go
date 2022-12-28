/*
Package rat is Inspired by Bryan Ford's PEG packrat parser paper and is
PEGN (version 2023-01) aware (without dependencies). The packrat
approach does not require a scanner. It is a happy medium between
maintainability and performance. Developers can use rat in
different ways depending on the performance requirements:


* rat.Gen(p,in)       - generate compilable Go code from PEGN or rat structs
* rat.Pack(in)        - dynamically create new grammars from PEGN or rat structs
* rat.Check(pegn,in)  - dynamically check input against ad-hoc grammars
* rat.Grammar         - build up performant grammars with Go directly
* rat.Rule            - create reusable, composable rule libraries

*/
package rat

import (
	"fmt"
	"io"
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
//     N        <:foo 'foo' > / Foo <- 'foo' / Foo <- ws* < 'foo' >  ws*
//     string   <inferred Lit>
//     Lit      "foo\n" -> ('foo' LF)
//     Seq      (rule1 rule2)
//     One      (rule1 / rule2)
//     Any      rune{n}
//     Opt      rule?
//     Min0     rule*
//     Min1     rule+
//     Rep      rule{n}
//     Is       &rule
//     Not      !rule
//     MMax     rule{min,max}
//     Min      rule{min,}
//     Max      rule{0,max}
//     Rng      [a-z] / [A-Z] / [x76-x34] / [u0000-u10FFFF] / [o20-o37]
//     Set      ('a' / 'b' / 'c') <invert with Not{Set}>
//     To       ...rule
//     Toi      ..rule
//     Grammar  <rules imported, overwrite existing>
//     PEGN     <dynamically parsed into rules and cached>
//
// Pointers to these types are not accepted (but the pointer to the
// underlying Check function always is).
//
// Errors during packing are pushed to Grammar.Errors (the Grammar
// itself can be used as an error as well).
//
func Pack(in ...any) *Grammar {
	g := new(Grammar)
	for _, it := range in {
		switch v := it.(type) {
		case Rule:
			g.Store(v.Text, v)
		case N:
			g.RuleNamed(v.Name, v.Rules...)
		case string:
			g.RuleLiteral(string(v))
		case Lit:
			g.RuleLiteral(string(v))
		case Seq:
			g.RuleSequence(v...)
		case One:
			g.RuleOneOf(v...)
		case Any:
			g.RuleAnyN(int(v))
		case Opt:
			g.RuleOptional(v...)
		case Min0:
			g.RuleMin0(v...)
		case Min1:
			g.RuleMin1(v...)
		case Rep:
			g.RuleRepeat(v.N, v.Rules...)
		case Is:
			g.RuleIs(v...)
		case Not:
			g.RuleNot(v...)
		case MMax:
			g.RuleMinMax(v.Min, v.Max, v.Rules...)
		case Min:
			g.RuleMin(v.Min, v.Rules...)
		case Max:
			g.RuleMax(v.Max, v.Rules...)
		case Rng:
			g.RuleRange(v.Beg, v.End)
		case Set:
			g.RuleSet(string(v))
		case To:
			g.RuleTo(v...)
		case Toi:
			g.RuleToInc(v...)
		case *Grammar:
			g.Import(v)
		case Grammar:
			g.Import(&v)
		case PEGN:
			g.PEGN(v.Name, v.Text)
		default:
			g.ErrPush(ErrBadType{v})
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
	Names    []string       // mapped to Result.T
	NamesMap map[string]int // lookup index
	Errors   []error
}

// ErrPush pushes an error onto the Errors slice.
func (g *Grammar) ErrPush(err error) {
	g.Errors = append(g.Errors, err)
}

// ErrPop pops last error off of the Errors slice returning nil if
// nothing to pop.
func (g *Grammar) ErrPop() error {
	if len(g.Errors) == 0 {
		return nil
	}
	err := g.Errors[len(g.Errors)-1]
	g.Errors = g.Errors[:len(g.Errors)-2]
	return err
}

// Check calls Check on the rule specified by key looking it up from the
// internal sync.Map cache. Sets result error (X) to ErrNotExist if rule
// cannot be found.
func (g *Grammar) Check(key string, r []rune, i int) Result {

	it, found := g.Load(key)
	if !found {
		return Result{B: i, E: i, X: ErrNotExist{key}}
	}
	rule, is := it.(Rule)
	if !is {
		return Result{B: i, E: i, X: ErrNotExist{key}}
	}
	return rule.Check(r, i)
}

// CheckString takes a string instead of []rune slice.
func (g *Grammar) CheckString(key string, rstr string, i int) Result {
	return g.Check(key, []rune(rstr), i)
}

// ------------------------------- Rule -------------------------------

// Rule combines one rule function (Check) with some identifying text
// providing rich possibilities for representing grammars textually.
// Rule functions are the fundamental building blocks of any functional
// PEG packrat parser.
//
// Usually the Text will be a unique identifier for a given Grammar. In
// PEGN this would be the left-hand side of a rule definition (ex: Foo
// <- 'foo' would be Foo). The Text can be anything, however, so long as
// it uniquely identifies the rule. Often the Text will be the
// grammar-specific syntax representing the Rule since many do not have
// individual identifiers (ex: 'foo', rune{4}).
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

// CheckFunc examines the []rune buffer at a specific position for
// a specific grammar rule and should generally only be used from an
// encapsulating Rule so that it has a Text identifier associated with
// it. One or more Rules may, however, encapsulate the same CheckFunc
// function.
//
// CheckFunc MUST return a Result indicating success or failure by
// setting Err for failure. (Note that this is unlike many packrat
// designs that return nil to indicate rule failure.)
//
// CheckFunc MUST set Err if unable to match the entire rule and MUST
// advance to the End to the farthest possible position in the []rune
// slice before failure occurred. This allows for better recovery and
// specific user-facing error messages while promoting succinct rule
// development.
type CheckFunc func(r []rune, i int) Result

// ------------------------------ Result ------------------------------

// Result contains the result of an evaluated Rule function along with
// its own []rune slice (R) (which refers to the same underlying array
// in memory as other rules).
//
// T (for "type") is an integer mapped to an enumeration of string
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
	N int      // integer id in enumeration of name strings
	R []rune   // reference data (copy of slice only)
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

	if m.N > 0 {
		s += fmt.Sprintf(`"N":%v,`, m.N)
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

// NamedExp is the same as Named but uses a regular expression match
// against the string version of the name.
func (m Result) NamedExp(exp *regexp.Regexp) []Result {
	res := []Result{}
	// TODO
	return res
}

// ------------------------------ ToPEGN ------------------------------

// ToPEGN returns a PEGN version of the argument passed to it.
//
//     * runes are converted to single quoted strings
//     * integers are converted to single quoted string version
//     * strings are converted phrases of literals
//
// literal or rune consisting of printable ASCII characters (graph + SP)
// except single quote and wrapped in single quotes. All other runes are
// PEGN hex escaped (ex: 😊 xe056) except for the following popular
// literals:
//
//     * TAB
//     * CR
//     * LF
//
// ToPEGN returns empty string if passed lit string as zero length.
func ToPEGN(it any) string {
	var lit string

	switch v := it.(type) {

	case string:
		lit = v

	case []rune:
		lit = string(v)

	case rune:
		if 'a' <= v && v <= 'z' || 'A' <= v && v <= 'Z' || '0' <= v && v <= '9' {
			return string(v)
		}
		return fmt.Sprintf(`x%x`, v)

	case int:
		return `'` + strconv.Itoa(v) + `'`

	case int8:
		return `'` + strconv.Itoa(int(v)) + `'`

	case int64:
		return `'` + strconv.Itoa(int(v)) + `'`

	case uint:
		return `'` + strconv.Itoa(int(v)) + `'`

	case uint8:
		return `'` + strconv.Itoa(int(v)) + `'`

	case uint32:
		return `'` + strconv.Itoa(int(v)) + `'`

	case uint64:
		return `'` + strconv.Itoa(int(v)) + `'`

	}

	// have literal to range through

	var s string
	var instr bool
	for _, r := range lit {

		if 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' {
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
		s += " x" + fmt.Sprintf(`%x`, r)

	}

	if instr {
		s += "'"
	}

	if strings.Index(s[1:], " ") > 0 {
		return `(` + s[1:] + `)`
	}
	return s[1:]
}

// ----------------------------- RuleAnyN -----------------------------

// Specific number (n) of any rune as in "rune{n}".
type Any int

func (n Any) String() string {
	if n < 0 {
		return ""
	}
	switch n {
	case 0:
		return `!rune`
	case 1:
		return `rune`
	default:
		return `rune{` + strconv.Itoa(int(n)) + `}`
	}
}

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

func (c *Grammar) RuleAnyN(n int) Rule {
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

func (c *Grammar) RuleLiteral(s string) Rule {
	rule := Lit(s).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ----------------------------- RuleOneOf ----------------------------

// One of rules matches as in "(foo / bar)".
type One []Rule

func (rules One) String() string {
	switch len(rules) {
	case 0:
		return ""
	case 1:
		return rules[0].Text
	default:
		str := "(" + rules[0].Text
		for _, v := range rules[1:] {
			str += " / " + v.Text
		}
		return str + ")"
	}
}

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

func (c *Grammar) RuleOneOf(rules ...Rule) Rule {
	rule := One(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// --------------------------- RuleSequence ---------------------------

// Sequence as in "(foo bar)".
type Seq []Rule

func (rules Seq) String() string {
	switch len(rules) {
	case 0:
		return ""
	case 1:
		return rules[0].Text
	default:
		str := "(" + rules[0].Text
		for _, v := range rules[1:] {
			str += " " + v.Text
		}
		return str + ")"
	}
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

func (c *Grammar) RuleSequence(rules ...Rule) Rule {
	rule := Seq(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// --------------------------- RuleOptional ---------------------------

// Optional as in "(foo bar)?".
type Opt []Rule

func (rules Opt) String() string {
	switch len(rules) {
	case 0:
		return ""
	case 1:
		return rules[0].Text + `?`
	default:
		str := "(" + rules[0].Text
		for _, v := range rules[1:] {
			str += " " + v.Text
		}
		return str + ")?"
	}
}

func (rules Opt) Rule() Rule {
	rule := Rule{
		Text: rules.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleOptional(rules ...Rule) Rule {
	rule := Opt(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ---------------------------- RuleRepeat ----------------------------

// Repeat n times as in "(foo bar){n}".
type Rep struct {
	N     int
	Rules []Rule
}

func (it Rep) String() string {
	if it.N < 0 {
		return ""
	}
	if it.N == 0 {
		return Not(it.Rules).String()
	}
	if it.N == 1 {
		return Seq(it.Rules).String()
	}
	if len(it.Rules) == 1 {
		return it.Rules[0].Text + `{` + strconv.Itoa(it.N) + `}`
	}
	str := `(` + it.Rules[0].Text
	for _, v := range it.Rules[1:] {
		str += " " + v.Text
	}
	return str + `){` + strconv.Itoa(it.N) + `}`
}

func (rules Rep) Rule() Rule {
	rule := Rule{
		Text: rules.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleRepeat(n int, rules ...Rule) Rule {
	rule := Rep{n, rules}.Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ---------------------------- RuleMinMax ----------------------------

// Minimum (m) and maximum (n) as in "(foo bar){m,n}".
type MMax struct {
	Min   int
	Max   int
	Rules []Rule
}

func (it MMax) String() string {
	if it.Min <= 0 || it.Max <= 0 || it.Max < it.Min {
		return ""
	}
	if it.Min == it.Max {
		return Rep{it.Min, it.Rules}.String()
	}
	str := `(` + it.Rules[0].Text
	for _, v := range it.Rules[1:] {
		str += " " + v.Text
	}
	return str + `){` + strconv.Itoa(it.Min) + `,` + strconv.Itoa(it.Max) + `}`
}

func (it MMax) Rule() Rule {
	rule := Rule{
		Text: it.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleMinMax(min, max int, rules ...Rule) Rule {
	rule := MMax{min, max, rules}.Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ------------------------------ RuleMin -----------------------------

// Minimum (m) as in "(foo bar){m,}".
type Min struct {
	Min   int
	Rules []Rule
}

func (it Min) String() string {
	if it.Min <= 0 {
		return ""
	}
	str := `(` + it.Rules[0].Text
	for _, v := range it.Rules[1:] {
		str += " " + v.Text
	}
	return str + `){` + strconv.Itoa(it.Min) + `,}`
}

func (it Min) Rule() Rule {
	rule := Rule{
		Text: it.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleMin(min int, rules ...Rule) Rule {
	rule := Min{min, rules}.Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ------------------------------ RuleMax -----------------------------

// Maximum (m) as in "(foo bar){0,m}".
type Max struct {
	Max   int
	Rules []Rule
}

func (it Max) String() string {
	if it.Max < 0 {
		return ""
	}
	if it.Max == 0 {
		return Not(it.Rules).String()
	}
	str := `(` + it.Rules[0].Text
	for _, v := range it.Rules[1:] {
		str += " " + v.Text
	}
	return str + `){0,` + strconv.Itoa(it.Max) + `}`
}

func (it Max) Rule() Rule {
	rule := Rule{
		Text: it.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleMax(max int, rules ...Rule) Rule {
	rule := Max{max, rules}.Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ----------------------------- RuleRange ----------------------------

// Range as in "[a-z]", "[A-Z]", ,"[0-7]", "[x40-x54]". Note that PEGN
// supports ranges outside of the Unicode valid ranges, but rat does
// not.
type Rng struct {
	Beg rune
	End rune
}

func (it Rng) String() string {
	if it.Beg == 0 && it.End == 0 {
		return ""
	}
	if it.Beg == it.End {
		return Lit(it.Beg).String()
	}
	return `[` + ToPEGN(it.Beg) + `-` + ToPEGN(it.End) + `]`
}

func (it Rng) Rule() Rule {
	rule := Rule{
		Text: it.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleRange(beg, end rune) Rule {
	rule := Rng{beg, end}.Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ----------------------------- RuleMin0 -----------------------------

// Zero minimum as in "(foo bar)*".
type Min0 []Rule

func (rules Min0) String() string {
	switch len(rules) {
	case 0:
		return ""
	case 1:
		return rules[0].Text + `*`
	default:
		str := `(` + rules[0].Text
		for _, v := range rules[1:] {
			str += " " + v.Text
		}
		return str + `)*`
	}
}

func (rules Min0) Rule() Rule {
	rule := Rule{
		Text: rules.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleMin0(rules ...Rule) Rule {
	rule := Min0(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ----------------------------- RuleMin1 -----------------------------

// One minimum as in "(foo bar)+"
type Min1 []Rule

func (rules Min1) String() string {
	switch len(rules) {
	case 0:
		return ""
	case 1:
		return rules[0].Text + `+`
	default:
		str := `(` + rules[0].Text
		for _, v := range rules[1:] {
			str += " " + v.Text
		}
		return str + `)+`
	}
}

func (rules Min1) Rule() Rule {
	rule := Rule{
		Text: rules.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleMin1(rules ...Rule) Rule {
	rule := Min1(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ---------------------------- RuleUnicode ---------------------------

// Unicode class as in p{L}, p{Arabic}, etc.
type U string

func (it U) String() string {
	if len(it) == 0 {
		return ""
	}
	return `p{` + string(it) + `}`
}

func (it U) Rule() Rule {
	rule := Rule{
		Text: it.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleUnicode(uni string) Rule {
	rule := U(uni).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// -------------------------------- Set -------------------------------

// Set from string as in "('a' / 'e' / 'i' / 'o' / 'u')" or "[aeiou]"
// from regular expressions. To invert combine with Not (ex: Not{Set})
type Set string

func (it Set) String() string {
	switch len(it) {
	case 0:
		return ""
	case 1:
		return ToPEGN(string(it[0]))
	default:
		str := `(` + ToPEGN(string(it[0]))
		for _, v := range it[1:] {
			str += " / " + ToPEGN(string(v))
		}
		return str + `)`
	}
}

func (it Set) Rule() Rule {
	rule := Rule{
		Text: it.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleSet(set string) Rule {
	rule := Set(set).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// --------------------------------- N --------------------------------

// Named capture as in "Foo <- 'foo'" or "<:Foo 'foo' >"
// or "Foo <- ws* < 'foo' > ws*" (rat.N{`Foo`,rules}). The Name becomes
// the Rule.Text that is cached when a Grammar is used.
type N struct {
	Name  string
	Rules []Rule
}

func (it N) String() string {

	if it.Name == "" {
		return Seq(it.Rules).String()
	}

	if len(it.Rules) == 0 {
		return ""
	}

	// _foo -> <:foo 'foo' >
	if it.Name[0] == '_' {
		str := `<:` + it.Name + ` ` + it.Rules[0].Text
		if len(it.Rules) > 1 {
			for _, v := range it.Rules[1:] {
				str += " " + v.Text
			}
		}
		return str + " >"
	}

	str := it.Name + ` <- ` + it.Rules[0].Text
	if len(it.Rules) > 1 {
		for _, v := range it.Rules[1:] {
			str += " " + v.Text
		}
	}

	return str
}

func (it N) Rule() Rule {
	rule := Rule{
		Text: it.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleNamed(name string, rules ...Rule) Rule {
	rule := N{name, rules}.Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ------------------------------ RuleIs ------------------------------

// Positive lookahead as in "&(foo bar)".
type Is []Rule

func (rules Is) String() string {
	switch len(rules) {
	case 0:
		return ""
	case 1:
		return rules[0].Text
	default:
		str := `&(` + rules[0].Text
		if len(rules) > 1 {
			for _, rule := range rules {
				str += " " + rule.Text
			}
		}
		return str + `)`
	}
}

func (rules Is) Rule() Rule {
	rule := Rule{
		Text: rules.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleIs(rules ...Rule) Rule {
	rule := Is(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// -------------------------------- Not -------------------------------

// Negative lookahead as in "!(foo bar)".
type Not []Rule

func (rules Not) String() string {
	switch len(rules) {
	case 0:
		return ""
	case 1:
		return `!` + rules[0].Text
	default:
		str := `!(` + rules[0].Text
		if len(rules) > 1 {
			for _, rule := range rules {
				str += " " + rule.Text
			}
		}
		return str + `)`
	}
}

func (rules Not) Rule() Rule {
	rule := Rule{
		Text: rules.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleNot(rules ...Rule) Rule {
	rule := Not(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// -------------------------------- To --------------------------------

// Up to non-inclusive as in "...(foo bar)".
type To []Rule

func (rules To) String() string {
	switch len(rules) {
	case 0:
		return ""
	case 1:
		return `...` + rules[0].Text
	default:
		str := `...(` + rules[0].Text
		if len(rules) > 1 {
			for _, rule := range rules {
				str += " " + rule.Text
			}
		}
		return str + `)`
	}
}

func (rules To) Rule() Rule {
	rule := Rule{
		Text: rules.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleTo(rules ...Rule) Rule {
	rule := To(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// -------------------------------- ToI -------------------------------

// Up to inclusive as in "..(foo bar)".
type Toi []Rule

func (rules Toi) String() string {
	switch len(rules) {
	case 0:
		return ""
	case 1:
		return `..` + rules[0].Text
	default:
		str := `..(` + rules[0].Text
		if len(rules) > 1 {
			for _, rule := range rules {
				str += " " + rule.Text
			}
		}
		return str + `)`
	}
}

func (rules Toi) Rule() Rule {
	rule := Rule{
		Text: rules.String(),
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

func (c *Grammar) RuleToInc(rules ...Rule) Rule {
	rule := Toi(rules).Rule()
	c.Store(rule.Text, rule)
	return rule
}

// ------------------------------- PEGN -------------------------------

// PEGN communicates to Pack that a string has specific meaning and PEGN
// syntax, otherwise Pack assumes all strings as Lit types.
type PEGN struct {
	Name string
	Text string
}

func (p PEGN) String() string {
	// TODO run through formatter and output pretty form
	return ""
}

func (p PEGN) Rule() Rule {
	rule := Rule{
		Text: p.String(), // for pretty print formatting
	}
	rule.Check = func(r []rune, i int) Result {
		// TODO
		return Result{B: i, E: i}
	}
	return rule
}

// PEGN parses the input (string, []byte, []rune, or PEGN) and returns
// a single top level Rule covering the mini-grammar passed as an
// argument. A name is required to avoid unnecessarily long caching
// keys.
func (g *Grammar) PEGN(name string, in any) Rule {
	var rule Rule

	switch v := in.(type) {
	case string:
		rule = PEGN{name, v}.Rule()
	case []byte:
		rule = PEGN{name, string(v)}.Rule()
	case []rune:
		rule = PEGN{name, string(v)}.Rule()
	case PEGN:
		rule.Text = name
		rule = v.Rule()
	}

	g.Store(rule.Text, rule)
	return rule
}

// ------------------------------- Check ------------------------------

// Check is a convenience function that takes a PEGN rule definition,
// dynamically compiles it, and checks the input returning the Result.
// The PEGN grammar (first argument) may be any of the following types:
//
//     string
//     []byte
//     []rune
//     PEGN (name plus string)
//     io.Reader
//
// The UTF-8 text to check (second argument) may be any of the following:
//
//     string
//     []byte
//     []rune
//     io.Reader
//
// This function is slower than alternatives (like a regular
// expression without compilation).
func Check(pegn any, in any) Result {
	var rule Rule
	var runes []rune

	switch v := pegn.(type) {
	case string:
		rule = PEGN{`_dynamic`, v}.Rule()
	case []byte:
		rule = PEGN{`_dynamic`, string(v)}.Rule()
	case []rune:
		rule = PEGN{`_dynamic`, string(v)}.Rule()
	case PEGN:
		rule = v.Rule()
	}

	switch v := pegn.(type) {
	case string:
		runes = []rune(v)
	case []byte:
		runes = []rune(string(v))
	case []rune:
		runes = v
	case PEGN:
		rule = v.Rule()
	case io.Reader:
		buf, _ := io.ReadAll(v)
		runes = []rune(string(buf))
	}

	switch v := in.(type) {
	case string:
		runes = []rune(v)
	case []byte:
		runes = []rune(string(v))
	case []rune:
		runes = v
	case io.Reader:
		buf, _ := io.ReadAll(v)
		runes = []rune(string(buf))
	}

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

// -------------------------------- Gen -------------------------------

type Generator interface {
	// TODO
}

type GenParams struct {
	Generator
	Out io.Writer
}

// Gen creates compilable source code to be used in projects that
// need a pre-compiled parser/scanner to implement domain-specific
// languages and such. Since code generation preference differ widely
// based on the needs of the project and language, GenParams passed
// allows any Generator implementation to be passed.
func (g *Grammar) Gen(p GenParams) error {
	// TODO
	return nil
}

// Gen is a shortcut for rat.Pack().Gen()
func Gen(in ...any) error {
	//g := Pack(in...)
	// TODO
	return nil
}

// ------------------------------ Errors ------------------------------

type ErrExpected struct {
	It any
}

func (e ErrExpected) Error() string {
	switch v := e.It.(type) {
	case rune:
		e.It = ToPEGN(string(v))
	}
	return fmt.Sprintf(_ErrExpected, e.It)
}

type ErrNotExist struct{ It any }

func (e ErrNotExist) Error() string {
	return fmt.Sprintf(_ErrNotExist, e.It)
}

type ErrBadType struct{ It any }

func (e ErrBadType) Error() string {
	return fmt.Sprintf(_ErrBadType, e.It)
}
