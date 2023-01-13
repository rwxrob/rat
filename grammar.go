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

// additive only
type Grammar struct {
	Trace   int
	names   []string
	rules   map[string]*Rule
	rule    *Rule
	rulenum int
}

// TODO MarshalJSON - just the meta stuff to help others implement compatible
// TODO UnmarshalJSON - just the meta stuff with Checks unassigned unset

func (g Grammar) String() string {
	var str string
	if g.rule != nil {
		str += g.rule.Text
	}
	return str
}

func (g Grammar) Print() { fmt.Println(g) }

// SetScanRule sets the primary entry rule used when g.Scan is called.
func (g *Grammar) SetScanRule(rule *Rule) *Grammar {
	g.rule = rule
	return g
}

// Scan checks the input against the current scan rule (see SetScanRule).
func (g *Grammar) Scan(in any) Result {
	if g.rule == nil {
		return Result{X: ErrIsZero{g.rule}}
	}
	var runes []rune
	switch v := in.(type) {
	case string:
		runes = []rune(v)
	case []byte:
		runes = []rune(string(v))
	case []rune:
		runes = v
	}
	return g.rule.Check(runes, 0)
}

// Pack is shorthand for g.SetScanRule(g.MakeRule(x.Seq(seq))).
func (g *Grammar) Pack(seq ...any) *Grammar {
	xseq := x.Seq(seq)
	rule := g.MakeRule(xseq)
	return g.SetScanRule(rule)
}

// MakeRule fulfills the MakeRule interface. The input argument is
// usually a rat/x expression type. Anything else is interpreted as
// a literal string by converting it into a string using the %v or %q
// representation. This includes anything that implements the
// fmt.Stringer interface.
func (g *Grammar) MakeRule(in any) *Rule {

	if g.Trace > 0 || Trace > 0 {
		log.Printf("MakeRule(%v)", x.String(in))
	}

	switch v := in.(type) {

	// text
	case string:
		return g.makeLit(v)
	case []rune:
		return g.makeLit(string(v))
	case []byte:
		return g.makeLit(string(v))
	case rune:
		return g.makeLit(string(v))

	// rat/x types as expressions
	case x.Name:
		return g.makeName(v)
	case x.Seq:
		return g.makeSeq(v)
	case x.Any:
		return g.makeAny(v)

	case fmt.Stringer:
		return g.makeLit(v.String())

	// anything that has an %q form
	default:
		return g.makeLit(fmt.Sprintf(`%q`, v))

	}

	log.Printf(`unsupported type: %T`, in) // TODO better error
	return nil
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
		g.rulenum++
		rule.ID = g.rulenum
	}
	if rule.Name == "" {
		rule.Name = DefaultRuleName + strconv.Itoa(-rule.ID)
	}
	if g.rules == nil {
		g.rules = map[string]*Rule{}
	}
	g.rules[rule.Name] = rule
	return rule
}

func (g *Grammar) makeName(in x.Name) *Rule {

	// syntax check
	if len(in) != 2 {
		panic(x.UsageName)
	}
	name, isstring := in[0].(string)
	if !isstring {
		panic(x.UsageName)
	}

	// check the cache, return if found
	rule, has := g.rules[name]
	if has {
		return rule
	}

	rule = g.MakeRule(in[1])
	rule.Name = name

	return rule
}

func (g *Grammar) makeLit(in string) *Rule {
	rule, has := g.rules[in]
	if has {
		return rule
	}

	rule = new(Rule)
	rule.Name = in
	rule.Text = x.String(in)
	g.AddRule(rule)

	rule.Check = func(r []rune, i int) Result {
		var err error
		start := i
		runes := []rune(in)
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

func (g *Grammar) makeAny(in x.Any) *Rule {

	name := in.String()
	if r, have := g.rules[name]; have {
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
		if i+n > len(r) {
			return Result{R: r, B: start, E: i + m}
		}

		// we have less than max, but more than min
		return Result{R: r, B: start, E: len(r) - 1}
	}

	return rule

}

func (g *Grammar) makeSeq(seq x.Seq) *Rule {

	rule := g.NewRule()

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
