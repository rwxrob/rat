package rat_test

import (
	"fmt"

	"github.com/rwxrob/rat"
)

func ExampleResult() {

	r := []rune(`Something`)
	m := rat.Result{r, 3, 5, nil, nil}
	fmt.Println(m)

	// Output:
	// {"B":3,"E":5}
}

func ExampleRule() {

	LineFeed := func(r []rune, i int) rat.Result {
		var err error
		if r[i] != '\n' {
			err = fmt.Errorf(`expected line feed (\n)`)
		}
		return rat.Result{r, i, i + 1, nil, err}
	}

	buf := []rune("some\nthing\n")
	LineFeed(buf, 4).Print()
	LineFeed(buf, 1).Print()
	LineFeed(buf, 10).Print()

	// Output:
	// {"B":4,"E":5}
	// {"B":1,"E":2,"X":"expected line feed (\\n)"}
	// {"B":10,"E":11}

}

func ExampleLit() {

	Foo := rat.Lit(`foo`)
	buf := []rune("barfoobazfo")
	Foo(buf, 3).Print()
	Foo(buf, 0).Print()
	Foo(buf, 9).Print()

	// Output:
	// {"B":3,"E":6}
	// {"B":0,"E":0,"X":"expected literal \"foo\""}
	// {"B":9,"E":11,"X":"expected literal \"foo\""}

}

func ExampleSeq() {

	EndBlock := rat.Seq(rat.Lit("foo"), rat.Lit("baz"))

	buf := []rune("barfoobazfoobut")
	EndBlock(buf, 3).Print()
	EndBlock(buf, 0).Print()
	EndBlock(buf, 9).Print()

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
