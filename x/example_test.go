package x_test

import (
	"fmt"
	"unicode"

	"github.com/rwxrob/rat/x"
)

type foo struct{}

func (foo) String() string { return `foo` }

func ExampleString() {

	// type foo struct{}
	// func (foo) String() string { return `foo` }

	smile := int32('\u263A')
	types := []any{
		`string`, []byte(`bytes as string`), []rune(`runes as string`), 'x', '😀',
		true, false, -127, -32767, -9223372036854775808, -9785, int8(127), uint8(255),
		3.141592653589793238, int16(32767), int64(9223372036854775807), int(smile),
		'\x00', int32(9785), smile, '👩', foo{}, []any{1, false}, []string{`one`, `two`},
	}
	for _, it := range types {
		fmt.Println(x.String(it))
	}

	// Output:
	// x.Lit{"string"}
	// x.Lit{"bytes as string"}
	// x.Lit{"runes as string"}
	// x.Lit{"x"}
	// x.Lit{"😀"}
	// x.Lit{"true"}
	// x.Lit{"false"}
	// x.Lit{"-127"}
	// x.Lit{"-32767"}
	// x.Lit{"-9223372036854775808"}
	// x.Lit{"-9785"}
	// x.Lit{"127"}
	// x.Lit{"255"}
	// x.Lit{"3.141592653589793"}
	// x.Lit{"32767"}
	// x.Lit{"9223372036854775807"}
	// x.Lit{"9786"}
	// x.Lit{"\x00"}
	// x.Lit{"☹"}
	// x.Lit{"☺"}
	// x.Lit{"👩"}
	// foo
	// x.Seq{x.Lit{"1"}, x.Lit{"false"}}
	// x.Seq{x.Lit{"one"}, x.Lit{"two"}}

}

func ExampleString_any_Slice() {
	fmt.Println(x.String([]any{`foo`}))
	fmt.Println(x.String([]any{}))

	// Output:
	// x.Lit{"foo"}
	// "%!ERROR: invalid rat/x type or syntax"

}

// ------------------------------- Name -------------------------------

func ExampleName() {

	x.Name{`FooName`, `foo`}.Print()
	x.Name{`FooName`, `foo`, `toomuch`}.Print()
	x.Name{false, `foo`}.Print()

	// Output:
	// x.Name{"FooName", x.Lit{"foo"}}
	// "%!USAGE: x.Name{name, rule}"
	// "%!USAGE: x.Name{name, rule}"

}

// -------------------------------- Ref -------------------------------

func ExampleRef() {

	ref := x.Ref{`foo`}
	ref.Print()

	ref[0] = true
	ref.Print()

	x.Ref{}.Print()

	// Output:
	// x.Ref{"foo"}
	// "%!USAGE: x.Ref{name}"
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
	x.Is{x.IsFunc(foo)}.Print()
	x.Is{another}.Print()
	x.Is{myIsPrint}.Print()
	x.Is{false}.Print()

	// anonymous functions not allowed for classes
	x.Is{func(r rune) bool { return true }}.Print()
	anon := func(r rune) bool { return true }
	x.Is{anon}.Print()
	x.Is{}.Print()
	x.Is{anon, false}.Print()

	// Output:
	// x.Is{aclass}
	// x.Is{aclass}
	// x.Is{IsPrint}
	// x.Is{myIsPrint}
	// "%!USAGE: namedFunc or x.IsFunc or x.Is{namedFunc}"
	// "%!USAGE: namedFunc or x.IsFunc or x.Is{namedFunc}"
	// "%!USAGE: namedFunc or x.IsFunc or x.Is{namedFunc}"
	// "%!USAGE: namedFunc or x.IsFunc or x.Is{namedFunc}"
	// "%!USAGE: namedFunc or x.IsFunc or x.Is{namedFunc}"

}

// -------------------------------- Seq -------------------------------

