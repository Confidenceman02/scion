package advanced

import (
	"scion/pkg/elm"
	"scion/pkg/elm/char"
	"scion/pkg/elm/list"
	"scion/pkg/parser/internal"
)

type Parser[C any, Probelm, Value any] struct {
	Parse[C, Value]
}

type Token[C any, X any] struct {
	Value     string
	Expecting X
}

type Parse[C any, V any] func(s State[C]) PStep[V]

type PStep[V any] interface {
	pstep() _PStep[V]
}

type _PStep[V any] struct {
	value V
}

// PStep Variants

type Good[C any, V any] struct {
	_PStep[V]
	State[C]
}

type Bad[C any, V any, P any] struct {
	_PStep[V]
	problem DeadEnd[C, P]
}

type State[C any] struct {
	src     string
	offset  int
	row     int
	col     int
	context []Located[C]
}

type DeadEnd[C any, P any] struct {
	Row          int
	Col          int
	Problem      P
	ContextStack []Located[C]
}

type Located[C any] struct {
	Row     int
	Col     int
	Context C
}

// PStep Constructor
func NewGood[C any, V any](s State[C], v V) PStep[V] {
	return Good[C, V]{State: s, _PStep: _PStep[V]{value: v}}
}

func NewBad[C any, V any, P any](de DeadEnd[C, P]) PStep[V] {
	return Bad[C, V, P]{problem: de}
}

// State Constructor
func NewState[C any](src string, offset int, row int, col int) State[C] {
	return State[C]{src: src, offset: offset, row: row, col: col}
}

/*
We call `InContext` so that any dead end that occurs during parsing
will get this extra context information.

That way you can say things like, “I was expecting an equals sign in the `view` definition.” Context!
*/
func InContext[C any, V any, P any](c C, p Parser[C, P, V]) Parser[C, P, V] {
	return Parser[C, P, V]{
		Parse: func(s0 State[C]) PStep[V] {
			newContext := changeContext[C](list.Cons(Located[C]{Row: s0.row, Col: s0.col, Context: c}, s0.context), s0)
			parsed := p.Parse(newContext)

			return PStepWith(
				parsed,
				func(g *Good[C, V]) PStep[V] {
					return NewGood(changeContext(s0.context, g.State), g.value)
				},
				func(b *Bad[C, V, P]) PStep[V] {
					return NewBad[C, V, P](b.problem)
				},
			)
		},
	}
}

func changeContext[C any](c []Located[C], s2 State[C]) State[C] {
	return State[C]{
		src:     s2.src,
		offset:  s2.offset,
		row:     s2.row,
		col:     s2.col,
		context: c,
	}
}

func fromState[C any, P any](s State[C], p P) DeadEnd[C, P] {
	return DeadEnd[C, P]{Row: s.row, Col: s.col, Problem: p, ContextStack: s.context}

}

func Symbol[C any, X any](t Token[C, X]) Parser[C, X, struct{}] {
	return t.Token()
}

/*
Just like parser.Keyword but you provide a 'Token' to clearly indicate your
custom type of problems:

	func letParser() Parser[Context Problem struct{}] {
	  return Keyword(Token{"let", ExpectingLet})
	}

Note that this would fail to chomp `letter` because of the subsequent
characters. Use `Token` if you do not want that last letter check.
*/
func Keyword[C any, X any](t Token[C, X]) Parser[C, X, struct{}] {
	return Parser[C, X, struct{}]{
		Parse: func(s State[C]) PStep[struct{}] {
			newOffset, newRow, newCol := isSubString(t.Value, s.offset, s.row, s.col, s.src)
			if newOffset == -1 || 0 <= isSubChar(func(c int32) bool { return char.IsAlphaNum(c) || c == '_' }, newOffset, s.src) {
				return NewBad[C, struct{}, X](fromState[C, X](s, t.Expecting))

			} else {
				return NewGood(State[C]{s.src, newOffset, newRow, newCol, s.context}, struct{}{})
			}
		},
	}
}
func isSubString(ss string, offset int, row int, col int, src string) (int, int, int) {
	return internal.IsSubString(ss, offset, row, col, src)
}
func isSubChar(predicate func(c int32) bool, offset int, s string) int {
	return internal.IsSubChar(predicate, offset, s)
}

// Methods
func (ps _PStep[V]) pstep() _PStep[V] {
	return ps
}

/* Parse exactly the given string, without any regard to what comes next.
 */
func (t *Token[C, X]) Token() Parser[C, X, struct{}] {
	return Parser[C, X, struct{}]{
		Parse: func(s State[C]) PStep[struct{}] {
			nOffset, nRow, nCol := internal.IsSubString(t.Value, s.offset, s.row, s.col, s.src)
			nState := NewState[C](s.src, nOffset, nRow, nCol)

			if nOffset == -1 {
				return NewBad[C, struct{}, X](fromState[C, X](s, t.Expecting))
			} else {
				return NewGood(nState, struct{}{})
			}
		},
	}
}

/*
A parser on it's own doesn't do anything, we need to run it.
*/
func Run[C any, V any, P any](p Parser[C, P, V], source string) elm.Result[V, []DeadEnd[C, P]] {
	init := p.Parse(NewState[C](source, 0, 1, 1))

	return PStepWith(
		init,
		func(g *Good[C, V]) elm.Result[V, []DeadEnd[C, P]] { return elm.Ok[V, []DeadEnd[C, P]]{Value: g.value} },
		func(b *Bad[C, V, P]) elm.Result[V, []DeadEnd[C, P]] {
			de := []DeadEnd[C, P]{b.problem}
			return elm.Err[V, []DeadEnd[C, P]]{Value: de}
		},
	)
}

// Matcher
func PStepWith[C any, R any, V any, P any](
	pstep PStep[V],
	good func(*Good[C, V]) R,
	bad func(*Bad[C, V, P]) R,
) R {
	switch d := pstep.(type) {
	case Good[C, V]:
		return good(&d)

	case Bad[C, V, P]:
		return bad(&d)
	}
	panic("unreachable")
}
