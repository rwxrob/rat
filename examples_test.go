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
	// x.One{x.Str{"foo"}, x.Str{"bar"}}
	// foo
	// {"B":0,"E":3,"C":[{"B":0,"E":3}],"R":"foobar"}
	// bar
	// {"B":0,"E":3,"C":[{"B":0,"E":3}],"R":"barfoo"}
	// {"B":0,"E":0,"X":"expected: x.One{x.Str{\"foo\"}, x.Str{\"bar\"}}","R":"baz"}

}

func ExamplePach_lit_Boolean() {

	g := rat.Pack(true)
	g.Print()

	g.Scan(`true`).PrintText()
	g.Scan(`true`).Print()
	g.Scan(`false`).Print()
	g.Scan(`TRUE`).Print()

	// Output:
	// x.Str{"true"}
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
	// x.Name{"foo", x.Str{"true"}}
	// {"N":"foo","B":0,"E":4,"R":"true"}

}
func ExamplePack_ref() {

	g := rat.Pack(x.Ref{`Foo`})
	g.MakeRule(x.Name{`Foo`, `foo`})
	g.Print()

	g.Rules[`x.Str{"foo"}`].Print()
	g.Rules[`Foo`].Print()

	g.Scan(`foo`).Print()
	g.Rules[`x.Str{"foo"}`].Scan(`foo`).Print()

	// Output:
	// x.Ref{"Foo"}
	// x.Str{"foo"}
	// x.Name{"Foo", x.Str{"foo"}}
	// {"N":"Foo","B":0,"E":3,"R":"foo"}
	// {"B":0,"E":3,"R":"foo"}

}

func ExamplePack_save() {

	g := new(rat.Grammar).Init()
	g.MakeRule(x.Name{`Post`, x.Mmx{3, 8, '`'}})
	g.Pack(x.Name{`Fenced`, x.Seq{x.Save{`Post`}, x.To{x.Val{`Post`}}, x.Val{`Post`}}})
	g.Print()

	// one step at a time
	g.Rules[`x.Save{"Post"}`].Scan("````").Print()
	g.Rules[`Post`].Scan("````").Print()
	g.Rules[`x.Val{"Post"}`].Scan("````````").Print()
	g.Rules[`x.To{x.Val{"Post"}}`].Scan("....``````").Print()

	// combined
	g.Scan("```.......`````").PrintText()
	g.Scan("```.......`````").Print()

	// Output:
	// x.Name{"Fenced", x.Seq{x.Save{"Post"}, x.To{x.Val{"Post"}}, x.Val{"Post"}}}
	// {"N":"Post","B":0,"E":4,"C":[{"B":0,"E":1},{"B":1,"E":2},{"B":2,"E":3},{"B":3,"E":4}],"R":"````"}
	// {"N":"Post","B":0,"E":4,"C":[{"B":0,"E":1},{"B":1,"E":2},{"B":2,"E":3},{"B":3,"E":4}],"R":"````"}
	// {"B":0,"E":4,"R":"````````"}
	// {"B":0,"E":4,"R":"....``````"}
	// ```.......```
	// {"N":"Fenced","B":0,"E":13,"C":[{"N":"Post","B":0,"E":3,"C":[{"B":0,"E":1},{"B":1,"E":2},{"B":2,"E":3}]},{"B":3,"E":10},{"B":10,"E":13}],"R":"```.......`````"}

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
	// x.Name{"Foo", x.One{x.Str{"foo"}, x.Str{"bar"}}}
	// {"N":"Foo","B":0,"E":3,"C":[{"B":0,"E":3}],"R":"foobar"}
	// foo
	// {"N":"Foo","B":0,"E":3,"C":[{"B":0,"E":3}],"R":"barrr"}
	// bar
	// {"N":"Foo","B":0,"E":0,"X":"expected: x.One{x.Str{\"foo\"}, x.Str{\"bar\"}}","R":"fobar"}

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
	// x.Mmx{1, 3, x.Str{"foo"}}
	// foo
	// {"B":0,"E":3,"C":[{"B":0,"E":3}],"R":"foo"}
	// foofoo
	// {"B":0,"E":6,"C":[{"B":0,"E":3},{"B":3,"E":6}],"R":"foofoo"}
	// foofoofoo
	// {"B":0,"E":9,"C":[{"B":0,"E":3},{"B":3,"E":6},{"B":6,"E":9}],"R":"foofoofoo"}
	// foofoofoo
	// {"B":0,"E":9,"C":[{"B":0,"E":3},{"B":3,"E":6},{"B":6,"E":9},{"B":9,"E":12}],"R":"foofoofoofoo"}
	// {"B":0,"E":0,"X":"expected: x.Mmx{1, 3, x.Str{\"foo\"}}","R":"barfoofoo"}

}

