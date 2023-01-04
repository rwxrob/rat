package x_test

import (
	"fmt"

	"github.com/rwxrob/rat/x"
)

func ExampleSeq() {

	seq := x.Seq{`foo`, `bar`}
	fmt.Println(seq)
	//	seq.Print()

	// Output:
	// x.Seq{"foo", "bar"}

}

func ExampleSeq_from_Slice() {

	vals := []any{`foo`, `bar`}
	seq := x.Seq{vals}
	fmt.Println(seq)
	seq.Print()

	// Output:
	// x.Seq{"foo", "bar"}
	// x.Seq{"foo", "bar"}

}

func ExampleRule() {

	x.Rule{`foo`}.Print()
	x.Rule{`foo`, `Foo`}.Print()
	x.Rule{`foo`, `Foo`, 1}.Print()
	x.Rule{[]any{`foo`, `bar`}, `Foo`, 1}.Print()
	x.Rule{[]string{`foo`, `bar`}, `Foo`, 1}.Print()
	x.Rule{x.Seq{`foo`, `bar`}, `Foo`, 1}.Print()

	// Output:
	// "foo"
	// x.Rule{"foo", "Foo"}
	// x.Rule{"foo", "Foo", 1}
	// x.Rule{x.Seq{"foo", "bar"}, "Foo", 1}
	// x.Rule{x.Seq{"foo", "bar"}, "Foo", 1}
	// x.Rule{x.Seq{"foo", "bar"}, "Foo", 1}

}
