package rat_test

import (
	"strings"

	"github.com/rwxrob/rat"
)

func ExampleRule() {

	rule := new(rat.Rule)
	rule.Name = `Foo`
	rule.Text = `x.Lit{"foo"}`

	rule.Check = func(r []rune, i int) rat.Result {
		start := i
		if !strings.HasPrefix(string(r[i:]), `foo`) {
			return rat.Result{R: r, B: start, E: i, X: rat.ErrExpected{`foo`}}
		}
		return rat.Result{R: r, B: start, E: i}
	}

	buf := []rune(`foobar`)
	rule.Print()
	rule.Check(buf, 0).Print()
	rule.Check(buf, 1).Print()

	// Output:
	// Foo
	// {"B":0,"E":0}
	// {"B":1,"E":1,"X":"expected: foo"}

}

/*
func ExamplePack() {
	g := rat.Pack(`foo`, x.Any{2}, `bar`)
	g.Check(`fooisbar`).Print()
	// Output:
	// {"B":0,"E":8}
}
*/
