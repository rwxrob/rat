package rat

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/rwxrob/rat/x"
)

// Trace enables tracing parsing and checks while they happen. Output is
// with the log package.
var Trace int

// DefaultRuleName is used by NewRule and AddRule as the prefix for new,
// unnamed rules and is combined with the internal integer number that
// increments for each new such rule added.
var DefaultRuleName = `Rule`

// Grammar is an aggregation of cached (memoized) rules with methods to
// add rules and check data against them. The Main rule is the entry
// point. Grammars may have multiple unrelated rules added and change
// the Main entry point dynamically as needed, but most will use the
// sequence rule created by calling the Pack method or it's rat.Pack
// equivalent constructor. Trace may be incremented during debugging to
// gain performant visibility into grammar construction and scanning.
type Grammar struct {
	Trace   int              // activate logs for debug visibility
	Rules   map[string]*Rule // keyed to Rule.Name (not Text)
	Main    *Rule            // entry point for Check or Scan
	RuleNum int              // auto-incrementing for ever unnamed rule added.
}

// String fulfills the fmt.Stringer interface by producing compilable Go
// code containing the Main rule (usually a rat/x.Seq). In this way,
// code generators for specific, dynamically created grammars can easily
// be created.
func (g Grammar) String() string {
	var str string
	if g.Main != nil {
		str += g.Main.Text
	}
	return str
}

func (g Grammar) Print() { fmt.Println(g) }

// Check delegates to g.Main.Check.
func (g *Grammar) Check(r []rune, i int) Result { return g.Main.Check(r, i) }

// Scan checks the input against the current g.Main rule. It is
// functionally identical to Check but accepts []rune, string, []byte,
// and io.Reader as input. The error (X) on Result is set if there is
// a problem.
func (g *Grammar) Scan(in any) Result {
	if g.Main == nil {
		return Result{X: ErrIsZero{g.Main}}
	}

	var runes []rune
	switch v := in.(type) {

	case string:
		runes = []rune(v)

	case []byte:
		runes = []rune(string(v))

	case []rune:
		runes = v

	case io.Reader:
		buf, err := io.ReadAll(v)
		if err != nil {
			return Result{X: err}
		}
		runes = []rune(string(buf))
	}

	return g.Main.Check(runes, 0)
}

// Pack creates a x.Seq rule from the input and assigns it as the Main
// rule returning a reference to the updated Grammar itself.
func (g *Grammar) Pack(seq ...any) *Grammar {
	xseq := x.Seq(seq)
	rule := g.MakeRule(xseq)
	g.Main = rule
	return g
}

// MakeRule fulfills the MakeRule interface. The input argument is
// usually a rat/x ("ratex") expression type. Anything else is
// interpreted as a literal string by using it's String method or
// converting it into a string using the %v (string, []rune, []byte,
// rune) or %q representation.
func (g *Grammar) MakeRule(in any) *Rule {

	if g.Trace > 0 || Trace > 0 {
		log.Printf("MakeRule(%v)", x.String(in))
	}

	switch v := in.(type) {

	// text (most common)
	case string, []rune, []byte, rune, x.Lit:
		return g.MakeLit(v)

	// rat/x ("ratex") types as expressions
	case x.Name:
		return g.MakeName(v)
	case x.ID:
		return g.MakeID(v)
	case x.Ref:
		return g.MakeRef(v)
	case x.Rid:
		return g.MakeRid(v)
	case x.Is:
		return g.MakeIs(v)
	case x.Seq:
		return g.MakeSeq(v)
	case x.One:
		return g.MakeOne(v)
	case x.Opt:
		return g.MakeOpt(v)
	case x.Mn1:
		return g.MakeMn1(v)
	case x.Mn0:
		return g.MakeMn0(v)
	case x.Min:
		return g.MakeMin(v)
	case x.Max:
		return g.MakeMax(v)
	case x.Mmx:
		return g.MakeMmx(v)
	case x.Rep:
		return g.MakeRep(v)
	case x.Pos:
		return g.MakePos(v)
	case x.Neg:
		return g.MakeNeg(v)
	case x.Any:
		return g.MakeAny(v)
	case x.Toi:
		return g.MakeToi(v)
	case x.Tox:
		return g.MakeTox(v)
	case x.Rng:
		return g.MakeRng(v)
	case x.End:
		return g.MakeEnd(v)

	case fmt.Stringer:
		return g.MakeLit(v.String())

	// anything that has an %q form
	default:
		return g.MakeLit(fmt.Sprintf(`%q`, v))

	}
}