func ExamplePack_mn1() {

	g := rat.Pack(x.Mn1{`foo`})
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
	// x.Mn1{x.Str{"foo"}}
	// foo
	// {"B":0,"E":3,"C":[{"B":0,"E":3}],"R":"foo"}
	// foofoo
	// {"B":0,"E":6,"C":[{"B":0,"E":3},{"B":3,"E":6}],"R":"foofoo"}
	// foofoofoo
	// {"B":0,"E":9,"C":[{"B":0,"E":3},{"B":3,"E":6},{"B":6,"E":9}],"R":"foofoofoo"}
	// foofoofoofoo
	// {"B":0,"E":12,"C":[{"B":0,"E":3},{"B":3,"E":6},{"B":6,"E":9},{"B":9,"E":12}],"R":"foofoofoofoo"}
	// {"B":0,"E":0,"X":"expected: x.Mn1{x.Str{\"foo\"}}","R":"barfoofoo"}

}

func ExamplePack_mn0() {

	g := rat.Pack(x.Mn0{`foo`})
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
	// x.Mn0{x.Str{"foo"}}
	// foo
	// {"B":0,"E":3,"C":[{"B":0,"E":3}],"R":"foo"}
	// foofoo
	// {"B":0,"E":6,"C":[{"B":0,"E":3},{"B":3,"E":6}],"R":"foofoo"}
	// foofoofoo
	// {"B":0,"E":9,"C":[{"B":0,"E":3},{"B":3,"E":6},{"B":6,"E":9}],"R":"foofoofoo"}
	// foofoofoofoo
	// {"B":0,"E":12,"C":[{"B":0,"E":3},{"B":3,"E":6},{"B":6,"E":9},{"B":9,"E":12}],"R":"foofoofoofoo"}
	// {"B":0,"E":0,"R":"barfoofoo"}

}

func ExamplePack_min() {

	g := rat.Pack(x.Min{2, `foo`})
	g.Print()

	g.Scan(`foofoo`).PrintText()
	g.Scan(`foofoo`).Print()

	g.Scan(`foofoofoo`).PrintText()
	g.Scan(`foofoofoo`).Print()

	g.Scan(`foo`).Print()
	g.Scan(`barfoofoo`).Print()

	// Output:
	// x.Min{2, x.Str{"foo"}}
	// foofoo
	// {"B":0,"E":6,"C":[{"B":0,"E":3},{"B":3,"E":6}],"R":"foofoo"}
	// foofoofoo
	// {"B":0,"E":9,"C":[{"B":0,"E":3},{"B":3,"E":6},{"B":6,"E":9}],"R":"foofoofoo"}
	// {"B":0,"E":3,"X":"expected: x.Min{2, x.Str{\"foo\"}}","C":[{"B":0,"E":3}],"R":"foo"}
	// {"B":0,"E":0,"X":"expected: x.Min{2, x.Str{\"foo\"}}","R":"barfoofoo"}

}

func ExamplePack_max() {

	g := rat.Pack(x.Max{2, `foo`})
	g.Print()

	g.Scan(`foo`).PrintText()
	g.Scan(`foo`).Print()

	g.Scan(`foofoo`).PrintText()
	g.Scan(`foofoo`).Print()

	g.Scan(`barfoofoo`).PrintText()
	g.Scan(`barfoofoo`).Print()

	g.Scan(`foofoofoo`).Print()

	// Output:
	// x.Max{2, x.Str{"foo"}}
	// foo
	// {"B":0,"E":3,"C":[{"B":0,"E":3}],"R":"foo"}
	// foofoo
	// {"B":0,"E":6,"C":[{"B":0,"E":3},{"B":3,"E":6}],"R":"foofoo"}
	//
	// {"B":0,"E":0,"R":"barfoofoo"}
	// {"B":0,"E":6,"X":"expected: x.Max{2, x.Str{\"foo\"}}","C":[{"B":0,"E":3},{"B":3,"E":6}],"R":"foofoofoo"}

}

func ExamplePack_pos() {

	g := rat.Pack(x.See{`foo`})
	g.Print()

	g.Scan(`fooooo`).PrintText()
	g.Scan(`fooooo`).Print()

	g.Scan(`fo`).Print()
	g.Scan(`bar`).Print()

	// Output:
	// x.See{x.Str{"foo"}}
	//
	// {"B":0,"E":0,"R":"fooooo"}
	// {"B":0,"E":0,"X":"expected: x.See{x.Str{\"foo\"}}","R":"fo"}
	// {"B":0,"E":0,"X":"expected: x.See{x.Str{\"foo\"}}","R":"bar"}

}

