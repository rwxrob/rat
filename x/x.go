/*
Package x (as in "expressions") contains the rat/x (pronounced "ratex")
language in the form of Go []any types. See rat.Pack examples to get
started using them quickly.

These type definitions allow any grammar that can be expresses in PEGN
to be implemented entirely in highly performant, compilable Go
code and easily rendered to any other language. This makes rat/x highly
useful for generating PEG parsing code in any programming language, or
as a replacement for regular expressions (which lack lookarounds and
other PEG-capable constructs).

Typed []any slices are used by convention to keep the syntax consistent.
These can be thought of as the tokens that would result from tokenizing
a higher-level grammar. All types implement the fmt.Stringer interface
producing valid Go code that can be used when creating generators. When
types are used incorrectly the string representation contains the
%!ERROR or %!USAGE prefix. Each type also implements a Print() method
that is shorthand for fmt.Println(self).

    Name - Foo <- rule
	  Save - =rule
	  Val	 - $rule
    Ref  - Bar <- Foo
    Is   - boolean class function
    Seq  - (rule1 rule2)
    One  - (rule1 / rule2)
    Opt  - rule?
    Lit  - ('foo' SP x20 u2563 CR LF)
    Mn1  - rule+
    Mn0  - rule*
    Min  - rule{n,}
    Max  - rule{0,n}
    Mmx  - rule{m,n}
    Rep  - rule{n}
    See  - &rule
    Not  - !rule
    To   - .. rule
    Any  - .{n} or .{m,n} or .{m,}
		Rng  - [a-f] / [x43-x54] / [u3243-u4545]
    End  - !.

See the documentation for each type for a details on syntax. Also see the included Examples.

Greedy matching

All checks are greedy (like PEG/PEGN). This means the longest possible progression is always returned as the result.

Errors included

Every rule in this package (and accompanying CheckFunc) always includes every sub-rule (child) within the results even if it fails (producing a Result.X). The error of the final sub-rule is set to the error for the parent as well.

First error stops

All rules stop evaluating when the first result with an error is detected (no inherent attempt to recover).

*/
package x

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// String returns a valid rat/x type for anything passed including all
// valid Go primitives as Lit (string) types. Generally, this is the
// fmt.Sprintf %v values wrapped in double quotes. Anything with an
// fmt.Stringer implementation is assumed to already be acceptable rat/x
// notation. If input is an []any or []string slice it is interpreted as
// an x.Seq type. []byte, []rune, and single runes are all interpreted
// as double-quoted strings. Invalid types are returned with the special
// %! prefix indicating an error of some kind (similar to the fmt
// package).
func String(it any) string {

	switch v := it.(type) {

	case fmt.Stringer:
		return v.String()

	case []any:

		switch len(v) {
		case 0:
			return SyntaxError
		case 1:
			return String(v[0])
		default:
			str := `x.Seq{` + String(v[0])
			for _, it := range v[1:] {
				str += `, ` + String(it)
			}
			return str + `}`
		}

	case []string:
		str := `x.Seq{` + String(v[0])
		for _, it := range v[1:] {
			str += `, ` + String(it)
		}
		return str + `}`

	case string:
		return fmt.Sprintf(`x.Lit{%q}`, v)

	case []rune:
		return fmt.Sprintf(`x.Lit{%q}`, string(v))

	case []byte:
		return fmt.Sprintf(`x.Lit{%q}`, string(v))

	case rune:
		return fmt.Sprintf(`x.Lit{%q}`, string(v))

	case bool:
		return fmt.Sprintf(`x.Lit{"%v"}`, v)

	case func(r rune) bool:
		return `x.Is{` + FuncName(v) + `}`

	case IsFunc:
		return `x.Is{` + FuncName(v) + `}`

	default:
		return fmt.Sprintf(`x.Lit{"%v"}`, v)

	}
}

// JoinLit takes the string form of each argument (by passing to String)
// and joins. Assumes types passed are literals. Does not work for other
// rat/x expressions.
func JoinLit(args ...any) string {
	var str string
	for _, it := range args {
		buf := String(it)
		str += buf[7 : len(buf)-2]
	}
	return str
}

