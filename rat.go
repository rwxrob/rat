/*
Package rat implements a PEG packrat parser that takes advantage of Go's
unique ability to switch on type to create an interpreted meta-language
consisting entirely of Go types. These types serve as tokens that any
lexer would create when lexing a higher level meta language such as PEG,
PEGN, or regular expressions. Passing these types to rat.Pack compiles
(memoizes) them into a grammar of parsing functions not unlike
compilation of regular expressions. The string representations of these
structs consists of entirely of valid, compilable Go code suitable for
parser code generation for parsers of any type, in any language. Simply
printing a Grammar instance to a file is suitable to generate such
a parser.

Prefer Pack over Make*

Although the individual Make* methods for each of the supported types
have been exported publicly allowing developers to call them directly
from within their own Rule implementations, most should use Pack
instead. Consider it the equivalent of compiling a regular expression.

*/
package rat

import (
	"fmt"
	"io"
)

// Pack interprets a sequence of any valid Go types into a Grammar
// suitable for checking and parsing any UTF-8 input. All arguments
// passed are assumed to rat/x expressions or literal forms (%q). These
// have special meaning corresponding to the fundamentals expressions of
// an PEG or regular expression. Consider these the tokens created by
// any lexer when processing any meta-language. Any grammar or
// structured data format that uses UTF-8 encoding can be fully
// expressed as compilable Go code using this method of interpretation.
//
// Alternative Call
//
// Pack is actually just shorthand equivalent to the following:
//
//     g := new(Grammar)
//     rule := g.MakeRule()
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

// Rule encapsulates a CheckFunc with a Name and Text representation.
// The Name is use as the unique key in the Grammar.Rules cache. Text
// can be anything, but it is strongly recommended that it contain rat/x
// compatible expression so that it can be used directly for code
// generation.
//
// Rules are created by implementations of RuleMaker the most important
// of which is Grammar. Almost every Rule encapsulates a different set
// of arguments enclosed in its CheckFunc. Once created, a Rule should
// be considered immutable. Field values must not change so that they
// correspond with the enclosed values within the CheckFunc closure and
// so that the Name can be used to uniquely identify the Rule.
//
type Rule struct {
	Name  string    // uniquely identifying name (sometimes dynamically assigned)
	Text  string    // prefer rat/x compatible expression (ex: x.Seq{"foo", "bar"})
	Check CheckFunc // closure created with a RuleMaker
}

// String implements the fmt.Stringer interface by returning the
// Rule.Text.
func (r Rule) String() string { return r.Text }

// Print is a shortcut for fmt.Println(rule) which calls String.
func (r Rule) Print() { fmt.Println(r) }

func (r Rule) Scan(in any) Result {
	var runes []rune
	switch v := in.(type) {
	case string:
		runes = []rune(v)
	case []byte:
		runes = []rune(string(v))
	case []rune:
		runes = v
	case io.Reader:
		buf, err := io.ReadAll(v)
		if err != nil {
			return Result{X: err}
		}
		runes = []rune(string(buf))
	}
	return r.Check(runes, 0)
}

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
// When a CheckFunc is composed of multiple sub-rules, each MUST be
// added to the Result.C (children) slice including any that generated
// errors. Some functions may opt to continue even if the result
// contained an error allowing recovery. Usually, however, a CheckFunc
// should stop on the first error and include it with the children.
// Usually, a CheckFunc should also set its error Result.X to that of
// the final Result that failed.
//
type CheckFunc func(r []rune, i int) Result

// IsFunc functions return true if the passed rune is contained in a set
// of runes. The unicode package contains several examples.
type IsFunc func(r rune) bool
