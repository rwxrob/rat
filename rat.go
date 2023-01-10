/*
Package rat implements a PEG packrat parser that takes advantage of Go's
unique ability to switch on type to create an interpreted meta-language
consisting entirely of Go types. These types serve as tokens that any
lexer would create when lexing a higher level meta language such as PEG,
PEGN, or regular expressions. Passing these types to rat.Pack compiles
(memoizes) them into a grammar of parsing functions not unlike
compilation of regular expressions. The string representations of these
structs consists of entirely of valid, compilable Go code suitable for
parser code generation for parsers of any type, in any language.
*/
package rat

import "fmt"

// Pack interprets a sequence of any valid Go types into a Grammar
// suitable for checking and parsing any UTF-8 input. All arguments
// passed are assumed to be literals of their string forms except for
// all the types defined in the x ("expression") subpackage. These have
// special meaning corresponding to the fundamentals expressions of an
// PEG or regular expression. Consider these the tokens created by any
// lexer when processing any meta-language. Any grammar or structured
// data format that uses UTF-8 encoding can be fully expressed as
// compilable Go code using this method of interpretation.
//
// Memoization
//
// Memoization is a fundamental requirement for any PEG packrat parser.
// Pack automatically memoizes all expressions using closure functions
// and map entries matching the specific arguments to a specific
// expression. Results are always integer pointers to specific positions
// within the data passed so there is never wasteful redundancy. This
// maximizes performance and minimizes memory utilization.
//
func Pack(seq ...any) *Grammar { return new(Grammar).Pack(seq...) }

// RuleMaker implementations must return a new Rule created from any
// input (but usually from rat/x expressions and other Go types).
// Implementations may choose to cache the newly created rule and simply
// return a previously cached rule if the input arguments are identified
// as representing an identical previous rule. This fulfills the
// PEG packrat parsing requirement for functional memoization.
//
type RuleMaker interface {
	MakeRule(in any) *Rule
}

// Rule encapsulates a CheckFunc with a unique ID and Name without the
// scope of a Grammar. Text must be rat/x compatible expression so that
// it can be used directly for code generation. Rules are created
// implementations of RuleMaker since almost every Rule encapsulates
// a different set of arguments enclosed in its CheckFunc. Once created,
// a Rule is immutable. Field values must not change so that they
// correspond with the enclosed values within the CheckFunc closure and
// so that the Name can be used to uniquely identify the Rule from among
// others in a rat.Map.
//
type Rule struct {
	Name  string    // name corresponding to ID (sometimes dynamically assigned)
	ID    int       // unique ID for Result one-one with Name
	Text  string    // rat/x compatible expression (ex: x.Seq{"foo", "bar"})
	Check CheckFunc // usually closure
}

// String implements the fmt.Stringer interface by returning the
// Rule.Text.
func (r Rule) String() string { return r.Text }

// Print is a shortcut for fmt.Println(rule) which calls String.
func (r Rule) Print() { fmt.Println(r) }

// CheckFunc examines the []rune buffer at a specific position for
// a specific grammar rule and should generally only be used from an
// encapsulating Rule so that it has a Text identifier associated with
// it. One or more Rules may, however, encapsulate the same CheckFunc
// function.
//
// CheckFunc MUST return a Result indicating success or failure by
// setting Result.X (error) for failure. (Note that this is unlike many
// packrat designs that return nil to indicate rule failure.)
//
// CheckFunc MUST set Result.X (error) if unable to match the entire
// rule and MUST advance to the E (end) to the farthest possible position in
// the []rune slice before failure occurred. This allows for better
// recovery and specific user-facing error messages while promoting
// succinct rule development.
//
type CheckFunc func(r []rune, i int) Result

// IsFunc functions return true if the passed rune is contained in a set
// of runes. The unicode package contains several examples.
type IsFunc func(r rune) bool
