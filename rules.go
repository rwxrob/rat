package rat

import (
	"fmt"
)

// ----------------------------- CheckFunc ----------------------------

type CheckFunc func(r []rune, i int) Result

// ------------------------------- Rule -------------------------------

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
// and MUST be structs with at least a V value

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
		return fmt.Sprintf(`rat.Seq{%v}`, Quoted(s[0]))
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
