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
				err = fmt.Errorf(`expected: LF`)
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
	// {"B":1,"E":2,"X":"expected: LF"}
	// {"B":10,"E":11}

}

func ExampleGrammar_Lit() {

	g := rat.NewGrammar()
	g.Lit(`foo`)

	buf := []rune("barfoobazfo")
	g.Check(`'foo'`, buf, 3).Print()
	g.Check(`'foo'`, buf, 0).Print()
	g.Check(`'foo'`, buf, 9).Print()

	// Output:
	// {"B":3,"E":6}
	// {"B":0,"E":0,"X":"expected: 'f'"}
	// {"B":9,"E":11,"X":"expected: 'o'"}

}

func ExampleGrammar_Seq() {

	g := rat.NewGrammar()
	foo := g.Lit(`foo`)
	baz := g.Lit(`baz`)
	g.Seq(foo, baz)

	buf := []rune("barfoobazfoobut")

	g.Check(`('foo' 'baz')`, buf, 3).Print()
	g.Check(`('foo' 'baz')`, buf, 0).Print()
	g.Check(`('foo' 'baz')`, buf, 9).Print()

	// Output:
	// {"B":3,"E":9}
	// {"B":0,"E":0,"X":"expected: 'f'"}
	// {"B":9,"E":13,"X":"expected: 'a'"}

}

func ExampleGrammar_One() {

	g := new(rat.Grammar)
	foo := g.Lit(`foo`)
	bar := g.Lit(`bar`)
	baz := g.Lit(`baz`)
	one_foobarbaz := g.One(foo, bar, baz)

	str := `foobarbaz`
	one_foobarbaz.Check([]rune(str), 0).Print()
	g.CheckString(`('foo' / 'bar' / 'baz')`, str, 3).Print()
	g.CheckString(`('foo' / 'bar' / 'baz')`, str, 2).Print()

	// Output:
	// {"B":0,"E":3}
	// {"B":3,"E":6}
	// {"B":2,"E":2,"X":"expected: ('foo' / 'bar' / 'baz')"}

}

func ExampleToPEGN() {

	fmt.Printf("%q\n", rat.ToPEGN("some\tthing\nuh\rwhat\r\nsmile😈"))
	fmt.Printf("%q\n", rat.ToPEGN("some"))

	// Output:
	// "('some' TAB 'thing' LF 'uh' CR 'what' CR LF 'smile' x1f608)"
	// "'some'"

}
