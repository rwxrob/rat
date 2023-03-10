package rat

import (
	"fmt"
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
//
// Memoization
//
// All Make* methods check the Rules map/cache for a match for the
// String form of the rat/x expression and return it directly if found
// rather than create a new Rule with an identical CheckFunc. The
// MakeNamed creates an additional entry (pointing to the same *Rule)
// for the specified name.
//
type Grammar struct {
	Trace int              // activate logs for debug visibility
	Rules map[string]*Rule // keyed to Rule.Name (not Text)
	Saved map[string]*Rule // dynamically created literals from Sav
	Main  *Rule            // entry point for Check or Scan

	ruleid int // auto-incrementing for ever unnamed rule added.
}

// Init initializes the Grammar emptying the Rules if any or creating
// a new Rules map and setting internal rule ID to 0 and disabling
// Trace, and emptying Main.
func (g *Grammar) Init() *Grammar {
	g.Trace = 0
	g.Rules = map[string]*Rule{}
	g.Saved = map[string]*Rule{}
	g.Main = nil
	g.ruleid = 0
	return g
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
	return g.Main.Scan(in)
}

// Pack allows multiple rules to be passed (unlike MakeRule). If one
// argument, returns MakeRule for it. If more than one argument,
// delegates to MakeSeq. Pack is called from the package function of the
// same name, which describes the valid argument types. As a convenience,
// a self-reference is returned.
func (g *Grammar) Pack(in ...any) *Grammar {
	var rule *Rule
	switch len(in) {
	case 0:
		panic(ErrArgs{in})
	case 1:
		rule = g.MakeRule(in[0])
	default:
		rule = g.MakeSeq(x.Seq(in))
	}
	g.Main = rule
	g.AddRule(rule)
	return g
}

// MakeRule fulfills the MakeRule interface. The input argument is
// usually a rat/x ("ratex") expression type including x.IsFunc functions.
// Anything else is interpreted as a literal string by using its String
// method or converting it into a string using the %v (string, []rune, []
// byte, rune) or %q representation. Note that MakeRule itself does not
// check the Rules cache for existing Rules not does it add the rule to
// that cache. This work is left to the Make* methods themselves or to
// the AddRule method. The result, however, is the same since MakeRule
// delegates to those Make* methods.
func (g *Grammar) MakeRule(in any) *Rule {

	if g.Trace > 0 || Trace > 0 {
		log.Printf("MakeRule(%v)", x.String(in))
	}

	switch v := in.(type) {

	// text (most common)
	case string, []rune, []byte, rune, x.Str:
		return g.MakeStr(v)

	// rat/x ("ratex") types as expressions
	case x.N:
		return g.MakeNamed(v)
	case x.Sav:
		return g.MakeSave(v)
	case x.Val:
		return g.MakeVal(v)
	case x.Ref:
		return g.MakeRef(v)
	case x.Is:
		return g.MakeIs(v)
	case func(r rune) bool:
		return g.MakeIs(x.Is{v})
	case x.IsFunc:
		return g.MakeIs(x.Is{v})
	case x.Seq:
		return g.MakeSeq(v)
	case x.One:
		return g.MakeOne(v)
	case x.Mmx:
		return g.MakeMmx(v)
	case x.See:
		return g.MakeSee(v)
	case x.Not:
		return g.MakeNot(v)
	case x.To:
		return g.MakeTo(v)
	case x.Any:
		return g.MakeAny(v)
	case x.Rng:
		return g.MakeRng(v)
	case x.End:
		return g.MakeEnd(v)

	case fmt.Stringer:
		return g.MakeStr(v.String())

	case bool:
		return g.MakeStr(fmt.Sprintf(`%v`, v))

	// anything that has an %q form
	default:
		return g.MakeStr(fmt.Sprintf(`%q`, v))

	}
}

// NewRule creates a new rule in the grammar cache using the defaults.
// It is a convenience when a Name is not needed. See AddRule
// for details.
func (g *Grammar) NewRule() *Rule {
	rule := new(Rule)
	g.AddRule(rule)
	return rule
}

// AddRule adds a new rule to the grammar cache keyed to the rule.Name.
// If a rule was already keyed to that name it is overwritten.
// If rule.Name is empty a new incremental name is created with the
// DefaultRuleName prefix.  Avoid changing the rule.Name values after added
// since the key in the grammar cache is hard-coded to the rule.Name
// when called. If the rule.Name is not important consider
// NewRule instead (which uses these defaults and requires no argument).
// Returns self for convenience.
func (g *Grammar) AddRule(rule *Rule) *Rule {
	if rule.Name == "" {
		g.ruleid++
		rule.Name = DefaultRuleName + strconv.Itoa(g.ruleid)
	}
	g.Rules[rule.Name] = rule
	return rule
}

// MakeNamed makes two rules pointing to the same CheckFunc, one unnamed
// and other named (first argument). Both produce results that have the
// Name field set.
func (g *Grammar) MakeNamed(in x.N) *Rule {

	text := in.String()

	rule, has := g.Rules[text]
	if has {
		return rule
	}

	if len(in) != 2 {
		panic(x.UsageN)
	}

	name, is := in[0].(string)
	if !is {
		panic(x.UsageN)
	}

	// check the cache for the encapsulated rule, else make one
	iname := x.String(in[1])
	irule, has := g.Rules[iname]
	if !has {
		irule = g.MakeRule(in[1])
	}

	rule = &Rule{Name: name, Text: in.String()}
	g.AddRule(rule)

	rule.Check = func(r []rune, i int) Result {
		unnamed := irule.Check(r, i)
		unnamed.N = name
		return unnamed
	}

	return rule
}

func (g *Grammar) MakeRef(in x.Ref) *Rule {

	name := in.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	if len(in) != 1 {
		panic(x.UsageRef)
	}

	key, is := in[0].(string)
	if !is {
		panic(x.UsageRef)
	}

	rule = &Rule{Name: name, Text: name}
	g.AddRule(rule)

	rule.Check = func(r []rune, i int) Result {
		rule, has := g.Rules[key]
		if has {
			return rule.Check(r, i)
		}
		return Result{R: r, B: i, E: i, X: ErrExpected{in}}
	}

	return rule
}

func (g *Grammar) MakeSave(in x.Sav) *Rule {

	name := in.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}
	g.AddRule(rule)

	if len(in) != 1 {
		panic(x.UsageSav)
	}

	key, is := in[0].(string)
	if !is {
		panic(x.UsageSav)
	}

	rule.Check = func(r []rune, i int) Result {
		rule, has := g.Rules[key]
		if has {
			res := rule.Check(r, i)
			if res.X == nil {
				g.Saved[key] = g.MakeStr(res.Text())
			}
			return res
		}
		return Result{R: r, B: i, E: i, X: ErrExpected{in}}
	}

	return rule
}

