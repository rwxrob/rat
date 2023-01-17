package rat_test

import (
	"fmt"
	"unicode"

	"github.com/rwxrob/rat"
	"github.com/rwxrob/rat/x"
)

func ExampleFlatFunc_ByDepth() {

	r1 := rat.Result{N: `r1`, B: 1, E: 3}
	r2 := r1
	r1a := rat.Result{N: `r1a`, B: 1, E: 2}
	r1b := rat.Result{N: `r1b`, B: 2, E: 3}
	r1.C = []rat.Result{r1a, r1b}
	r2.N = `r2`

	root := rat.Result{
		N: `Root`, B: 1, E: 3, C: []rat.Result{r1, r2},
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
	r1b := rat.Result{N: `r1b`, B: 2, E: 3, C: []rat.Result{foo}}
	foo.I = 2
	r1.C = []rat.Result{r1a, r1b}
	r2.N = `r2`
	r2.C = []rat.Result{foo}
	foo.I = 3

	root := rat.Result{
		N: `Root`, B: 1, E: 3, C: []rat.Result{r1, r2, foo},
	}

	for _, result := range root.WithName(`foo`) {
		result.Print()
	}

	// Output:
	// {"N":"foo","I":1,"B":2,"E":3}
	// {"N":"foo","I":2,"B":2,"E":3}
	// {"N":"foo","I":3,"B":2,"E":3}

}

func ExamplePack_one() {

	g := rat.Pack(x.One{`foo`, `bar`})
	g.Print()

	g.Scan(`foobar`).PrintText()
	g.Scan(`foobar`).Print()

	g.Scan(`barfoo`).PrintText()
	g.Scan(`barfoo`).Print()

	g.Scan(`baz`).Print()

	// Output:
	// x.One{"foo", "bar"}
	// foo
	// {"B":0,"E":3,"C":[{"B":0,"E":3}],"R":"foobar"}
	// bar
	// {"B":0,"E":3,"C":[{"B":0,"E":3}],"R":"barfoo"}
	// {"B":0,"E":0,"X":"expected: x.One{\"foo\", \"bar\"}","R":"baz"}

}

func ExamplePach_lit_Boolean() {

	g := rat.Pack(true)
	g.Print()

	g.Scan(`true`).PrintText()
	g.Scan(`true`).Print()
	g.Scan(`false`).Print()
	g.Scan(`TRUE`).Print()

	// Output:
	// "true"
	// true
	// {"B":0,"E":4,"R":"true"}
	// {"B":0,"E":0,"X":"expected: t","R":"false"}
	// {"B":0,"E":0,"X":"expected: t","R":"TRUE"}

}

func ExamplePack_named() {

	g := rat.Pack(x.Name{`foo`, true})
	g.Print()

	g.Scan(`true`).Print()

	// Output:
	// x.Name{"foo", "true"}
	// {"N":"foo","B":0,"E":4,"R":"true"}

}
func ExamplePack_ref() {

	g := rat.Pack(x.Ref{`foo`})
	g.MakeRule(x.Name{`foo`, `foo`})
	g.Print()

	g.Scan(`foo`).Print()

	// Output:
	// x.Ref{"foo"}
	// {"N":"foo","B":0,"E":3,"R":"foo"}

}

func ExamplePack_one_Named() {

	one := x.One{`foo`, `bar`}
	Foo := x.Name{`Foo`, one}
	g := rat.Pack(Foo)
	g.Print()

	// foo
	g.Scan(`foobar`).Print()
	g.Scan(`foobar`).PrintText()

	// bar
	g.Scan(`barrr`).Print()
	g.Scan(`barrr`).PrintText()

	// bork
	g.Scan(`fobar`).Print()

	// Output:
	// x.Name{"Foo", x.One{"foo", "bar"}}
	// {"N":"Foo","B":0,"E":3,"C":[{"B":0,"E":3}],"R":"foobar"}
	// foo
	// {"N":"Foo","B":0,"E":3,"C":[{"B":0,"E":3}],"R":"barrr"}
	// bar
	// {"N":"Foo","B":0,"E":0,"X":"expected: x.One{\"foo\", \"bar\"}","R":"fobar"}

}

func ExamplePack_isfunc() {

	IsPrint := unicode.IsPrint

	g := rat.Pack(IsPrint)
	g.Print()

	g.Scan(`foo`).PrintText()
	g.Scan(`foo`).Print()

	g.Scan("\x00foo").Print()

	// Output:
	// x.Is{IsPrint}
	// f
	// {"B":0,"E":1,"R":"foo"}
	// {"B":0,"E":0,"X":"expected: x.Is{IsPrint}","R":"\x00foo"}

}

func ExamplePack_mmx() {

	g := rat.Pack(x.Mmx{1, 3, `foo`})
	g.Print()

	g.Scan(`foo`).PrintText()
	g.Scan(`foo`).Print()

	g.Scan(`foofoo`).PrintText()
	g.Scan(`foofoo`).Print()

	g.Scan(`foofoofoo`).PrintText()
	g.Scan(`foofoofoo`).Print()

	g.Scan(`foofoofoofoo`).PrintText()
	g.Scan(`foofoofoofoo`).Print()

	g.Scan(`barfoofoo`).Print()

	// Output:
	// x.Mmx{1, 3, "foo"}
	// foo
	// {"B":0,"E":3,"C":[{"B":0,"E":3}],"R":"foo"}
	// foofoo
	// {"B":0,"E":6,"C":[{"B":0,"E":3},{"B":3,"E":6}],"R":"foofoo"}
	// foofoofoo
	// {"B":0,"E":9,"C":[{"B":0,"E":3},{"B":3,"E":6},{"B":6,"E":9}],"R":"foofoofoo"}
	// foofoofoo
	// {"B":0,"E":9,"C":[{"B":0,"E":3},{"B":3,"E":6},{"B":6,"E":9},{"B":9,"E":12}],"R":"foofoofoofoo"}
	// {"B":0,"E":0,"X":"expected: x.Mmx{1, 3, \"foo\"}","R":"barfoofoo"}

}

func ExamplePack_end() {

	g := rat.Pack(x.Any{2}, x.End{})
	g.Print()

	g.Scan(`fo`).PrintText()
	g.Scan(`fo`).Print()

	g.Scan(`foo`).Print()

	// Output:
	// x.Seq{x.Any{2}, x.End{}}
	// fo
	// {"B":0,"E":2,"C":[{"B":0,"E":2},{"B":2,"E":2}],"R":"fo"}
	// {"B":0,"E":2,"X":"expected: x.End{}","C":[{"B":0,"E":2},{"B":2,"E":2,"X":"expected: x.End{}"}],"R":"foo"}

}

func ExamplePack_opt() {

	g := rat.Pack(x.Opt{`foo`})
	g.Print()

	g.Scan(`foo`).PrintText()
	g.Scan(`foo`).Print()

	g.Scan(`bar`).PrintText()
	g.Scan(`bar`).Print()

	// Output:
	// x.Opt{"foo"}
	// foo
	// {"B":0,"E":3,"R":"foo"}
	//
	// {"B":0,"E":0,"R":"bar"}

}

func ExamplePack_rep() {

	g := rat.Pack(x.Rep{2, `foo`})
	g.Print()

	g.Scan(`foofoofoo`).PrintText()
	g.Scan(`foofoofoo`).Print()

	g.Scan(`foobar`).Print()

	// Output:
	// x.Rep{2, "foo"}
	// foofoo
	// {"B":0,"E":6,"C":[{"B":0,"E":3},{"B":3,"E":6}],"R":"foofoofoo"}
	// {"B":0,"E":3,"X":"expected: x.Rep{2, \"foo\"}","C":[{"B":0,"E":3}],"R":"foobar"}

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

func ExampleMakeAny() {

	g := new(rat.Grammar)
	rule := g.MakeAny(x.Any{3})
	rule.Check([]rune(`..`), 0).Print()
	rule.Check([]rune(`...`), 0).Print()
	rule.Check([]rune(`....`), 0).Print()
	rule.Check([]rune(`....`), 2).Print()
	fmt.Println(g.Rules[`x.Any{3}`].Name)
	fmt.Println(g.Rules[`x.Any{3}`].Text)

	//Output:
	// {"B":0,"E":1,"X":"expected: x.Any{3}","R":".."}
	// {"B":0,"E":3,"R":"..."}
	// {"B":0,"E":3,"R":"...."}
	// {"B":2,"E":3,"X":"expected: x.Any{3}","R":"...."}
	// x.Any{3}
	// x.Any{3}

}

func ExampleMakeLit() {

	g := new(rat.Grammar)
	foo := g.MakeLit(`foo`)
	oo := g.MakeLit(`oo`)
	foo.Check([]rune(`foo`), 0).Print()
	foo.Check([]rune(`fooo`), 0).Print()
	foo.Check([]rune(`fo`), 0).Print()
	oo.Check([]rune(`fooo`), 0).Print()
	oo.Check([]rune(`fooo`), 1).Print()
	oo.Check([]rune(`fooo`), 2).Print()
	fmt.Println(g.Rules[`foo`].Name)
	fmt.Println(g.Rules[`foo`].Text)
	fmt.Println(g.Rules[`oo`].Name)
	fmt.Println(g.Rules[`oo`].Text)

	//Output:
	// {"B":0,"E":3,"R":"foo"}
	// {"B":0,"E":3,"R":"fooo"}
	// {"B":0,"E":2,"X":"expected: o","R":"fo"}
	// {"B":0,"E":0,"X":"expected: o","R":"fooo"}
	// {"B":1,"E":3,"R":"fooo"}
	// {"B":2,"E":4,"R":"fooo"}
	// foo
	// "foo"
	// oo
	// "oo"

}
