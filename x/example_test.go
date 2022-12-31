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

	vals := []string{`foo`, `bar`}
	seq := x.Seq{vals}
	fmt.Println(seq)
	//	seq.Print()

	// Output:
	// x.Seq{"foo", "bar"}

}
