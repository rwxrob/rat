/*
Package rat...



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
func Pack(seq ...any) *Grammar { return new(Grammar).Pack(seq...) }

// RuleMaker implementations must return a new Rule created from any
// input (but usually from rat/x expressions and other Go types).
// Implementations may choose to cache the newly created rule and simply
// return a previously cached rule if the input arguments are identified
// as representing an identical previous rule. This fulfills the
// PEG packrat parsing requirement for functional memoization.
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
type Rule struct {
	Name  string    // name corresponding to ID (sometimes dynamically assigned)
	ID    int       // unique ID for Result one-one with Name
	Text  string    // rat/x compatible expression (ex: x.Seq{"foo", "bar"})
	Check CheckFunc // usually closure
}

func (r Rule) String() string { return r.Text }
func (r Rule) Print()         { fmt.Println(r) }

type CheckFunc func(r []rune, i int) Result