func ExampleSeq() {

	x.Seq{`foo`, false, `bar`}.Print()
	x.Seq{[]any{`foo`, `bar`}}.Print()
	x.Seq{[]any{`foo`}}.Print()
	x.Seq{`foo`}.Print()
	x.Seq{}.Print()
	x.Seq{[]any{}}.Print()

	// Output:
	// x.Seq{x.Lit{"foo"}, x.Lit{"false"}, x.Lit{"bar"}}
	// x.Seq{x.Lit{"foo"}, x.Lit{"bar"}}
	// x.Lit{"foo"}
	// x.Lit{"foo"}
	// "%!USAGE: x.Seq{...rule}"
	// "%!USAGE: x.Seq{...rule}"

}

// -------------------------------- One -------------------------------

func ExampleOne() {

	x.One{`foo`, false, `bar`}.Print()
	x.One{`foo`}.Print()
	x.One{}.Print()

	// Output:
	// x.One{x.Lit{"foo"}, x.Lit{"false"}, x.Lit{"bar"}}
	// x.Lit{"foo"}
	// "%!USAGE: x.One{...rule}"

}

// -------------------------------- Opt -------------------------------

func ExampleOpt() {

	x.Opt{`foo`}.Print()
	x.Opt{}.Print()
	x.Opt{`foo`, false}.Print()

	// Output:
	// x.Opt{x.Lit{"foo"}}
	// "%!USAGE: x.Opt{rule}"
	// "%!USAGE: x.Opt{rule}"

}

// -------------------------------- Lit -------------------------------

func ExampleLit() {

	smile := int32('\u263A')

	x.Lit{
		`string`, []byte(`bytes as string`), []rune(`runes as string`), 'x', '😀',
		true, false, -127, -32767, -9223372036854775808, -9785, int8(127), uint8(255),
		3.141592653589793238, int16(32767), int64(9223372036854775807), int(smile),
		'\x00', int32(9785), smile, '👩',
	}.Print()

	// Output:
	// x.Lit{"stringbytes as stringrunes as stringx😀truefalse-127-32767-9223372036854775808-97851272553.1415926535897933276792233720368547758079786\x00☹☺👩"}

}

func ExampleLit_any_Slice() {

	smile := int32('\u263A')
	types := []any{
		`string`, []byte(`bytes as string`), []rune(`runes as string`), 'x', '😀',
		true, false, -127, -32767, -9223372036854775808, -9785, int8(127), uint8(255),
		3.141592653589793238, int16(32767), int64(9223372036854775807), int(smile),
		'\x00', int32(9785), smile, '👩',
	}
	x.Lit{types}.Print()
	x.Lit{}.Print()
	x.Lit{false}.Print()
	x.Lit{[]any{}}.Print()

	// Output:
	// x.Lit{"stringbytes as stringrunes as stringx😀truefalse-127-32767-9223372036854775808-97851272553.1415926535897933276792233720368547758079786\x00☹☺👩"}
	// "%!USAGE: x.Lit{...any}"
	// x.Lit{"false"}
	// "%!USAGE: x.Lit{...any}"

}

// -------------------------------- Mn1 -------------------------------

func ExampleMn1() {

	x.Mn1{`foo`}.Print()
	x.Mn1{`foo`, `bar`}.Print()
	x.Mn1{}.Print()

	// Output:
	// x.Mn1{x.Lit{"foo"}}
	// "%!USAGE: x.Mn1{rule}"
	// "%!USAGE: x.Mn1{rule}"
}

// -------------------------------- Mn0 -------------------------------

func ExampleMn0() {

	x.Mn0{`foo`}.Print()
	x.Mn0{}.Print()
	x.Mn0{`foo`, `bar`}.Print()

	// Output:
	// x.Mn0{x.Lit{"foo"}}
	// "%!USAGE: x.Mn0{rule}"
	// "%!USAGE: x.Mn0{rule}"

}

// -------------------------------- Min -------------------------------

func ExampleMin() {

	x.Min{2, `foo`}.Print()
	x.Min{}.Print()
	x.Min{2, `foo`, `bar`}.Print()
	x.Min{`two`, `foo`}.Print()

	// Output:
	// x.Min{2, x.Lit{"foo"}}
	// "%!USAGE: x.Min{n, rule}"
	// "%!USAGE: x.Min{n, rule}"
	// "%!USAGE: x.Min{n, rule}"

}

// -------------------------------- Max -------------------------------

