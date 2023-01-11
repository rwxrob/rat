package rat

import (
	"fmt"

	"github.com/rwxrob/rat/x"
)

func ExampleMakeAny() {

	g := new(Grammar)
	rule := g.makeAny(x.Any{3})
	rule.Check([]rune(`..`), 0).Print()
	rule.Check([]rune(`...`), 0).Print()
	rule.Check([]rune(`....`), 0).Print()
	rule.Check([]rune(`....`), 2).Print()
	fmt.Println(g.rules[`x.Any{3}`].Name)
	fmt.Println(g.rules[`x.Any{3}`].Text)

	//Output:
	// {"B":0,"E":1,"X":"expected: x.Any{3}","R":".."}
	// {"B":0,"E":3,"R":"..."}
	// {"B":0,"E":3,"R":"...."}
	// {"B":2,"E":3,"X":"expected: x.Any{3}","R":"...."}
	// x.Any{3}
	// x.Any{3}

}

func ExampleMakeLit() {

	g := new(Grammar)
	foo := g.makeLit(`foo`)
	oo := g.makeLit(`oo`)
	foo.Check([]rune(`foo`), 0).Print()
	foo.Check([]rune(`fooo`), 0).Print()
	foo.Check([]rune(`fo`), 0).Print()
	oo.Check([]rune(`fooo`), 0).Print()
	oo.Check([]rune(`fooo`), 1).Print()
	oo.Check([]rune(`fooo`), 2).Print()
	fmt.Println(g.rules[`foo`].Name)
	fmt.Println(g.rules[`foo`].Text)
	fmt.Println(g.rules[`oo`].Name)
	fmt.Println(g.rules[`oo`].Text)

	//Output:
	// {"B":0,"E":3,"R":"foo"}
	// {"B":0,"E":3,"R":"fooo"}
	// {"B":0,"E":2,"X":"expected: o","R":"fo"}
	// {"B":0,"E":0,"X":"expected: o","R":"fooo"}
	// {"B":1,"E":3,"R":"fooo"}
	// {"B":2,"E":4,"R":"fooo"}
	// foo
	// "foo"
	// oo
	// "oo"

}
