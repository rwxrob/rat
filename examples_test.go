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

func ExampleGrammar_Literal() {

	g := rat.NewGrammar()
	foo := g.Literal(`foo`)

	buf := []rune("barfoobazfo")

	g.Check(`'foo'`, buf, 3).Print()
	g.Check(`'foo'`, buf, 0).Print()
	g.Check(`'foo'`, buf, 9).Print()

	foo.Check(buf, 3).Print()
	foo.Check(buf, 0).Print()
	foo.Check(buf, 9).Print()

	// Output:
	// {"B":3,"E":6}
	// {"B":0,"E":0,"X":"expected: 'f'"}
	// {"B":9,"E":11,"X":"expected: 'o'"}
	// {"B":3,"E":6}
	// {"B":0,"E":0,"X":"expected: 'f'"}
	// {"B":9,"E":11,"X":"expected: 'o'"}

}

func ExampleGrammar_Sequence() {

	g := rat.NewGrammar()
	foo := g.Literal(`foo`)
	baz := g.Literal(`baz`)
	foobaz := g.Sequence(foo, baz)

	buf := []rune("barfoobazfoobut")

	g.Check(`'foo' 'baz'`, buf, 3).Print()
	g.Check(`'foo' 'baz'`, buf, 0).Print()
	g.Check(`'foo' 'baz'`, buf, 9).Print()

	foobaz.Check(buf, 3).Print()
	foobaz.Check(buf, 0).Print()
	foobaz.Check(buf, 9).Print()

	// Output:
	// {"B":3,"E":9}
	// {"B":0,"E":0,"X":"expected: 'f'"}
	// {"B":9,"E":13,"X":"expected: 'a'"}
	// {"B":3,"E":9}
	// {"B":0,"E":0,"X":"expected: 'f'"}
	// {"B":9,"E":13,"X":"expected: 'a'"}

}

func ExampleGrammar_OneOf() {

	g := new(rat.Grammar)
	foo := g.Literal(`foo`)
	bar := g.Literal(`bar`)
	baz := g.Literal(`baz`)
	oneof_foobarbaz := g.OneOf(foo, bar, baz)

	str := `foobarbaz`
	oneof_foobarbaz.Check([]rune(str), 0).Print()

	g.CheckString(`'foo' / 'bar' / 'baz'`, str, 3).Print()

	g.CheckString(`'foo' / 'bar' / 'baz'`, str, 2).Print()

	// Output:
	// {"B":0,"E":3}
	// {"B":3,"E":6}
	// {"B":2,"E":2,"X":"expected: 'foo' / 'bar' / 'baz'"}

}
