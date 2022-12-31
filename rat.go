/*
Package rat...



*/
package rat

import "fmt"

// creates and caches rules from any valid sequence value, not safe for
// concurrency to keep performant (combine with semaphore when concurrency needed)
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

func (r Rule) String() string { return r.Name }
func (r Rule) Print()         { fmt.Println(r) }

type CheckFunc func(r []rune, i int) Result