func ExamplePack_neg() {

	g := rat.Pack(x.Not{`foo`})
	g.Print()

	g.Scan(`fo`).PrintText()
	g.Scan(`fo`).Print()

	g.Scan(`bar`).PrintText()
	g.Scan(`bar`).Print()

	g.Scan(`fooooo`).Print()

	// Output:
	// x.Not{x.Str{"foo"}}
	//
	// {"B":0,"E":0,"R":"fo"}
	//
	// {"B":0,"E":0,"R":"bar"}
	// {"B":0,"E":0,"X":"expected: x.Not{x.Str{\"foo\"}}","R":"fooooo"}

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

func ExamplePack_rng() {

	g := rat.Pack(x.Rng{'😀', '🙏'})
	g.Print()

	g.Scan(`🙉`).PrintText()
	g.Scan(`🙉`).Print()

	g.Scan(`🚆`).Print()

	// Output:
	// x.Rng{'😀', '🙏'}
	// 🙉
	// {"B":0,"E":1,"R":"🙉"}
	// {"B":0,"E":0,"X":"expected: x.Rng{'😀', '🙏'}","R":"🚆"}

}

func ExamplePack_opt() {

	g := rat.Pack(x.Opt{`foo`})
	g.Print()

	g.Scan(`foo`).PrintText()
	g.Scan(`foo`).Print()

	g.Scan(`bar`).PrintText()
	g.Scan(`bar`).Print()

	// Output:
	// x.Opt{x.Str{"foo"}}
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
	// x.Rep{2, x.Str{"foo"}}
	// foofoo
	// {"B":0,"E":6,"C":[{"B":0,"E":3},{"B":3,"E":6}],"R":"foofoofoo"}
	// {"B":0,"E":3,"X":"expected: x.Rep{2, x.Str{\"foo\"}}","C":[{"B":0,"E":3}],"R":"foobar"}

}

func ExamplePack_to() {

	g := rat.Pack(x.To{`foo`})
	g.Print()

	g.Scan(`...foo`).PrintText()
	g.Scan(`...foo`).Print()
	g.Scan(`foofoo`).PrintText()
	g.Scan(`foofoo`).Print()
	g.Scan(`.foofo`).PrintText()
	g.Scan(`.foofo`).Print()

	g.Scan(`...fo`).Print()
	g.Scan(`...bar`).Print()

	// Output:
	// x.To{x.Str{"foo"}}
	// ...
	// {"B":0,"E":3,"R":"...foo"}
	//
	// {"B":0,"E":0,"R":"foofoo"}
	// .
	// {"B":0,"E":1,"R":".foofo"}
	// {"B":0,"E":5,"X":"expected: x.To{x.Str{\"foo\"}}","R":"...fo"}
	// {"B":0,"E":6,"X":"expected: x.To{x.Str{\"foo\"}}","R":"...bar"}

}

func ExampleMakeAny() {

	g := new(rat.Grammar).Init()
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

func ExampleMakeStr() {

	g := new(rat.Grammar).Init()
	foo := g.MakeStr(`foo`)
	foo.Print()
	oo := g.MakeStr(`oo`)
	oo.Print()
	foo.Check([]rune(`foo`), 0).Print()
	foo.Check([]rune(`fooo`), 0).Print()
	foo.Check([]rune(`fo`), 0).Print()
	oo.Check([]rune(`fooo`), 0).Print()
	oo.Check([]rune(`fooo`), 1).Print()
	oo.Check([]rune(`fooo`), 2).Print()

	for k, v := range g.Rules {
		fmt.Printf("key: %q name: %q text: %q\n", k, v.Name, v.Text)
	}

	// Unordered Output:
	// x.Str{"foo"}
	// x.Str{"oo"}
	// {"B":0,"E":3,"R":"foo"}
	// {"B":0,"E":3,"R":"fooo"}
	// {"B":0,"E":2,"X":"expected: o","R":"fo"}
	// {"B":0,"E":0,"X":"expected: o","R":"fooo"}
	// {"B":1,"E":3,"R":"fooo"}
	// {"B":2,"E":4,"R":"fooo"}
	// key: "x.Str{\"foo\"}" name: "x.Str{\"foo\"}" text: "x.Str{\"foo\"}"
	// key: "x.Str{\"oo\"}" name: "x.Str{\"oo\"}" text: "x.Str{\"oo\"}"

}