// NewRule creates a new rule in the grammar cache using the defaults.
// It is a convenience when the Name and ID are not needed. See AddRule
// for details.
func (g *Grammar) NewRule() *Rule {
	rule := new(Rule)
	g.AddRule(rule)
	return rule
}

// AddRule adds a new rule to the grammar cache keyed to the rule.Name.
// If a rule was already keyed to that name it is overwritten. If the
// rule.ID is 0, a new arbitrary ID is assigned that begins with -1 and
// decreases for every new rule added with a 0 value. If the
// rule.Name is empty the positive value of the ID is combined with the
// DefaultRuleName prefix to provide a generic name.
// Avoid changing the rule.Name or rule.ID values after added
// since the key in the grammar cache is hard-coded to the rule.Name
// when called. If the rule.Name and rule.ID are not important consider
// NewRule instead (which uses these defaults and requires no argument).
// Returns self for convenience.
func (g *Grammar) AddRule(rule *Rule) *Rule {
	if rule.ID == 0 {
		g.RuleNum++
		rule.ID = g.RuleNum
	}
	if rule.Name == "" {
		rule.Name = DefaultRuleName + strconv.Itoa(-rule.ID)
	}
	if g.Rules == nil {
		g.Rules = map[string]*Rule{}
	}
	g.Rules[rule.Name] = rule
	return rule
}

func (g *Grammar) MakeName(in x.Name) *Rule {

	// syntax check
	if len(in) != 2 {
		panic(x.UsageName)
	}
	name, isstring := in[0].(string)
	if !isstring {
		panic(x.UsageName)
	}

	// check the cache for a rule with this name
	rule, has := g.Rules[name]
	if has {
		return rule
	}

	// check the cache for the encapsulated rule, else make one
	text := x.String(in[1])
	irule, rulecached := g.Rules[text]
	if !rulecached {
		irule = g.MakeRule(in[1])
	}

	rule = new(Rule)
	rule.Name = name
	rule.Text = in.String()
	rule.Check = irule.Check

	return g.AddRule(rule)
}

func (g *Grammar) MakeID(in x.ID) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}

func (g *Grammar) MakeRef(in x.Ref) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}
func (g *Grammar) MakeRid(in x.Rid) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}
func (g *Grammar) MakeIs(in x.Is) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}

func (g *Grammar) MakeSeq(seq x.Seq) *Rule {

	name := seq.String()
	if r, have := g.Rules[name]; have {
		return r
	}

	rule := new(Rule)
	rule.Name = name
	rule.Text = name

	rules := []*Rule{}

	for _, it := range seq {

		rule := g.MakeRule(it)

		if rule == nil || rule.Check == nil {
			log.Printf(`skipping invalid Rule: %v`, rule)
			continue
		}

		rules = append(rules, rule)
	}

	rule.Check = func(r []rune, i int) Result {
		start := i
		results := []Result{}
		for _, rule := range rules {
			result := rule.Check(r, i)
			i = result.E
			results = append(results, result)
			if result.X != nil {
				return Result{R: r, B: start, E: i, C: results, X: result.X}
			}
		}
		return Result{R: r, B: start, E: i, C: results}
	}

	return rule
}

func (g *Grammar) MakeOne(one x.One) *Rule {

	ln := len(one)
	if ln < 1 {
		panic(x.UsageOne)
	}

	name := one.String()
	if r, have := g.Rules[name]; have {
		return r
	}

	// just one is same as rule by itself
	if ln == 1 {
		return g.MakeRule(one[0])
	}

	// create/fetch rules for every possibility to cache and enclose in Check
	rules := make([]*Rule, ln)
	for n, exp := range one {
		rules[n] = g.MakeRule(exp)
	}

	rule := new(Rule)
	rule.Name = name
	rule.Text = name

	rule.Check = func(r []rune, i int) Result {
		start := i
		for _, it := range rules {
			result := it.Check(r, i)
			if result.X == nil {
				return result
			}
		}
		return Result{R: r, B: start, E: i, X: ErrExpected{one}}
	}

	return g.AddRule(rule)

}

