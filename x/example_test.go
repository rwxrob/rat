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
	// x.Str{"string"}
	// x.Str{"bytes as string"}
	// x.Str{"runes as string"}
	// x.Str{"x"}
	// x.Str{"😀"}
	// x.Str{"true"}
	// x.Str{"false"}
	// x.Str{"-127"}
	// x.Str{"-32767"}
	// x.Str{"-9223372036854775808"}
	// x.Str{"-9785"}
	// x.Str{"127"}
	// x.Str{"255"}
	// x.Str{"3.141592653589793"}
	// x.Str{"32767"}
	// x.Str{"9223372036854775807"}
	// x.Str{"9786"}
	// x.Str{"\x00"}
	// x.Str{"☹"}
	// x.Str{"☺"}
	// x.Str{"👩"}
	// foo
	// x.Seq{x.Str{"1"}, x.Str{"false"}}
	// x.Seq{x.Str{"one"}, x.Str{"two"}}

}

func ExampleString_any_Slice() {
	fmt.Println(x.String([]any{`foo`}))
	fmt.Println(x.String([]any{}))

	// Output:
	// x.Str{"foo"}
	// "%!ERROR: invalid rat/x type or syntax"

}

func ExampleJoinStr() {

	fmt.Println(x.JoinStr("foo", "bar"))
	fmt.Println(x.JoinStr(x.Str{"foo"}, x.Str{"bar"}))
	fmt.Println(x.JoinStr(true, false))

	// Output:
	// foobar
	// foobar
	// truefalse

}

func ExampleCombineStr() {

	bunch := []any{false, true, 42, x.Val{`Foo`}, "foo", x.Str{"bar"}}
	comb := x.CombineStr(bunch...)
	for i, it := range comb {
		fmt.Printf("%v: %v\n", i, it)
	}

	// Output:
	// 0: x.Str{"falsetrue42"}
	// 1: x.Val{"Foo"}
	// 2: x.Str{"foobar"}

}

// ------------------------------- N -------------------------------

