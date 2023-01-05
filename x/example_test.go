package x_test

import (
	"fmt"
	"unicode"

	"github.com/rwxrob/rat/x"
)

func ExampleString() {

	smile := int32('\u263A')
	types := []any{
		`string`, []byte(`bytes as string`), []rune(`runes as string`), 'x', '😀',
		true, false, -127, -32767, -9223372036854775808, -9785, int8(127), uint8(255),
		3.141592653589793238, int16(32767), int64(9223372036854775807), int(smile),
		'\x00', int32(9785), smile, '👩',
	}
	for _, it := range types {
		fmt.Println(x.String(it))
	}

	// Output:
	// "string"
	// "bytes as string"
	// "runes as string"
	// "x"
	// "😀"
	// "true"
	// "false"
	// "-127"
	// "-32767"
	// "-9223372036854775808"
	// "-9785"
	// "127"
	// "255"
	// "3.141592653589793"
	// "32767"
	// "9223372036854775807"
	// "9786"
	// "\x00"
	// "☹"
	// "☺"
	// "👩"

}

func ExampleRule() {

	x.Rule{`foo`}.Print()
	x.Rule{`foo`, `Foo`}.Print()
	x.Rule{`foo`, `Foo`, 1}.Print()
	x.Rule{[]any{`foo`, `bar`}, `Foo`, 1}.Print()
	x.Rule{[]string{`foo`, `bar`}, `Foo`, 1}.Print()
	x.Rule{x.Seq{`foo`, `bar`}, `Foo`, 1}.Print()
	x.Rule{`foo`, `Foo`, 1, false}.Print()

	// Output:
	// "foo"
	// x.Rule{"foo", "Foo"}
	// x.Rule{"foo", "Foo", 1}
	// x.Rule{x.Seq{"foo", "bar"}, "Foo", 1}
	// x.Rule{x.Seq{"foo", "bar"}, "Foo", 1}
	// x.Rule{x.Seq{"foo", "bar"}, "Foo", 1}
	// "%!USAGE: x.Rule{rule,name,id} or x.Rule{rule,name}"

}

// -------------------------------- Ref -------------------------------

func ExampleRef() {

	ref := x.Ref{`foo`}
	ref.Print()

	ref[0] = true
	ref.Print()

	// Output:
	// x.Ref{"foo"}
	// "%!USAGE: x.Ref{name}"

}

// -------------------------------- Is --------------------------------

func aclass(r rune) bool { return true }

func myIsPrint(r rune) bool { return unicode.IsPrint(r) }

func ExampleIs() {

	foo := aclass              // named function in the x_test package
	another := unicode.IsPrint // retains IsPrint original name
	myIsPrint := myIsPrint     // full wraps in own function to retain name

	x.Is{foo}.Print()
	x.Is{another}.Print()
	x.Is{myIsPrint}.Print()

	// anonymous functions not allowed for classes
	x.Is{func(r rune) bool { return true }}.Print()
	anon := func(r rune) bool { return true }
	x.Is{anon}.Print()

	// Output:
	// x.Is{aclass}
	// x.Is{IsPrint}
	// x.Is{myIsPrint}
	// "%!USAGE: x.Is{namedfunc}"
	// "%!USAGE: x.Is{namedfunc}"

}

/*

// -------------------------------- Seq -------------------------------

func ExampleSeq() {

	seq := x.Seq{`foo`, `bar`}
	seq.Print()

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

// -------------------------------- One -------------------------------

func ExampleOne() {

	one := x.One{}
	one.Print()

	// Output:
	// x.One{}

}

// -------------------------------- Opt -------------------------------

func ExampleOpt() {

	opt := x.Opt{}
	opt.Print()

	// Output:
	// x.Opt{}

}

// -------------------------------- Lit -------------------------------

func ExampleLit() {

	lit := x.Lit{}
	lit.Print()

	// Output:
	// x.Lit{}

}

// -------------------------------- Mn1 -------------------------------

func ExampleMn1() {

	mn1 := x.Mn1{}
	mn1.Print()

	// Output:
	// x.Mn1{}

}

// -------------------------------- Mn0 -------------------------------

func ExampleMn0() {

	mn0 := x.Mn0{}
	mn0.Print()

	// Output:
	// x.Mn0{}

}

// -------------------------------- Min -------------------------------

func ExampleMin() {

	min := x.Min{}
	min.Print()

	// Output:
	// x.Min{}

}

// -------------------------------- Max -------------------------------

func ExampleMax() {

	max := x.Max{}
	max.Print()

	// Output:
	// x.Max{}

}

// -------------------------------- Mmx -------------------------------

func ExampleMmx() {

	mmx := x.Mmx{}
	mmx.Print()

	// Output:
	// x.Mmx{}

}

// -------------------------------- Pos -------------------------------

func ExamplePos() {

	pos := x.Pos{}
	pos.Print()

	// Output:
	// x.Pos{}

}

// -------------------------------- Neg -------------------------------

func ExampleNeg() {

	neg := x.Neg{}
	neg.Print()

	// Output:
	// x.Neg{}

}

// -------------------------------- Any -------------------------------

func ExampleAny() {

	any := x.Any{}
	any.Print()

	// Output:
	// x.Any{}

}

// -------------------------------- Toi -------------------------------

func ExampleToi() {

	toi := x.Toi{}
	toi.Print()

	// Output:
	// x.Toi{}

}

// -------------------------------- Tox -------------------------------

func ExampleTox() {

	tox := x.Tox{}
	tox.Print()

	// Output:
	// x.Tox{}

}

// -------------------------------- Rng -------------------------------

func ExampleRng() {

	rng := x.Rng{}
	rng.Print()

	// Output:
	// x.Rng{}

}

// -------------------------------- End -------------------------------

func ExampleEnd() {

	end := x.End{}
	end.Print()

	// Output:
	// x.End{}

}
*/
