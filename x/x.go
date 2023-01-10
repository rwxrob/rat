/*
Package x (as in "expressions") contains the rat/x grammar in the form
of Go []any types.  These type definitions allow any grammar that can be
expresses in PEG (or PEGN) to be implemented entirely in highly
performant, compilable Go code and easily rendered to any other
language. This makes rat/x highly useful for generating PEG parsing code
in any programming language, or as a replacement for regular expressions
(which lack lookarounds and other PEG-capable constructs).

Typed []any slices are used by convention to keep the syntax consistent.
These can be thought of as the tokens that would result from tokenizing
a higher-level grammar. All types implement the fmt.Stringer interface
producing valid Go code that can be used when creating generators. When
types are used incorrectly the string representation contains the
%!ERROR or %!USAGE prefix. Each type also implements a Print() method
that is shorthand for fmt.Println(self).

    Rule - Foo <- rule or <:Foo rule >
    Ref  - reference another rule by name
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
    Pos  - &rule
    Neg  - !rule
    Any  - .{n}
    Toi  - ..rule
    Tox  - ...rule
    Rng  - [a-f] / [x43-x54] / [u3243-u4545]
    End  - !.

See the documentation for each type for a details on syntax. Also see the included Examples.

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
			return _SyntaxError
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
		return fmt.Sprintf(`%q`, v)

	case []rune:
		return fmt.Sprintf(`%q`, string(v))

	case []byte:
		return fmt.Sprintf(`%q`, string(v))

	case rune:
		return fmt.Sprintf(`%q`, string(v))

	case bool:
		return fmt.Sprintf(`"%v"`, v)

	default:
		return fmt.Sprintf(`"%v"`, v)

	}
}

// ------------------------------- Rule -------------------------------

// Rule encapsulates another rule with a name and optional integer ID.
// Names can be any valid Go string. The first argument must be the
// unique name of the rule to encapsulate, the second argument must be
// an integer ID. If the ID is 0 (Unknown) an integer will be assigned
// automatically. Negative integer IDs are allowed. The third argument
// is the rule to encapsulate.
//
// PEGN
//
//    Foo <- rule / <:Foo rule >
//
type Rule []any

func (it Rule) String() string {
	if len(it) != 3 {
		return _UsageRule
	}
	if _, is := it[0].(string); !is {
		return _UsageRule
	}
	if _, is := it[1].(int); !is {
		return _UsageRule
	}
	return fmt.Sprintf(`x.Rule{%q, %v, %v}`, it[0], it[1], String(it[2]))
}

func (it Rule) Print() { fmt.Println(it) }

// -------------------------------- Ref -------------------------------

// Ref refers to another rule by name. This prevents having to assign
// rules to variables and use them in subsequent rules.
//
// PEGN
//
//     Foo     <- 'some' 'thing'
//     Another <- Foo 'else'
//
type Ref []any // EndOfLine <- CR? LF; Block <- rune+ EndOfLine

func (args Ref) String() string {
	switch len(args) {
	case 1:
		name, isstring := args[0].(string)
		if !isstring {
			return _UsageRef
		}
		return `x.Ref{"` + name + `"}`
	default:
		return _UsageRef
	}
}

func (it Ref) Print() { fmt.Println(it) }

// -------------------------------- Is --------------------------------

func funcName(it any) string {
	fp := reflect.ValueOf(it).Pointer()
	long := runtime.FuncForPC(fp).Name()
	parts := strings.Split(long, `.`)
	return parts[len(parts)-1]
}

// IsFunc functions return true if the passed rune is contained in a set
// of runes. The unicode package contains several examples.
type IsFunc func(r rune) bool

// Is takes a single IsFunc argument which must refer to a non-anonymous
// function. The function is encapsulated within the CheckFunc of the
// resulting rule.
//
// PEGN
//
//     ws   <- SP CR LF TAB
//     word <- (!ws rune)+
//
type Is []any

func (it Is) String() string {
	switch len(it) {
	case 1:
		var name string
		switch v := it[0].(type) {
		case func(r rune) bool:
			name = funcName(v)
		case IsFunc:
			name = funcName(v)
		default:
			return _UsageIs
		}
		if strings.HasPrefix(name, `func`) {
			return _UsageIs
		}
		return `x.Is{` + name + `}`
	default:
		return _UsageIs
	}
}

func (it Is) Print() { fmt.Println(it) }

// -------------------------------- Seq -------------------------------

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
			return _UsageSeq
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
		return _UsageSeq
	default:
		str := `x.Seq{` + String(rules[0])
		for _, rule := range rules[1:] {
			str += `, ` + String(rule)
		}
		return str + `}`
	}
}

func (rules Seq) Print() { fmt.Println(rules) }

// -------------------------------- One -------------------------------

// One represents one of a set of possible matching rules. If more than
// one value assume combined values are an array of []any. If only
// a single value and that value is an []any slice assume each value are
// the values of the set (somewhat like []any{}...). If just a single value
// that is anything but an []any slice, unwrap and handle as if just a single
// rule.
//
// PEGN
//
//     (rule1 / rule2)
//
type One []any

func (rules One) String() string {
	switch len(rules) {
	case 0:
		return _UsageOne
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

// -------------------------------- Opt -------------------------------

// Opt represents a single optional rule.
//
// PEGN
//
//     rule?
//
type Opt []any

func (it Opt) String() string {
	if len(it) != 1 {
		return _UsageOpt
	}
	return `x.Opt{` + String(it[0]) + `}`
}

func (rules Opt) Print() { fmt.Println(rules) }

// -------------------------------- Lit -------------------------------

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
		return _UsageLit

	case 1:
		it, isslice := rules[0].([]any)
		if !isslice {
			return String(rules[0])
		}
		if len(it) == 0 {
			return _UsageLit
		}
		var str string
		for _, rule := range it {
			it := String(rule)
			str += it[1 : len(it)-1]
		}
		return `"` + str + `"`

	default:
		var str string
		for _, rule := range rules {
			it := String(rule)
			str += it[1 : len(it)-1]
		}
		return `"` + str + `"`
	}

}

func (s Lit) Print() { fmt.Println(s) }

// -------------------------------- Mn1 -------------------------------

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
		return _UsageMn1
	}
	return `x.Mn1{` + String(it[0]) + `}`
}

func (it Mn1) Print() { fmt.Println(it) }

// -------------------------------- Mn0 -------------------------------

// Mn0 represents zero or more of a single rule.
//
// PEGN
//
//     rule*
//
type Mn0 []any // rule*

func (it Mn0) String() string {
	if len(it) != 1 {
		return _UsageMn0
	}
	return `x.Mn0{` + String(it[0]) + `}`
}

func (it Mn0) Print() { fmt.Println(it) }

// -------------------------------- Min -------------------------------

// Min represents a minimum number (n) of a single rule.
//
// PEGN
//
//     rule{n,}
//
type Min []any

func (it Min) String() string {
	if len(it) != 2 {
		return _UsageMin
	}
	if _, isint := it[0].(int); !isint {
		return _UsageMin
	}
	return fmt.Sprintf(`x.Min{%v, %v}`, it[0], String(it[1]))
}

func (it Min) Print() { fmt.Println(it) }

// -------------------------------- Max -------------------------------

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
		return _UsageMax
	}
	if _, isint := it[0].(int); !isint {
		return _UsageMax
	}
	return fmt.Sprintf(`x.Max{%v, %v}`, it[0], String(it[1]))
}

func (it Max) Print() { fmt.Println(it) }

// -------------------------------- Mmx -------------------------------

// Mmx represents a minimum and maximum number (n) of a single rule.
//
// PEGN
//
//     rule{m,n}
//
type Mmx []any

func (it Mmx) String() string {
	if len(it) != 3 {
		return _UsageMmx
	}
	if _, isint := it[0].(int); !isint {
		return _UsageMmx
	}
	if _, isint := it[1].(int); !isint {
		return _UsageMmx
	}
	return fmt.Sprintf(`x.Mmx{%v, %v, %v}`, it[0], it[1], String(it[2]))
}

func (it Mmx) Print() { fmt.Println(it) }

// -------------------------------- Rep -------------------------------

// Rep represents a minimum and maximum number (n) of a single rule.
//
// PEGN
//
//     rule{n}
//
type Rep []any

func (it Rep) String() string {
	if len(it) != 2 {
		return _UsageRep
	}
	if _, isint := it[0].(int); !isint {
		return _UsageRep
	}
	return fmt.Sprintf(`x.Rep{%v, %v}`, it[0], String(it[1]))
}

func (it Rep) Print() { fmt.Println(it) }

// -------------------------------- Pos -------------------------------

// Pos represents a positive lookahead assertion. The end of the result
// is always unchanged, but an error set if rule assertion fails.
//
// PEGN
//
//     &rule
//
type Pos []any

func (it Pos) String() string {
	if len(it) != 1 {
		return _UsagePos
	}
	return fmt.Sprintf(`x.Pos{%v}`, String(it[0]))
}

func (it Pos) Print() { fmt.Println(it) }

// -------------------------------- Neg -------------------------------

// Neg represents a negative lookahead assertion. The end of the result
// is always unchanged, but an error set if the rule assertion is true.
//
// PEGN
//
//     !rule
//
type Neg []any

func (it Neg) String() string {
	if len(it) != 1 {
		return _UsageNeg
	}
	return fmt.Sprintf(`x.Neg{%v}`, String(it[0]))
}

func (it Neg) Print() { fmt.Println(it) }

// -------------------------------- Any -------------------------------

// Any represents a specific number of any valid rune.
//
// PEGN
//
//    .{n}
//
type Any []any

func (it Any) String() string {
	if len(it) != 1 {
		return _UsageAny
	}
	if _, isint := it[0].(int); !isint {
		return _UsageAny
	}
	return fmt.Sprintf(`x.Any{%v}`, it[0])
}

func (it Any) Print() { fmt.Println(it) }

// -------------------------------- Toi -------------------------------

// Toi represents any rune up to and including the specified rule.
//
// PEGN
//
//    ..rule
//
type Toi []any // ..rule

func (it Toi) String() string {
	if len(it) != 1 {
		return _UsageToi
	}
	return fmt.Sprintf(`x.Toi{%v}`, String(it[0]))
}

func (it Toi) Print() { fmt.Println(it) }

// -------------------------------- Tox -------------------------------

// Tox represents any rune up to the specified rule, but excluding it.
//
// PEGN
//
//    ...rule
//
type Tox []any // ..rule

func (it Tox) String() string {
	if len(it) != 1 {
		return _UsageTox
	}
	return fmt.Sprintf(`x.Tox{%v}`, String(it[0]))
}

func (it Tox) Print() { fmt.Println(it) }

// -------------------------------- Rng -------------------------------

// Rng represents an inclusive range between any two valid runes.
//
// PEGN
//
//     [a-f] / [x43-x54] / [u3243-u4545]
//
type Rng []any

func (it Rng) String() string {
	if len(it) != 2 {
		return _UsageRng
	}
	if _, isrune := it[0].(rune); !isrune {
		return _UsageRng
	}
	if _, isrune := it[1].(rune); !isrune {
		return _UsageRng
	}
	return fmt.Sprintf(`x.Rng{%q, %q}`, it[0], it[1])
}

func (it Rng) Print() { fmt.Println(it) }

// -------------------------------- End -------------------------------

// End represents the end of data, that there are no more runes to
// examine.
//
// PEGN
//
//     !.
//
type End []any

func (it End) String() string {
	if len(it) != 0 {
		return _UsageEnd
	}
	return `x.End{}`
}

func (it End) Print() { fmt.Println(it) }
