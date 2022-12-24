package rat_test

import (
	"fmt"

	"github.com/rwxrob/rat"
)

func ExampleMatch() {

	r := []rune(`Something`)
	m := rat.Match{r, 3, 5, nil}
	fmt.Println(m)

	// Output:
	// {"B":3,"E":5}
}
