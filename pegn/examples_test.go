package pegn_test

import (
	"fmt"

	"github.com/rwxrob/rat/pegn"
)

func ExampleFromString() {

	fmt.Printf("%q", pegn.FromString("some\tthing\nuh\rwhat\r\nsmile😈"))

	// Output:
	// "'some' TAB 'thing' LF 'uh' CR 'what' CR LF 'smile' x1f608"

}
