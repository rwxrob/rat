package rat

import (
	"fmt"
	"strconv"
	"sync"
)

type Grammar struct {
	sync.Map
	ResultNames    []string
	ResultNamesMap map[string]int

	anoncount int
	main      Rule
}

func NewGrammar() *Grammar {
	g := new(Grammar)
	g.ResultNames = []string{}
	g.ResultNamesMap = map[string]int{}
	return g
}

func (g Grammar) IsZero() bool { return g.main == nil }

func (g Grammar) String() string {
	if g.main == nil {
		return `<empty>`
	}
	return g.main.String()
}

func (g Grammar) Print() { fmt.Println(g) }

func (g *Grammar) Cache(r Rule) {
	key := r.String()
	if _, cached := g.Load(key); cached {
		return
	}
	g.Store(key, r)
	return
}

func (g *Grammar) Import(in *Grammar) {
	in.Range(func(k any, v any) bool {
		txt, isstring := k.(string)
		fn, isrule := v.(Rule)
		if isstring && isrule {
			g.Store(txt, fn)
		}
		return true
	})
}

func (g *Grammar) Pack(in ...any) {
	gg := Pack(in...)
	g.Import(gg)
}

func (g Grammar) Check(in any) Result {
	//TODO check against main
	return Result{}
}

func (g Grammar) Add(in any) Rule {
	var rule Rule

	switch v := in.(type) {

	case string:
		rule = Lit{v}
	case []rune:
		rule = Lit{string(v)}
	case []byte:
		rule = Lit{string(v)}

	case Rule:
		rule = v

	case func(r []rune, i int) Result:
		g.anoncount++
		txt := `Rule` + strconv.Itoa(g.anoncount)
		rule = arule{txt, v}
	case CheckFunc:
		g.anoncount++
		txt := `Rule` + strconv.Itoa(g.anoncount)
		rule = arule{txt, v}

	case rune: // int32 alias
		rule = Lit{string(v)}

	case int:
		rule = Lit{fmt.Sprintf(`%v`, v)}
	case int8:
		rule = Lit{fmt.Sprintf(`%v`, v)}
	case int16:
		rule = Lit{fmt.Sprintf(`%v`, v)}
	case int64:
		rule = Lit{fmt.Sprintf(`%v`, v)}

	case uint:
		rule = Lit{fmt.Sprintf(`%v`, v)}
	case uint8:
		rule = Lit{fmt.Sprintf(`%v`, v)}
	case uint16:
		rule = Lit{fmt.Sprintf(`%v`, v)}
	case uint32:
		rule = Lit{fmt.Sprintf(`%v`, v)}
	case uint64:
		rule = Lit{fmt.Sprintf(`%v`, v)}

	case float32:
		rule = Lit{fmt.Sprintf(`%v`, v)}
	case float64:
		rule = Lit{fmt.Sprintf(`%v`, v)}

	case bool:
		rule = Lit{fmt.Sprintf(`%v`, v)}
	}

	g.Cache(rule)
	return rule

}