func (g *Grammar) MakeOpt(in x.Opt) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}

func (g *Grammar) MakeLit(in any) *Rule {

	var val string

	switch v := in.(type) {
	case string:
		val = v
	case []rune:
		val = string(v)
	case []byte:
		val = string(v)
	case rune:
		val = string(v)
	case x.Lit:
		return g.MakeLit(x.JoinLit(v...))

	}

	rule, has := g.Rules[val]
	if has {
		return rule
	}

	rule = new(Rule)
	rule.Name = val
	rule.Text = x.String(val)
	g.AddRule(rule)

	rule.Check = func(r []rune, i int) Result {
		var err error
		start := i
		runes := []rune(val)
		var n int
		runeslen := len(runes)
		for i < len(r) && n < runeslen {
			if r[i] != runes[n] {
				err = ErrExpected{r[n]}
				break
			}
			i++
			n++
		}
		if n < runeslen {
			err = ErrExpected{string(runes[n])}
		}
		return Result{R: r, B: start, E: i, X: err}
	}

	return rule
}

func (g *Grammar) MakeMn1(in x.Mn1) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}
func (g *Grammar) MakeMn0(in x.Mn0) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}

func (g *Grammar) MakeMin(in x.Min) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}

func (g *Grammar) MakeMax(in x.Max) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}

func (g *Grammar) MakeMmx(in x.Mmx) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}

func (g *Grammar) MakeRep(in x.Rep) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}

func (g *Grammar) MakePos(in x.Pos) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}
func (g *Grammar) MakeNeg(in x.Neg) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}

func (g *Grammar) MakeAny(in x.Any) *Rule {

	name := in.String()
	if r, have := g.Rules[name]; have {
		return r
	}

	switch len(in) {
	case 1:
		return g.makeAnyN(in)
	case 2:
		return g.makeAnyMmx(in)
	default:
		panic(x.UsageAny)
	}
}

// only call from makeAny
func (g *Grammar) makeAnyN(in x.Any) *Rule {

	n, is := in[0].(int)
	if !is {
		panic(x.UsageAny)
	}

	name := in.String()
	rule := new(Rule)
	rule.Name = name
	rule.Text = name
	g.AddRule(rule)

	rule.Check = func(r []rune, i int) Result {
		start := i
		if i+n > len(r) {
			return Result{R: r, B: start, E: len(r) - 1, X: ErrExpected{rule.Name}}
		}
		return Result{R: r, B: start, E: i + n}
	}

	return rule
}

// only call from makeAny
func (g *Grammar) makeAnyMmx(in x.Any) *Rule {

	m, is := in[0].(int)
	if !is {
		panic(x.UsageAny)
	}

	n, is := in[1].(int)
	if !is {
		panic(x.UsageAny)
	}

	if m >= n || n <= 0 {
		panic(x.UsageAny)
	}

	name := in.String()
	rule := new(Rule)
	rule.Name = name
	rule.Text = name
	g.AddRule(rule)

	rule.Check = func(r []rune, i int) Result {
		start := i

		// minimum is more than we have
		if i+m > len(r) {
			return Result{R: r, B: start, E: len(r) - 1, X: ErrExpected{rule.Name}}
		}

		// we have enough for max
		if i+n < len(r) {
			return Result{R: r, B: start, E: i + n}
		}

		// we have less than max, but more than min
		return Result{R: r, B: start, E: len(r)}
	}

	return rule

}

func (g *Grammar) MakeToi(in x.Toi) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}
func (g *Grammar) MakeTox(in x.Tox) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}
func (g *Grammar) MakeRng(in x.Rng) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}
func (g *Grammar) MakeEnd(in x.End) *Rule {
	rule := new(Rule)
	// TODO
	return rule
}
