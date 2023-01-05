/*
Package x (as in "expressions") contains most of the rat.Pack-able
types. These type definitions allow PEG grammars to be defined entirely
in compilable Go code and easily rendered to any other grammar
meta-language (such as PEGN). Typed []any slices are used by convention
to keep the syntax consistent. These can thought of as the tokens that
would be created after having tokenizing a higher-level grammar. All
types implement the fmt.Stringer interface producing valid Go code that
can be used when creating generators. Each type also implements a Print()
method that is shorthand for fmt.Println(self). These cover most regular
expression engines as well. It is common and expected for developers to
create collections of rules (into grammars) that are comprised of these
basic expression components.

    Rule - Foo <- rule or <:Foo rule >
    Ref  - reference another rule by name
    Is   - any PEGN or Unicode or POSIX class function
    Seq  - (rule1 rule2)
    One  - (rule1 / rule2)
    Opt  - rule?
    Lit  - ('foo' SP x20 u2563 CR LF)
    Mn1  - rule+
    Mn0  - rule*
    Min  - rule{n,}
    Max  - rule{0,n}
    Mmx  - rule{m,n}
    Pos  - &rule
    Neg  - !rule
    Any  - .{n}
    Toi  - ..rule
    Tox  - ...rule
    Rng  - [a-f] / [x43-x54] / [u3243-u4545]
    End  - !.

Note that rat.Pack automatically converts any unrecognized expression
argument into a literal (Lit) expression based on its fmt.Sprintf
representation. Also note that these assume that the data being checked
consists entirely of UTF-8 unicode code points ([]rune slice).

Usage and Syntax Errors

This package implements an interpreted language entirely as Go types of
[]any. As such usage and syntax errors need to be communicated
consistently. The fmt.Sprintf convention of prefixing strings with %!
has therefore been adopted. All types implement the fmt.Stringer
interface with a method that returns either the proper, compilable Go
code of that instance, or a string prefixed with %!. Only errors will
every begin with this prefix.

*/
package x

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
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
	// TODO
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

	default:
		return fmt.Sprintf(`"%v"`, v)
		//	return _SyntaxError

	}
}

// ------------------------------- Rule -------------------------------

// Rule encapsulated another rule with a name and optional integer ID.
// The first argument must be the rule to encapsulate (a valid rat/x
// expression type), the second the unique string name to use, and the
// third, the integer ID to associate with that string name and be used
// as the result type for that rule. Rules without an ID will be
// assigned one automatically.  A Rule with only one argument is
// interpreted as if the encapsulated rule was used directly (unwrapping
// it effectively from the x.Rule{}).
//
// PEGN
//
//    Foo <- rule / <:Foo rule >
//
type Rule []any

func (it Rule) String() string {
	switch len(it) {

	case 2: // rule, name
		name, isstring := it[1].(string)
		if !isstring {
			return _UsageRule
		}
		return `x.Rule{` + String(it[0]) + `, ` + String(name) + `}`

	case 1: // rule (uncommon, but acceptable)
		return String(it[0])

	case 0:
		return _UsageRule

	case 3: // rule, name, id
		name, isstring := it[1].(string)
		id, isint := it[2].(int)
		if !isstring || !isint {
			return ""
		}
		return `x.Rule{` + String(it[0]) + `, ` + String(name) + `, ` + strconv.Itoa(id) + `}`

	default:
		return _UsageRule

	}
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
// of runes. The unicode package is full of examples of these.
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
	return "IsFunc"
}

func (it Is) Print() { fmt.Println(it) }

// -------------------------------- Seq -------------------------------

// Seq represents a sequence of expressions. If more than one value
// assume combined values are an array of []any. If only a single value,
// assume it is a []any slice and needs to be processed.
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
			return String(it)
		}
		switch len(it) {
		case 0:
			return ""
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
		return ""
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

type One []any // (rule1 / rule2)

// -------------------------------- Opt -------------------------------

type Opt []any // rule?

// -------------------------------- Lit -------------------------------

type Lit []any // ('foo' SP x20 u2563 CR LF)

func (s Lit) String() string {
	if len(s) == 0 {
		return ""
	}
	it, isstring := s[0].(string)
	if !isstring {
		return _UsageLit
	}
	return fmt.Sprintf(`%q`, it)
}

func (s Lit) Print() { fmt.Println(s) }

// -------------------------------- Mn1 -------------------------------

type Mn1 []any // rule+

// -------------------------------- Mn0 -------------------------------

type Mn0 []any // rule*

// -------------------------------- Min -------------------------------

type Min []any // rule{n,}

// -------------------------------- Max -------------------------------

type Max []any // rule{0,n}

// -------------------------------- Mmx -------------------------------

type Mmx []any // rule{m,n}

// -------------------------------- Pos -------------------------------

type Pos []any // &rule

// -------------------------------- Neg -------------------------------

type Neg []any // !rule

// -------------------------------- Any -------------------------------

type Any []any // rune{n}

func (args Any) String() string {
	switch len(args) {
	case 1:
		n, isint := args[0].(int)
		if !isint {
			return _UsageAny
		}
		return `x.Any{` + strconv.Itoa(n) + `}`
	default:
		return _UsageAny
	}
}

func (a Any) Print() { fmt.Println(a) }

// -------------------------------- Toi -------------------------------

type Toi []any // ..rule

// -------------------------------- Tox -------------------------------

type Tox []any // ...rule

// -------------------------------- Rng -------------------------------

type Rng []any // [a-f] / [x43-x54] / [u3243-u4545]

// -------------------------------- End -------------------------------

type End []any // !.
