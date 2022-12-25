package rat_test

import (
	"fmt"

	"github.com/rwxrob/rat"
)

func ExampleResult() {

	r := []rune(`Something`)
	m := rat.Result{R: r, B: 3, E: 5}
	fmt.Println(m)

	// Output:
	// {"B":3,"E":5}
}

func ExampleCheck() {

	LF := func(r []rune, i int) rat.Result {
		var err error
		if r[i] != '\n' {
			err = fmt.Errorf(`expected line feed (\n)`)
		}
		return rat.Result{R: r, B: i, E: i + 1, X: err}
	}

	buf := []rune("some\nthing\n")
	LF(buf, 4).Print()
	LF(buf, 1).Print()
	LF(buf, 10).Print()

	// Output:
	// {"B":4,"E":5}
	// {"B":1,"E":2,"X":"expected line feed (\\n)"}
	// {"B":10,"E":11}

}

func ExampleRule() {

	LF := rat.Rule{
		Text: "LF",
		Check: func(r []rune, i int) rat.Result {
			var err error
			if r[i] != '\n' {
				err = fmt.Errorf(`expected line feed (\n)`)
			}
			return rat.Result{R: r, B: i, E: i + 1, X: err}
		},
	}

	buf := []rune("some\nthing\n")
	LF.Check(buf, 4).Print()
	LF.Check(buf, 1).Print()
	LF.Check(buf, 10).Print()

	// Output:
	// {"B":4,"E":5}
	// {"B":1,"E":2,"X":"expected line feed (\\n)"}
	// {"B":10,"E":11}

}

func ExampleGrammar() {

	LF := rat.Rule{
		Text: "LF",
		Check: func(r []rune, i int) rat.Result {
			var err error
			if r[i] != '\n' {
				err = fmt.Errorf(`expected line feed (\n)`)
			}
			return rat.Result{R: r, B: i, E: i + 1, X: err}
		},
	}

	g := rat.NewGrammar(LF)

	buf := []rune("some\nthing\n")
	g.Check("LF", buf, 4).Print()
	g.Check("LF", buf, 1).Print()
	g.Check("LF", buf, 10).Print()

	// Output:
	// {"B":4,"E":5}
	// {"B":1,"E":2,"X":"expected line feed (\\n)"}
	// {"B":10,"E":11}

}

func ExampleLit() {

	g := rat.NewGrammar()
	g.Lit(`foo`)

	buf := []rune("barfoobazfo")
	g.Check(`'foo'`, buf, 3).Print()
	//Foo(buf, 0).Print()
	//Foo(buf, 9).Print()

	// Output:
	// {"B":3,"E":6}
	// {"B":0,"E":0,"X":"expected literal \"foo\""}
	// {"B":9,"E":11,"X":"expected literal \"foo\""}

}

/*
func ExampleSeq() {

	FooBaz := rat.Seq(rat.Lit("foo"), rat.Lit("baz"))

	buf := []rune("barfoobazfoobut")
	FooBaz(buf, 3).Print()
	FooBaz(buf, 0).Print()
	FooBaz(buf, 9).Print()

	// Output:
	// {"B":3,"E":9}
	// {"B":0,"E":0,"X":"expected literal \"foo\""}
	// {"B":9,"E":13,"X":"expected literal \"baz\""}

}

func ExampleFuncName() {

	fmt.Println(rat.FuncName(ExampleFuncName))
	fmt.Println(rat.FuncName(func() {}))

	// Output:
	// ExampleFuncName
	// func1
}

func ExampleErrOneOf() {

	Foo := rat.Rule{
		Text: `'foo'`,
		Func: func(r []rune, i int) rat.Result {
			return rat.Result{T: 1, R: r, B: i, E: i}
		},
	}

	Bar := Foo

	Baz := rat.Rule{
		Text: `'baz'`,
		Func: func(r []rune, i int) rat.Result { return rat.Result{} },
	}

	g := new(rat.Grammar)
	rule := g.OneOf(Foo, Bar, Baz)

	rule.Check(`foobarbaz`, 0).Print()

	// Output:
	// expected one of [foo foo func1]
}

func ExampleOneOf() {

	FooBarBaz := rat.OneOf(rat.Lit("foo"), rat.Lit("bar"), rat.Lit("baz"))
	buf := []rune("barfoobazfoobut")
	FooBarBaz(buf, 3).Print()
	FooBarBaz(buf, 0).Print()
	FooBarBaz(buf, 12).Print()

	// Output:
	// foo

}
*/
