/*
Package x (as in "expressions") contains most of the rat.Pack-able
types. These type definitions allow PEG grammars to be defined entirely
in compilable Go code and easily rendered to any other grammar
meta-language (such as PEGN). Typed []any slices are used by convention
to keep the syntax consistent. These can thought of as the tokens that
would be created after having tokenizing a higher-level grammar. All
types implement the fmt.Stringer interface producing valid Go code that
can be used when creating generators. These cover most regular
expression engines as well. It is common and expected for developers to
create collections of rules (into grammars) that are comprised of these
basic expression components.

    N   - Foo <- rule or <:Foo rule >
    Is  - any PEGN or Unicode or POSIX class
    Seq - (rule1 rule2)
    One - (rule1 / rule2)
    Opt - rule?
		Lit - ('foo' SP x20 u2563 CR LF)
    Mn1 - rule+
    Mn0 - rule*
    Min - rule{n,}
    Max - rule{0,n}
    Mmx - rule{m,n}
    Pos - &rule
    Neg - !rule
    Any - .{n}
    Toi - ..rule
    Tox - ...rule
		Rng - [a-f] / [x43-x54] / [u3243-u4545]
    End - !.

Note that rat.Pack automatically converts any unrecognized expression
argument into a literal (Lit) expression based on its fmt.Sprintf
representation. Also note that these assume that the data being checked
consists entirely of UTF-8 unicode code points ([]rune slice).

*/
package x

type N []any  // Foo <- rule / <:Foo rule >
type Is []any // any PEGN or Unicode or POSIX class

// Seq represents a sequence of expressions.
type Seq []any // (rule1 rule2)

type One []any // (rule1 / rule2)
type Opt []any // rule?
type Lit []any // ('foo' SP x20 u2563 CR LF)
type Mn1 []any // rule+
type Mn0 []any // rule*
type Min []any // rule{n,}
type Max []any // rule{0,n}
type Mmx []any // rule{m,n}
type Pos []any // &rule
type Neg []any // !rule
type Any []any // rune{n}
type Toi []any // ..rule
type Tox []any // ...rule
type Rng []any // [a-f] / [x43-x54] / [u3243-u4545]
type End []any // !.