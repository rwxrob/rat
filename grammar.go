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

// SetCheckRule sets the primary entry rule used when g.Check is called.
func (g *Grammar) SetCheckRule(rule *Rule) *Grammar {
	g.rule = rule
	return g
}

func (g *Grammar) Check(in any) Result {
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

// Pack is shorthand for g.SetCheckRule(g.MakeRule(x.Seq{[]seq})).
func (g *Grammar) Pack(seq ...any) *Grammar {
	return g.SetCheckRule(g.MakeRule(x.Seq(seq)))
}

// MakeRule fulfills the MakeRule interface by returning the same Rule
// created from a rat/x.Seq{in}.
func (g *Grammar) MakeRule(in any) *Rule {
	switch v := in.(type) {

	case string:
		return g.makeLit(v)
	case []rune:
		return g.makeLit(string(v))
	case []byte:
		return g.makeLit(string(v))
	case rune:
		return g.makeLit(string(v))

	case x.Seq:
		return g.makeSeq(v)

	case x.Any:
		return g.makeAny(v)

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

func (g *Grammar) makeLit(in string) *Rule {
	rule, has := g.rules[in]
	if has {
		return rule
	}
	rule = new(Rule)
	rule.Name = in
	rule.Text = fmt.Sprintf(`%q`, in)
	rule.Check = func(r []rune, i int) Result {
		var err error
		var n int
		e := i
		runes := []rune(in)
		for e < len(r) && n < len(runes) {
			if r[e] != runes[n] {
				err = ErrExpected{r[e]}
				break
			}
			e++
			n++
		}
		if n < len(runes) {
			err = ErrExpected{string(runes[n])}
		}
		return Result{R: r, B: i, E: e, X: err}
	}
	g.addRule(rule)
	return rule
}

func (g *Grammar) makeAny(in x.Any) *Rule {
	rule := new(Rule)
	if len(in) != 1 {
		return nil
	}
	n, isint := in[0].(int)
	if !isint {
		return nil
	}
	rule.Name = `x.Any{` + strconv.Itoa(n) + `}`
	rule.Text = rule.Name

	// check cache

	if r, has := g.rules[rule.Name]; has {
		return r
	}

	rule.Check = func(r []rune, i int) Result {
		start := i
		if i+n > len(r) {
			return Result{R: r, B: start, E: len(r) - 1, X: ErrExpected{rule.Name}}
		}
		return Result{R: r, B: start, E: i + n}
	}

	g.addRule(rule)
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
		e := i
		results := []Result{}
		for _, rule := range rules {
			result := rule.Check(r, i)
			e = result.E
			results = append(results, result)
			if result.X != nil {
				break
			}
		}
		return Result{R: r, B: start, E: e, S: results}
	}

	return rule
}
