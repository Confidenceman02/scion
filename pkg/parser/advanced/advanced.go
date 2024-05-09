package advanced

import (
	"scion/pkg/elm"
	"scion/pkg/parser/internal"
)

type Parser[C any, Value any, Problem any] struct {
	Parse[C]
}

type Token[C any, X any] struct {
	Value     string
	Expecting X
}

type Parse[C any] func(s State[C]) PStep

type PStep interface {
	pstep() _PStep
}

// Variants
type Good[C any, T any] struct {
	_PStep
	State[C]
	value T
}

type Bad[C any, X any] struct {
	_PStep
	problem DeadEnd[C, X]
}

type _PStep struct{}

// TODO add context
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
	ContextStack []struct {
		Row     int
		Col     int
		Context C
	}
}

type Located[C any] struct {
	Row     int
	Col     int
	Context C
}

// Constructors
func NewGood[C any, T any](s State[C], v T) PStep {
	return Good[C, T]{State: s, value: v}
}

func NewState[C any](src string, offset int, row int, col int) State[C] {
	return State[C]{src: src, offset: offset, row: row, col: col}
}

// Methods
func (ps _PStep) pstep() _PStep {
	return ps
}

/* Parse exactly the given string, without any regard to what comes next.
 */
func (t *Token[C, X]) Token() Parser[C, string, X] {
	return Parser[C, string, X]{
		Parse: func(s State[C]) PStep {
			nOffset, nRow, nCol := internal.IsSubString(t.Value, s.offset, s.row, s.col, s.src)
			nState := NewState[C](s.src, nOffset, nRow, nCol)

			if nOffset == -1 {
				return Bad[C, X]{problem: DeadEnd[C, X]{Row: s.row, Col: s.col, Problem: t.Expecting}}
			} else {
				return NewGood(nState, t.Value)
			}
		},
	}
}

/*
A parser on it's own doesn't do anything, we need to run it.
*/
func Run[C any, T any, P any](p Parser[C, T, P], source string) elm.Result[T, []DeadEnd[C, P]] {
	init := p.Parse(NewState[C](source, 0, 1, 1))

	return PStepWith(
		init,
		func(g *Good[C, T]) elm.Result[T, []DeadEnd[C, P]] { return elm.Ok[T, []DeadEnd[C, P]]{Value: g.value} },
		func(b *Bad[C, P]) elm.Result[T, []DeadEnd[C, P]] {
			de := []DeadEnd[C, P]{b.problem}
			return elm.Err[T, []DeadEnd[C, P]]{Value: de}
		},
	)
}

// Matcher
func PStepWith[C any, R any, V any, X any](
	pstep PStep,
	good func(*Good[C, V]) R,
	bad func(*Bad[C, X]) R,
) R {
	switch d := pstep.(type) {
	case Good[C, V]:
		return good(&d)

	case Bad[C, X]:
		return bad(&d)
	}
	panic("unreachable")
}

// func MapChompedString[T any, R any](f func(v string, t T) R, p Parser[T]) Parser[R] {
// 	return Parser[R]{parse: func(s State) PStep {
// 		parsed := p.parse(s)
// 		pstep := WithPStep(
// 			parsed,
// 			func(g *Good) PStep { return Good{State: s} },
// 			func(b *Bad) PStep { return Bad{} },
// 		)
//
// 		return Good{}
// 	}}
// }

// mapChompedString : (String -> a -> b) -> Parser c x a -> Parser c x b
// mapChompedString func (Parser parse) =
//   Parser <| \s0 ->
//     case parse s0 of
//       Bad p x ->
//         Bad p x
//
//       Good p a s1 ->
//         Good p (func (String.slice s0.offset s1.offset s0.src) a) s1
