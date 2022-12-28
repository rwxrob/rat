package rat_test

import (
	"fmt"

	"github.com/rwxrob/rat"
)

func ExampleResult() {

	names := []string{`Unknown`, `First`, `LetterM`}

	r := []rune(`Something`)
	m := rat.Result{N: 2, R: r, B: 3, E: 5}
	fmt.Println(m)
	fmt.Println(names[m.N])

	mx := rat.Result{R: r, B: 4, E: 4, X: fmt.Errorf(`bork`)}
	fmt.Println(mx)
	fmt.Println(names[mx.N])

	// Output:
	// {"N":2,"B":3,"E":5}
	// LetterM
	// {"B":4,"E":4,"X":"bork"}
	// Unknown
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

	g := rat.Pack(LF)

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

	g := new(rat.Grammar)
	g.RuleLiteral(`foo`)

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

	g := new(rat.Grammar)
	foo := g.RuleLiteral(`foo`)
	baz := g.RuleLiteral(`baz`)
	g.RuleSequence(foo, baz)

	buf := []rune("barfoobazfoobut")

	g.Check(`('foo' 'baz')`, buf, 3).Print()
	g.Check(`('foo' 'baz')`, buf, 0).Print()
	g.Check(`('foo' 'baz')`, buf, 9).Print()

	// Output:
	// {"B":3,"E":9}
	// {"B":0,"E":0,"X":"expected: 'f'"}
	// {"B":9,"E":13,"X":"expected: 'a'"}

}

func ExampleGrammar_RuleOneOf() {

	g := new(rat.Grammar)
	foo := g.RuleLiteral(`foo`)
	bar := g.RuleLiteral(`bar`)
	baz := g.RuleLiteral(`baz`)
	one_foobarbaz := g.RuleOneOf(foo, bar, baz)

	str := `foobarbaz`
	one_foobarbaz.Check([]rune(str), 0).Print()
	g.CheckString(`('foo' / 'bar' / 'baz')`, str, 3).Print()
	g.CheckString(`('foo' / 'bar' / 'baz')`, str, 2).Print()

	// Output:
	// {"B":0,"E":3}
	// {"B":3,"E":6}
	// {"B":2,"E":2,"X":"expected: ('foo' / 'bar' / 'baz')"}

}

func ExamplePack_One() {

	g := rat.Pack(rat.One{`foo`, `bar`, `baz`})

	str := `foobarbaz`
	g.CheckString(`('foo' / 'bar' / 'baz')`, str, 3).Print()
	g.CheckString(`('foo' / 'bar' / 'baz')`, str, 2).Print()

	// Output:
	// {"B":0,"E":3}
	// {"B":3,"E":6}
	// {"B":2,"E":2,"X":"expected: ('foo' / 'bar' / 'baz')"}

}

/*

func ExampleToPEGN() {

	fmt.Printf("%q\n", rat.ToPEGN("some\tthing\nuh\rwhat\r\nsmile😈"))
	fmt.Printf("%q\n", rat.ToPEGN("some"))

	// Output:
	// "('some' TAB 'thing' LF 'uh' CR 'what' CR LF 'smile' x1f608)"
	// "'some'"

}
*/
