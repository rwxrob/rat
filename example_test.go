package pegn

import "fmt"

func ExampleFromString() {

	fmt.Println(pegn.FromString("some\tthing\nuh\rsmile😊"))

	// Output:
	// 'some' TAB 'thing' LR 'uh' CR 'smile' xe0fh

}