func ExampleMax() {

	x.Max{2, `foo`}.Print()
	x.Max{}.Print()
	x.Max{2, `foo`, `bar`}.Print()
	x.Max{`two`, `foo`}.Print()

	// Output:
	// x.Max{2, x.Lit{"foo"}}
	// "%!USAGE: x.Max{n, rule}"
	// "%!USAGE: x.Max{n, rule}"
	// "%!USAGE: x.Max{n, rule}"

}

// -------------------------------- Mmx -------------------------------

func ExampleMmx() {

	x.Mmx{2, 4, `foo`}.Print()
	x.Mmx{}.Print()
	x.Mmx{2, 4, `foo`, `bar`}.Print()
	x.Mmx{`two`, 4, `foo`}.Print()
	x.Mmx{2, `four`, `foo`}.Print()

	// Output:
	// x.Mmx{2, 4, x.Lit{"foo"}}
	// "%!USAGE: x.Mmx{n, m, rule}"
	// "%!USAGE: x.Mmx{n, m, rule}"
	// "%!USAGE: x.Mmx{n, m, rule}"
	// "%!USAGE: x.Mmx{n, m, rule}"

}

// -------------------------------- Rep -------------------------------

func ExampleRep() {

	x.Rep{2, `foo`}.Print()
	x.Rep{}.Print()
	x.Rep{2, `foo`, `bar`}.Print()
	x.Rep{`two`, `foo`}.Print()

	// Output:
	// x.Rep{2, x.Lit{"foo"}}
	// "%!USAGE: x.Rep{n, rule}"
	// "%!USAGE: x.Rep{n, rule}"
	// "%!USAGE: x.Rep{n, rule}"

}

// -------------------------------- Pos -------------------------------

func ExamplePos() {

	x.Pos{`foo`}.Print()
	x.Pos{}.Print()
	x.Pos{`foo`, `bar`}.Print()

	// Output:
	// x.Pos{x.Lit{"foo"}}
	// "%!USAGE: x.Pos{rule}"
	// "%!USAGE: x.Pos{rule}"

}

// -------------------------------- Neg -------------------------------

func ExampleNeg() {

	x.Neg{`foo`}.Print()
	x.Neg{}.Print()
	x.Neg{`foo`, `bar`}.Print()

	// Output:
	// x.Neg{x.Lit{"foo"}}
	// "%!USAGE: x.Neg{rule}"
	// "%!USAGE: x.Neg{rule}"

}

// -------------------------------- Any -------------------------------

func ExampleAny() {

	x.Any{5}.Print()
	x.Any{}.Print()
	x.Any{`five`}.Print()

	// Output:
	// x.Any{5}
	// "%!USAGE: x.Any{n} or x.Any{n, m}"
	// "%!USAGE: x.Any{n} or x.Any{n, m}"

}

func ExampleAny_minmax() {

	x.Any{5, 10}.Print()
	x.Any{}.Print()
	x.Any{`five`, 10}.Print()
	x.Any{5, `ten`}.Print()

	// Output:
	// x.Any{5, 10}
	// "%!USAGE: x.Any{n} or x.Any{n, m}"
	// "%!USAGE: x.Any{n} or x.Any{n, m}"
	// "%!USAGE: x.Any{n} or x.Any{n, m}"

}

// -------------------------------- Rng -------------------------------

func ExampleRng() {

	x.Rng{'a', 'Z'}.Print()
	x.Rng{}.Print()
	x.Rng{'a', 'Z', `foo`}.Print()
	x.Rng{`an a`, 'Z'}.Print()
	x.Rng{'a', `a Z`}.Print()

	// Output:
	// x.Rng{'a', 'Z'}
	// "%!USAGE: x.Rng{beg, end}"
	// "%!USAGE: x.Rng{beg, end}"
	// "%!USAGE: x.Rng{beg, end}"
	// "%!USAGE: x.Rng{beg, end}"

}

// -------------------------------- End -------------------------------

func ExampleEnd() {

	x.End{}.Print()
	x.End{`nope`}.Print()

	// Output:
	// x.End{}
	// "%!USAGE: x.End{}"

}
