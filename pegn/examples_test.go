package pegn_test

import (
	"fmt"

	"github.com/rwxrob/rat/pegn"
)

func ExampleFromString() {

	fmt.Printf("%q\n", pegn.FromString("some\tthing\nuh\rwhat\r\nsmile😈"))
	fmt.Printf("%q\n", pegn.FromString("some"))

	// Output:
	// "('some' TAB 'thing' LF 'uh' CR 'what' CR LF 'smile' x1f608)"
	// "'some'"

}
