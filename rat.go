/*
Package rat implements a PEG packrat parser that takes advantage of Go's
unique ability to type switch to create an intermediary interpreted
language consisting entirely of Go types (mostly structs). These types
serve as tokens that any lexer would create when lexing a higher level
meta language such as PEG, PEGN, or regular expressions. Passing them to
rat.Pack compiles (memoizes into sync.Map) them into a grammar of
parsing functions not unlike how regular expressions are compiled before
use. The string representations of these structs (including rat.Grammar)
consist of completely valid, compilable Go code suitable for parser code
generation. More performant (VM) parsers can also be generated simply by
interpreting the rat.Pack language of typed parameters.
*/
package rat

import (
	"fmt"
)

const (
	SeqPackType = 1
	OnePackType = 2
)

var DefaultPackType = SeqPackType

func Pack(in ...any) *Grammar {
	g := NewGrammar()

	switch len(in) {

	case 0:
		return g

	case 1:

		switch v := in[0].(type) {

		case *Grammar:
			g.Import(v)

		case Grammar:
			g.Import(&v)

		default:
			g.main = g.Add(in)
		}

	default:
		switch DefaultPackType {
		case SeqPackType:
			g.main = g.Add(Seq(in))
		case OnePackType:
			g.main = g.Add(One(in))
		}

	}

	return g
}

func Quoted(in any) string {
	switch v := in.(type) {
	case fmt.Stringer:
		return v.String()
	case string:
		return fmt.Sprintf(`%q`, v)
	case []byte:
		return fmt.Sprintf(`%q`, string(v))
	case []rune:
		return fmt.Sprintf(`%q`, string(v))
	case rune:
		return fmt.Sprintf(`%q`, string(v))
	default:
		return fmt.Sprintf(`"%v"`, v)
	}
	return ""
}
