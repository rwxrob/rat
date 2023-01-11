package rat_test

import (
	"fmt"

	"github.com/rwxrob/rat"
)

func ExampleFlatFunc_ByDepth() {

	r1 := rat.Result{N: `r1`, B: 1, E: 3}
	r2 := r1
	r1a := rat.Result{N: `r1a`, B: 1, E: 2}
	r1b := rat.Result{N: `r1b`, B: 2, E: 3}
	r1.C = rat.Results{r1a, r1b}
	r2.N = `r2`

	root := rat.Result{
		N: `Root`, B: 1, E: 3, C: rat.Results{r1, r2},
	}

	for _, result := range rat.ByDepth(root) {
		fmt.Println(result.N)
	}

	// Output:
	// Root
	// r1
	// r1a
	// r1b
	// r2

}

func ExampleResult_WithName() {

	foo := rat.Result{N: `foo`, I: 1, B: 2, E: 3}
	r1 := rat.Result{N: `r1`, B: 1, E: 3}
	r2 := r1
	r1a := rat.Result{N: `r1a`, B: 1, E: 2}
	r1b := rat.Result{N: `r1b`, B: 2, E: 3, C: rat.Results{foo}}
	foo.I = 2
	r1.C = rat.Results{r1a, r1b}
	r2.N = `r2`
	r2.C = rat.Results{foo}
	foo.I = 3

	root := rat.Result{
		N: `Root`, B: 1, E: 3, C: rat.Results{r1, r2, foo},
	}

	for _, result := range root.WithName(`foo`) {
		result.Print()
	}

	// Output:
	// {"N":"foo","I":1,"B":2,"E":3}
	// {"N":"foo","I":2,"B":2,"E":3}
	// {"N":"foo","I":3,"B":2,"E":3}

}

/*
func ExampleRule() {

	rule := new(rat.Rule)
	rule.Name = `Foo`
	rule.ID = 1
	rule.Text = `x.Rule{ 1, "Foo", x.Lit{"foo"} }`

	rule.Check = func(r []rune, i int) rat.Result {
		start := i
		if !strings.HasPrefix(string(r[i:]), `foo`) {
			return rat.Result{T: rule.ID, R: r, B: start, E: i, X: rat.ErrExpected{`foo`}}
		}
		return rat.Result{T: rule.ID, R: r, B: start, E: i}
	}

	buf := []rune(`foobar`)
	rule.Print()
	rule.Check(buf, 0).Print()
	rule.Check(buf, 1).Print()

	// Output:
	// x.Rule{ 1, "Foo", x.Lit{"foo"} }
	// {"T":1,"B":0,"E":0,"R":"foobar"}
	// {"T":1,"B":1,"E":1,"X":"expected: foo","R":"foobar"}

}

func ExamplePack() {

	g := rat.Pack(`foo`, x.Any{2}, `bar`, `foo`)
	g.Print()
	//res := g.Check(`fooisbarfoo`)
	//res.Print()

	// Output:
	// x.Seq{x.Rule{"Foo", "foo"}, x.Any{2}, "bar", x.Ref{"Foo"}}
	// {"T":1,"B":0,"E":11,"C":[{"T":2,"B":0,"E":3},{"T":3,"B":3,"E":5},{"T":4,"B":5,"E":8},{"T":2,"B":8,"E":11}],"R":"fooisbarfoo"}

}

func ExamplePack_ref() {

	g := rat.Pack(x.Rule{`Foo`, `foo`}, x.Any{2}, `bar`, x.Ref{`Foo`})
	g.Print()
	res := g.Check(`fooisbarfoo`)
	res.Print()

	// Output:
	// some
}
*/