// Name encapsulates another Result with a name. In PEGN these are
// called "significant" (<=) because they can be easily found in the
// parsed results tree. Names can be any valid Go string but keeping to
// non-whitespace UTF-8 runes is strongly recommended (and required for
// rendering to PEG/PEGN). The first argument must be the unique name of
// the rule to encapsulate. The second argument is the rule to associate
// with the name. The name appears in the results output JSON if set
// (see rat.Result).
//
// Unlike other expressions, Name does not have a child result. This is
// effectively the same as if the encapsulated rule was called and
// simply had the Result.Name changed. The encapsulated rule also has an
// additional entry added to the Rules cache for the name pointing to
// the same value as the unnamed version of the rule.
//
// Note that both the encapsulated rule and the new rule are both cached
// using different Name and Text values but the same Check function. The
// encapsulated rule uses whatever the default name for that type of
// rule is, the new named rule uses the specific name. Depending on the
// default name of the encapsulated rule (for example a Seq can get
// quite long) giving them a name could decrease lookup time for cache
// checks. The grammar rule cache is checked first for the encapsulated
// rule and that is used if it already exists, otherwise a new rule is
// created and cached for the encapsulated rule as well.
//
// Also note that both the encapsulated rule and the named rule use the
// exact same closure function.
//
// PEGN
//
//    Foo <= rule
//    Bar <= Foo{2}
//
type Name []any

func (it Name) String() string {
	if len(it) != 2 {
		return UsageName
	}
	if _, is := it[0].(string); !is {
		return UsageName
	}
	return fmt.Sprintf(`x.Name{%q, %v}`, it[0], String(it[1]))
}

func (it Name) Print() { fmt.Println(it) }

// Save saves the results of a successful rule as a rule representing
// the literal output of that result allowing it to be used later with
// Val. There can only be one saved result for any saved rule at a time.
// Like Ref, the first argument must be a string matching the name of
// a rule.
//
// PEGN
//
//     FenceTok  <- ( '~' / BQ){3,8}
//     Fenced    <- =FenceTok .. $FenceTok
//
type Save []any

func (args Save) String() string {
	switch len(args) {
	case 1:
		_, is := args[0].(string)
		if !is {
			return UsageSave
		}
		return fmt.Sprintf(`x.Save{%q}`, args[0])
	default:
		return UsageSave
	}
}

func (it Save) Print() { fmt.Println(it) }

// Val uses a literal rule created with Save.
type Val []any

func (args Val) String() string {
	switch len(args) {
	case 1:
		_, is := args[0].(string)
		if !is {
			return UsageVal
		}
		return fmt.Sprintf(`x.Val{%q}`, args[0])
	default:
		return UsageVal
	}
}

func (it Val) Print() { fmt.Println(it) }

// Ref refers to another rule by name and is evaluated at runtime
// allowing reference to entirely different rules to be used before they
// are imported. This prevents having to assign rules to variables and
// use them in subsequent rules. This also allows looking up dynamically
// created rules such as those from Save ($Foo). The same cached lookup
// is just done at a different point during runtime.
//
// PEGN
//
//     Foo     <- 'some' 'thing'
//     Another <- Foo 'else'
//     WithVar <- =Foo 'else' $Foo
//
type Ref []any

func (args Ref) String() string {
	switch len(args) {
	case 1:
		_, isstring := args[0].(string)
		if !isstring {
			return UsageRef
		}
		return fmt.Sprintf(`x.Ref{%q}`, args[0])
	default:
		return UsageRef
	}
}

func (it Ref) Print() { fmt.Println(it) }

// IsFunc functions return true if the passed rune is contained in a set
// of runes. The unicode package contains several examples.
type IsFunc func(r rune) bool

// FuncName returns the best guess at the function name without the
// package. Note that this is generally only useful when passing named
// funcitons. This function is called by Is when creating names for
// Grammar caching.
func FuncName(it any) string {
	fp := reflect.ValueOf(it).Pointer()
	long := runtime.FuncForPC(fp).Name()
	parts := strings.Split(long, `.`)
	return parts[len(parts)-1]
}

// Is takes a single IsFunc argument which must refer to a non-anonymous
// function (see FuncName). The function is encapsulated within the
// CheckFunc of the resulting rule.
//
// Note that creating an explicit Is []any slice is not required. Any
// named function that matches the IsFunc type will be properly handled.
//
// PEGN
//
//     ws   <- SP CR LF TAB
//     word <- (!ws rune)+
//
type Is []any

func (it Is) String() string {
	if len(it) != 1 {
		return UsageIs
	}
	switch v := it[0].(type) {
	case func(r rune) bool:
		name := FuncName(v)
		if name[0:4] == "func" {
			return UsageIs
		}
		return `x.Is{` + name + `}`
	case IsFunc:
		name := FuncName(v)
		if name[0:4] == "func" {
			return UsageIs
		}
		return `x.Is{` + name + `}`
	default:
		return UsageIs
	}
}