func ExampleN() {

	x.N{`FooName`, `foo`}.Print()
	x.N{`FooName`, `foo`, `toomuch`}.Print()
	x.N{false, `foo`}.Print()

	// Output:
	// x.N{"FooName", x.Str{"foo"}}
	// "%!USAGE: x.N{name, rule}"
	// "%!USAGE: x.N{name, rule}"

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

// -------------------------------- Save -------------------------------

func ExampleSave() {

	_ = x.N{`Foo`, x.Str{`foo`}}
	save := x.Save{`Foo`}
	save.Print()

	x.Save{false}.Print()
	x.Save{}.Print()

	// Output:
	// x.Save{"Foo"}
	// "%!USAGE: x.Save{name}"
	// "%!USAGE: x.Save{name}"

}

// -------------------------------- Val -------------------------------

func ExampleVal() {

	_ = x.N{`Foo`, x.Str{`foo`}}
	_ = x.Save{`Foo`}
	val := x.Val{`Foo`}
	val.Print()

	x.Val{false}.Print()
	x.Val{}.Print()

	// Output:
	// x.Val{"Foo"}
	// "%!USAGE: x.Val{name}"
	// "%!USAGE: x.Val{name}"

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
	x.Seq{`foo`, x.One{false, true}, `bar`}.Print()
	x.Seq{[]any{`foo`, `bar`}}.Print()
	x.Seq{[]any{`foo`}}.Print()
	x.Seq{`foo`}.Print()
	x.Seq{}.Print()
	x.Seq{[]any{}}.Print()

	// Output:
	// x.Str{"foofalsebar"}
	// x.Seq{x.Str{"foo"}, x.One{x.Str{"false"}, x.Str{"true"}}, x.Str{"bar"}}
	// x.Str{"foobar"}
	// x.Str{"foo"}
	// x.Str{"foo"}
	// "%!USAGE: x.Seq{...rule}"
	// "%!USAGE: x.Seq{...rule}"

}

// -------------------------------- One -------------------------------

func ExampleOne() {

	x.One{`foo`, false, `bar`}.Print()
	x.One{`foo`}.Print()
	x.One{}.Print()

	// Output:
	// x.One{x.Str{"foo"}, x.Str{"false"}, x.Str{"bar"}}
	// x.Str{"foo"}
	// "%!USAGE: x.One{...rule}"

}

// -------------------------------- Opt -------------------------------

func ExampleOpt() {

	x.Opt{`foo`}.Print()
	x.Opt{}.Print()
	x.Opt{`foo`, false}.Print()

	// Output:
	// x.Opt{x.Str{"foo"}}
	// "%!USAGE: x.Opt{rule}"
	// "%!USAGE: x.Opt{rule}"

}

// -------------------------------- Str -------------------------------

func ExampleStr() {

	smile := int32('\u263A')

	x.Str{
		`string`, []byte(`bytes as string`), []rune(`runes as string`), 'x', '😀',
		true, false, -127, -32767, -9223372036854775808, -9785, int8(127), uint8(255),
		3.141592653589793238, int16(32767), int64(9223372036854775807), int(smile),
		'\x00', int32(9785), smile, '👩',
	}.Print()

	// Output:
	// x.Str{"stringbytes as stringrunes as stringx😀truefalse-127-32767-9223372036854775808-97851272553.1415926535897933276792233720368547758079786\x00☹☺👩"}

}

func ExampleStr_any_Slice() {

	smile := int32('\u263A')
	types := []any{
		`string`, []byte(`bytes as string`), []rune(`runes as string`), 'x', '😀',
		true, false, -127, -32767, -9223372036854775808, -9785, int8(127), uint8(255),
		3.141592653589793238, int16(32767), int64(9223372036854775807), int(smile),
		'\x00', int32(9785), smile, '👩',
	}
	x.Str{types}.Print()
	x.Str{}.Print()
	x.Str{false}.Print()
	x.Str{[]any{}}.Print()

	// Output:
	// x.Str{"stringbytes as stringrunes as stringx😀truefalse-127-32767-9223372036854775808-97851272553.1415926535897933276792233720368547758079786\x00☹☺👩"}
	// "%!USAGE: x.Str{...any}"
	// x.Str{"false"}
	// "%!USAGE: x.Str{...any}"

}

// -------------------------------- Mn1 -------------------------------

func ExampleMn1() {

	x.Mn1{`foo`}.Print()
	x.Mn1{`foo`, `bar`}.Print()
	x.Mn1{}.Print()

	// Output:
	// x.Mn1{x.Str{"foo"}}
	// "%!USAGE: x.Mn1{rule}"
	// "%!USAGE: x.Mn1{rule}"
}

// -------------------------------- Mn0 -------------------------------

func ExampleMn0() {

	x.Mn0{`foo`}.Print()
	x.Mn0{}.Print()
	x.Mn0{`foo`, `bar`}.Print()

	// Output:
	// x.Mn0{x.Str{"foo"}}
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
	// x.Min{2, x.Str{"foo"}}
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
	// x.Max{2, x.Str{"foo"}}
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
	// x.Mmx{2, 4, x.Str{"foo"}}
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
	// x.Rep{2, x.Str{"foo"}}
	// "%!USAGE: x.Rep{n, rule}"
	// "%!USAGE: x.Rep{n, rule}"
	// "%!USAGE: x.Rep{n, rule}"

}

// -------------------------------- See -------------------------------

func ExampleSee() {

	x.See{`foo`}.Print()
	x.See{}.Print()
	x.See{`foo`, `bar`}.Print()

	// Output:
	// x.See{x.Str{"foo"}}
	// "%!USAGE: x.See{rule}"
	// "%!USAGE: x.See{rule}"

}

// -------------------------------- Not -------------------------------

func ExampleNot() {

	x.Not{`foo`}.Print()
	x.Not{}.Print()
	x.Not{`foo`, `bar`}.Print()

	// Output:
	// x.Not{x.Str{"foo"}}
	// "%!USAGE: x.Not{rule}"
	// "%!USAGE: x.Not{rule}"

}

// -------------------------------- To --------------------------------

func ExampleTo() {

	x.To{`foo`}.Print()
	x.To{}.Print()
	x.To{`foo`, `bar`}.Print()

	// Output:
	// x.To{x.Str{"foo"}}
	// "%!USAGE: x.To{rule}"
	// "%!USAGE: x.To{rule}"

}

// -------------------------------- Any -------------------------------

func ExampleAny() {

	x.Any{5}.Print()
	x.Any{}.Print()
	x.Any{`five`}.Print()

	// Output:
	// x.Any{5}
	// "%!USAGE: x.Any{n} or x.Any{m, n} or x.Any{m, 0}"
	// "%!USAGE: x.Any{n} or x.Any{m, n} or x.Any{m, 0}"

}

func ExampleAny_minmax() {

	x.Any{5, 10}.Print()
	x.Any{}.Print()
	x.Any{`five`, 10}.Print()
	x.Any{5, `ten`}.Print()

	// Output:
	// x.Any{5, 10}
	// "%!USAGE: x.Any{n} or x.Any{m, n} or x.Any{m, 0}"
	// "%!USAGE: x.Any{n} or x.Any{m, n} or x.Any{m, 0}"
	// "%!USAGE: x.Any{n} or x.Any{m, n} or x.Any{m, 0}"

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
