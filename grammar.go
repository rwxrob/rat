package rat

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rwxrob/rat/x"
)

// additive only
type Grammar struct {
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

// MakeRule fulfills the MakeRule interface by returning the same Rule
// created from a rat/x.Seq{in}.
func (g *Grammar) MakeRule(in any) *Rule {

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

	// types as expressions
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

func (g *Grammar) newRule() *Rule {
	rule := new(Rule)
	g.addRule(rule)
	return rule
}

func (g *Grammar) addRule(rule *Rule) {
	if rule.ID == 0 {
		g.rulenum++
		rule.ID = g.rulenum
	}
	if rule.Name == "" {
		rule.Name = DefaultRuleName + strconv.Itoa(rule.ID)
	}
	if g.rules == nil {
		g.rules = map[string]*Rule{}
	}
	g.rules[rule.Name] = rule
}

func (g *Grammar) makeName(in x.Name) *Rule {
	name, isstring := in[0].(string)
	if !isstring {
		panic(x.UsageName)
	}
	rule, has := g.rules[name]
	if has {
		return rule
	}

	// TODO
	return rule
}

func (g *Grammar) makeLit(in string) *Rule {
	rule, has := g.rules[in]
	if has {
		return rule
	}
	rule = g.newRule()
	rule.Name = in
	rule.Text = fmt.Sprintf(`%q`, in)
	g.addRule(rule)

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

	n, isint := in[0].(int)
	if !isint || len(in) != 1 {
		panic(fmt.Sprintf(_ErrArgs, in))
	}

	name := `x.Any{` + strconv.Itoa(n) + `}`
	if r, have := g.rules[name]; have {
		return r
	}

	rule := new(Rule)
	rule.Name = name
	rule.Text = name
	g.addRule(rule)

	rule.Check = func(r []rune, i int) Result {
		start := i
		if i+n > len(r) {
			return Result{R: r, B: start, E: len(r) - 1, X: ErrExpected{rule.Name}}
		}
		return Result{R: r, B: start, E: i + n}
	}

	return rule
}

var DefaultRuleName = `Rule`

func (g *Grammar) makeSeq(seq x.Seq) *Rule {

	rule := g.newRule()

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
