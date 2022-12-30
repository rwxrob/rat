package rat_test

import (
	"fmt"

	"github.com/rwxrob/rat"
)

func ExampleLit() {

	foo := rat.Lit{"f\x23oo\t👩bar"}
	fmt.Println(foo)
	foo.Print()

	// Output:
	// "f#oo\t👩bar"
	// "f#oo\t👩bar"
}

func ExampleSeq() {

	foobar := rat.Seq{`foo`, 20 + 45, `bar`}
	fmt.Println(foobar)
	foobar.Print()

	// Output:
	// rat.Seq{"foo", "65", "bar"}
	// rat.Seq{"foo", "65", "bar"}
}

func ExampleGrammar_empty() {
	g := rat.NewGrammar()
	g.Print()
	fmt.Println(g.IsZero())
	// Output:
	// <empty>
	// true
}

func ExamplePack_primitives() {

	smile := int32('\u263A')
	g := rat.Pack(
		-127, -32767, -9223372036854775808, -9785, int8(127), uint8(255), 3.141592653589793238,
		int16(32767), int64(9223372036854775807), int(smile), '\x00', []byte(`bytes`), false,
		[]rune{'f', 'o', 'o'}, `some`, true, int32(9785), smile, '👩')
	g.Print()

	// Output:
	// rat.Seq{"-127", "-32767", "-9223372036854775808", "-9785", "127", "255", "3.141592653589793", "32767", "9223372036854775807", "9786", "\x00", "bytes", "false", "foo", "some", "true", "☹", "☺", "👩"}

}

func ExampleDefaultPack() {

	rat.Pack(`foo`, 42).Print()
	rat.DefaultPackType = rat.OnePackType
	rat.Pack(`foo`, 42).Print()

	// Output:
	// rat.Seq{"foo", "42"}
	// rat.One{"foo", "42"}
}

/*


func ExampleGrammar_Import() {

	g1 := rat.Pack(`foo`)
	g2 := rat.Pack(`bar`)
	g1.Import(g2)
	g1.Print()

	// Output:
	// rat.One{"bar", "foo"}

}

func ExampleGrammar_Pack() {

	g := rat.Pack(`foo`)
	g.Pack(`bar`)
	g.Print()

	// Output:
	// rat.One{"bar", "foo"}

}

func ExampleGrammar_Check_default_One() {

	g := rat.Pack(`foo`, `bar`)
	g.Check(`this is foobar`).Print()

	// Output:
	// out
}
*/