func (it Is) Print() { fmt.Println(it) }

// Seq represents a sequence of expressions. One represents one of
// a set of possible matching rules. If more than one value assume
// combined values are an array of []any. If only a single value and
// that value is an []any slice assume each value are the values of the
// set (somewhat like []any{}...). If just a single value that is anything
// but an [] any slice, unwrap and handle as if just a single rule.
//
// PEGN
//
//     Foo <- (rule1 rule2)
//
type Seq []any // (rule1 rule2)

func (rules Seq) String() string {
	switch len(rules) {
	case 1:
		it, isslice := rules[0].([]any)
		if !isslice {
			return String(rules[0])
		}
		switch len(it) {
		case 0:
			return UsageSeq
		case 1:
			return String(it[0])
		default:
			str := `x.Seq{` + String(it[0])
			for _, rule := range it[1:] {
				str += `, ` + String(rule)
			}
			return str + `}`
		}
	case 0:
		return UsageSeq
	default:
		str := `x.Seq{` + String(rules[0])
		for _, rule := range rules[1:] {
			str += `, ` + String(rule)
		}
		return str + `}`
	}
}

func (rules Seq) Print() { fmt.Println(rules) }

// One represents one of a set of possible matching rules in order from
// left to right. If more than one value assume combined values are an
// array of []any. If only a single value and that value is an []any
// slice assume each value are the values of the set (somewhat like []
// any{}...). If just a single value that is anything but an []any slice,
// unwrap and handle as if just a single rule.
//
// PEGN
//
//     (rule1 / rule2)
//
type One []any

func (rules One) String() string {
	switch len(rules) {
	case 0:
		return UsageOne
	case 1:
		return String(rules[0])
	default:
		str := `x.One{` + String(rules[0])
		for _, rule := range rules[1:] {
			str += `, ` + String(rule)
		}
		return str + `}`
	}
}

func (rules One) Print() { fmt.Println(rules) }

// Opt represents a single optional rule. Note that the result never
// fails, only advances on success.
//
// PEGN
//
//     rule?
//
type Opt []any

func (it Opt) String() string {
	if len(it) != 1 {
		return UsageOpt
	}
	return `x.Opt{` + String(it[0]) + `}`
}

func (rules Opt) Print() { fmt.Println(rules) }

// Lit represents any literal and allows combining literals from any
// other type into a single rule (see String). This is useful when
// a number of independent strings (or other types that are represented
// as strings) are wanted as a single new rule. Otherwise, each
// independent string will be considered a separate rule with its own
// result. This is akin to a dynamic join that is evaluated at the time
// of the check assertion. If the first and only value is an []any slice
// assume it is to be expanded ([]any{}...).
//
// PEGN
//
//     ('foo' SP x20 u2563 CR LF)
//
type Lit []any

func (rules Lit) String() string {
	switch len(rules) {

	case 0:
		return UsageLit

	case 1:
		it, is := rules[0].([]any)
		if !is {
			return String(rules[0])
		}
		if len(it) == 0 {
			return UsageLit
		}
		return `x.Lit{"` + JoinLit(it...) + `"}`

	default:
		return `x.Lit{"` + JoinLit(rules...) + `"}`
	}

}

func (s Lit) Print() { fmt.Println(s) }

// Mn1 represents one or more of a single rule. If the first
// and only value is an []any slice assume it is to be expanded
// ([]any{}...).
//
// PEGN
//
//     rule+
//
type Mn1 []any

func (it Mn1) String() string {
	if len(it) != 1 {
		return UsageMn1
	}
	return `x.Mn1{` + String(it[0]) + `}`
}

func (it Mn1) Print() { fmt.Println(it) }

// Mn0 represents zero or more of a single rule.
//
// PEGN
//
//     rule*
//
type Mn0 []any // rule*

func (it Mn0) String() string {
	if len(it) != 1 {
		return UsageMn0
	}
	return `x.Mn0{` + String(it[0]) + `}`
}

func (it Mn0) Print() { fmt.Println(it) }

// Min represents a minimum number (n) of a single rule.
//
// PEGN
//
//     rule{n,}
//
type Min []any

func (it Min) String() string {
	if len(it) != 2 {
		return UsageMin
	}
	if _, isint := it[0].(int); !isint {
		return UsageMin
	}
	return fmt.Sprintf(`x.Min{%v, %v}`, it[0], String(it[1]))
}

func (it Min) Print() { fmt.Println(it) }

// Max represents a maximum number (n) of a single rule. Minimum is
// assumed to be zero.
//
// PEGN
//
//     rule{0,n}
//
type Max []any