func (g *Grammar) MakeVal(in x.Val) *Rule {

	name := in.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}
	g.AddRule(rule)

	if len(in) != 1 {
		panic(x.UsageVal)
	}

	key, is := in[0].(string)
	if !is {
		panic(x.UsageVal)
	}

	rule.Check = func(r []rune, i int) Result {
		rule, has := g.Saved[key]
		if has {
			return rule.Check(r, i)
		}
		return Result{R: r, B: i, E: i, X: ErrExpected{in}}
	}

	return rule
}

// MakeIs takes an x.IsFunc (which is just a func(r rune) bool) or x.Is
// type and calls that function in its Check.
func (g *Grammar) MakeIs(in x.Is) *Rule {

	if len(in) != 1 {
		panic(x.UsageIs)
	}

	isfunc, is := in[0].(func(r rune) bool)
	if !is {
		panic(x.UsageIs)
	}

	name := in.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}

	rule.Check = func(r []rune, i int) Result {
		if i < len(r) && isfunc(r[i]) {
			return Result{R: r, B: i, E: i + 1}
		}
		return Result{R: r, B: i, E: i, X: ErrExpected{in}}
	}

	return g.AddRule(rule)
}

func (g *Grammar) MakeSeq(seq x.Seq) *Rule {

	if len(seq) == 1 {
		it, is := seq[0].([]any)
		if !is {
			return g.MakeRule(seq[0])
		}
		seq = it
	}

	seq = x.CombineStr(seq...)
	if len(seq) == 1 {
		return g.MakeRule(seq[0])
	}

	name := seq.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}
	g.AddRule(rule)

	rules := []*Rule{}

	for _, it := range seq {

		iname := x.String(it)
		irule, has := g.Rules[iname]
		if !has {
			irule = g.MakeRule(it)
		}

		rules = append(rules, irule)
	}

	rule.Check = func(r []rune, i int) Result {
		start := i
		results := []Result{}

		for _, rule := range rules {
			res := rule.Check(r, i)
			i = res.E
			results = append(results, res)
			if res.X != nil {
				return Result{R: r, B: start, E: i, C: results, X: res.X}
			}
		}

		return Result{R: r, B: start, E: i, C: results}
	}

	return rule
}

func (g *Grammar) MakeOne(one x.One) *Rule {

	name := one.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}
	g.AddRule(rule)

	ln := len(one)
	if ln < 1 {
		panic(x.UsageOne)
	}

	// just one is same as rule by itself
	if ln == 1 {
		return g.MakeRule(one[0])
	}

	// create/fetch rules for every possibility to cache and enclose in Check
	rules := make([]*Rule, ln)
	for n, exp := range one {
		name := x.String(exp)
		irule, has := g.Rules[name]
		if !has {
			irule = g.MakeRule(exp)
			g.AddRule(irule)
		}
		rules[n] = irule
	}

	rule.Check = func(r []rune, i int) Result {
		result := Result{R: r, B: i, E: i}
		for _, it := range rules {
			res := it.Check(r, i)
			if res.X == nil {
				result.E = res.E
				result.C = []Result{res}
				return result
			}
		}
		result.X = ErrExpected{one}
		return result
	}

	return rule

}

