package rat

import (
	"fmt"
)

// ----------------------------- CheckFunc ----------------------------

// CheckFunc examines the []rune buffer at a specific position for
// a specific grammar rule and should generally only be used from an
// encapsulating Rule so that it has a Text identifier associated with
// it. One or more Rules may, however, encapsulate the same CheckFunc
// function.
//
// CheckFunc MUST return a Result indicating success or failure by
// setting Err for failure. (Note that this is unlike many packrat
// designs that return nil to indicate rule failure.)
//
// CheckFunc MUST set Err if unable to match the entire rule and MUST
// advance to the End to the farthest possible position in the []rune
// slice before failure occurred. This allows for better recovery and
// specific user-facing error messages while promoting succinct rule
// development.
type CheckFunc func(r []rune, i int) Result

// ------------------------------- Rule -------------------------------

// Rule combines one rule function (Check) with some identifying text
// (by way of String and Print) providing rich possibilities for
// representing grammars textually.  Rule functions are the fundamental
// building blocks of any functional PEG packrat parser.
//
// The String must produce valid, compilable Go code. By default this is
// the same as how they rule would be passed to the Add or Pack
// functions. The prefix used is StringPrefix.
//
// Rules have no external dependencies allowing them to be safely
// combined from multiple packages. For best performance, however, Rules
// should be created and used from a central grammar with proper caching.
type Rule interface {
	String() string               // valid Go rat struct syntax
	Print()                       // must fmt.Println(it.String())
	Check(r []rune, i int) Result // must always set R,B,E
}

type arule struct {
	string
	CheckFunc
}

func (r arule) String() string               { return r.string }
func (r arule) Print()                       { fmt.Println(r.string) }
func (a arule) Check(r []rune, i int) Result { return a.CheckFunc(r, i) }

// all rule structs are designed to be used with composite definitions
// so that the pseudo-syntax uses braces for the most part

// ------------------------------ Lit ---------------------------------

type Lit struct{ V string }

func (s Lit) String() string { return fmt.Sprintf(`%q`, s.V) }

func (s Lit) Print() { fmt.Println(s) }

func (s Lit) Check(r []rune, i int) Result {
	var err error
	var n int
	e := i
	runes := []rune(s.V)
	for e < len(r) && n < len(runes) {
		if r[e] != runes[n] {
			err = ErrExpected{r[e]}
			break
		}
		e++
		n++
	}
	if n < len(runes) {
		err = ErrExpected{runes[n]}
	}
	return Result{R: r, B: i, E: e, X: err}
}

// ----------------------------- Seq -----------------------------------

type Seq []any

func (s Seq) String() string {
	switch len(s) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf(`%v.Seq{%v}`, StringPrefix, Quoted(s[0]))
	}
	str := `rat.Seq{` + fmt.Sprintf(`%v`, Quoted(s[0]))
	for _, it := range s[1:] {
		str += fmt.Sprintf(`, %v`, Quoted(it))
	}
	return str + `}`
}

func (s Seq) Print() { fmt.Println(s) }

func (s Seq) Check(r []rune, i int) Result {
	// TODO
	return Result{}
}

// ------------------------------- One --------------------------------

type One []any

func (s One) String() string {
	switch len(s) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf(`rat.One{%v}`, Quoted(s[0]))
	}
	str := `rat.One{` + fmt.Sprintf(`%v`, Quoted(s[0]))
	for _, it := range s[1:] {
		str += fmt.Sprintf(`, %v`, Quoted(it))
	}
	return str + `}`
}

func (s One) Print() { fmt.Println(s) }

func (s One) Check(r []rune, i int) Result {
	// TODO
	return Result{}
}