func (it Max) String() string {
	if len(it) != 2 {
		return UsageMax
	}
	if _, isint := it[0].(int); !isint {
		return UsageMax
	}
	return fmt.Sprintf(`x.Max{%v, %v}`, it[0], String(it[1]))
}

func (it Max) Print() { fmt.Println(it) }

// Mmx represents a minimum (m) and maximum number (n) of a single rule.
// The minimum must be greater than zero. The maximum must be greater
// than the minimum.
//
// PEGN
//
//     rule{m,n}
//
type Mmx []any

func (it Mmx) String() string {
	if len(it) != 3 {
		return UsageMmx
	}
	if _, isint := it[0].(int); !isint {
		return UsageMmx
	}
	if _, isint := it[1].(int); !isint {
		return UsageMmx
	}
	return fmt.Sprintf(`x.Mmx{%v, %v, %v}`, it[0], it[1], String(it[2]))
}

func (it Mmx) Print() { fmt.Println(it) }

// Rep represents a minimum and maximum number (n) of a single rule.
//
// PEGN
//
//     rule{n}
//
type Rep []any

func (it Rep) String() string {
	if len(it) != 2 {
		return UsageRep
	}
	if _, isint := it[0].(int); !isint {
		return UsageRep
	}
	return fmt.Sprintf(`x.Rep{%v, %v}`, it[0], String(it[1]))
}

func (it Rep) Print() { fmt.Println(it) }

// See represents a positive lookahead assertion. The end of the result
// is always unchanged, but an error set if rule assertion fails.
//
// PEGN
//
//     &rule
//
type See []any

func (it See) String() string {
	if len(it) != 1 {
		return UsageSee
	}
	return fmt.Sprintf(`x.See{%v}`, String(it[0]))
}

func (it See) Print() { fmt.Println(it) }

// Not represents a negative lookahead assertion. The end of the result
// is always unchanged, but an error set if the rule assertion is true.
//
// PEGN
//
//     !rule
//
type Not []any

func (it Not) String() string {
	if len(it) != 1 {
		return UsageNot
	}
	return fmt.Sprintf(`x.Not{%v}`, String(it[0]))
}

func (it Not) Print() { fmt.Println(it) }

// To represents every rune until the rule matches successfully. This is
// essentially shorthand for the same done with looping negative
// expression lookaheads. Note that the rule itself is never included in
// the results (but usually is scanned itself immediately after).
//
// PEGN
//
//    .. rule
//
type To []any

func (it To) String() string {
	if len(it) != 1 {
		return UsageTo
	}
	return fmt.Sprintf(`x.To{%v}`, String(it[0]))
}

func (it To) Print() { fmt.Println(it) }

// Any represents a specific number of any valid rune. If more than one
// argument, then first is minimum and second maximum. If the maximum
// is 0 then maximum is unlimited and consumes all runes remaining.
//
// PEGN
//
//    .{n}
//    .{m,n}
//
type Any []any

func (it Any) String() string {
	switch len(it) {

	case 1: // exact number
		if _, isint := it[0].(int); !isint {
			return UsageAny
		}
		return fmt.Sprintf(`x.Any{%v}`, it[0])

	case 2: // min to max
		if _, isint := it[0].(int); !isint {
			return UsageAny
		}
		if _, isint := it[1].(int); !isint {
			return UsageAny
		}
		return fmt.Sprintf(`x.Any{%v, %v}`, it[0], it[1])

	default:
		return UsageAny
	}
}

func (it Any) Print() { fmt.Println(it) }

// Rng represents an inclusive range between any two valid runes.
//
// PEGN
//
//     [a-f] / [x43-x54] / [u3243-u4545]
//
type Rng []any

func (it Rng) String() string {
	if len(it) != 2 {
		return UsageRng
	}
	if _, isrune := it[0].(rune); !isrune {
		return UsageRng
	}
	if _, isrune := it[1].(rune); !isrune {
		return UsageRng
	}
	return fmt.Sprintf(`x.Rng{%q, %q}`, it[0], it[1])
}

func (it Rng) Print() { fmt.Println(it) }

// End represents the end of data, that there are no more runes to
// examine. End must be an empty []any slice for consistency and to
// allow a String representation method to be attached.
//
// PEGN
//
//     !.
//
type End []any

func (it End) String() string {
	if len(it) > 0 {
		return (UsageEnd)
	}
	if len(it) != 0 {
		return UsageEnd
	}
	return `x.End{}`
}

func (it End) Print() { fmt.Println(it) }