func (g *Grammar) MakeStr(in any) *Rule {

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
	case x.Str:
		return g.MakeStr(x.JoinStr(v...))

	}

	name := `x.Str{"` + val + `"}`

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}
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

func (g *Grammar) MakeMmx(in x.Mmx) *Rule {

	name := in.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}
	g.AddRule(rule)

	if len(in) != 3 {
		panic(x.UsageMmx)
	}

	min, is := in[0].(int)
	if !is || min < 0 {
		panic(x.UsageMmx)
	}

	max, is := in[1].(int)
	if !is || (max < min && max != -1) {
		panic(x.UsageMmx)
	}

	iname := x.String(in[2])
	irule, has := g.Rules[iname]
	if !has {
		irule = g.MakeRule(in[2])
	}

	rule.Check = func(r []rune, i int) Result {
		result := Result{R: r, B: i, E: i, C: []Result{}}
		var count int
		var res Result
		for {
			res = irule.Check(r, i)
			if res.X != nil || count == max {
				break
			}
			result.C = append(result.C, res)
			i = res.E
			result.E = i
			count++
		}

		if min <= count && count <= max {
			if res.X == nil {
				result.C = append(result.C, res)
			}
			return result
		}
		result.X = ErrExpected{in}
		return result
	}

	return rule
}

func (g *Grammar) MakeSee(in x.See) *Rule {

	name := in.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}
	g.AddRule(rule)

	if len(in) != 1 {
		panic(x.UsageSee)
	}

	iname := x.String(in[0])
	irule, has := g.Rules[iname]
	if !has {
		irule = g.MakeRule(in[0])
	}

	rule.Check = func(r []rune, i int) Result {
		result := Result{R: r, B: i, E: i}
		res := irule.Check(r, i)
		if res.X == nil {
			return result
		}
		result.X = ErrExpected{in}
		return result
	}

	return rule

}

func (g *Grammar) MakeNot(in x.Not) *Rule {

	name := in.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}
	g.AddRule(rule)

	if len(in) != 1 {
		panic(x.UsageNot)
	}

	iname := x.String(in[0])
	irule, has := g.Rules[iname]
	if !has {
		irule = g.MakeRule(in[0])
	}

	rule.Check = func(r []rune, i int) Result {
		result := Result{R: r, B: i, E: i}
		res := irule.Check(r, i)
		if res.X != nil {
			return result
		}
		result.X = ErrExpected{in}
		return result
	}

	return rule

}

func (g *Grammar) MakeTo(in x.To) *Rule {

	name := in.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}
	g.AddRule(rule)

	if len(in) != 1 {
		panic(x.UsageTo)
	}

	iname := x.String(in[0])
	irule, has := g.Rules[iname]
	if !has {
		irule = g.MakeRule(in[0])
	}

	rule.Check = func(r []rune, i int) Result {
		result := Result{R: r, B: i, E: i}

		for ; i < len(r); i++ {
			res := irule.Check(r, i)
			if res.X == nil {
				return result
			}
			result.E++
		}

		result.X = ErrExpected{in}
		return result
	}

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
	rule := &Rule{Name: name, Text: name}
	g.AddRule(rule)

	rule.Check = func(r []rune, i int) Result {
		start := i
		if i+n > len(r) {
			return Result{R: r, B: start, E: len(r) - 1, X: ErrExpected{in}}
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
	rule := &Rule{Name: name, Text: name}
	g.AddRule(rule)

	rule.Check = func(r []rune, i int) Result {
		start := i

		// minimum is more than we have
		if i+m > len(r) {
			return Result{R: r, B: start, E: len(r) - 1, X: ErrExpected{in}}
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

func (g *Grammar) MakeRng(in x.Rng) *Rule {

	name := in.String()

	rule, has := g.Rules[name]
	if has {
		return rule
	}

	rule = &Rule{Name: name, Text: name}
	g.AddRule(rule)

	if len(in) != 2 {
		panic(x.UsageRng)
	}

	beg, is := in[0].(rune)
	if !is {
		panic(x.UsageRng)
	}

	end, is := in[1].(rune)
	if !is {
		panic(x.UsageRng)
	}

	rule.Check = func(r []rune, i int) Result {
		result := Result{R: r, B: i, E: i}
		cur := r[i]
		if beg <= cur && cur <= end {
			result.E++
			return result
		}
		result.X = ErrExpected{in}
		return result
	}

	return rule

}

func (g *Grammar) MakeEnd(in x.End) *Rule {

	if len(in) != 0 {
		panic(x.UsageEnd)
	}

	name := in.String()
	rule := new(Rule)
	rule.Name = name
	rule.Text = name

	rule.Check = func(r []rune, i int) Result {
		if i == len(r) {
			return Result{R: r, B: i, E: i}
		}
		return Result{R: r, B: i, E: i, X: ErrExpected{in}}
	}

	return g.AddRule(rule)
}
