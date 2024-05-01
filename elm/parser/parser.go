package parser

import "scion/elm/core"

type Parser[T any, P any] struct {
	parse
}

type Token[X any] struct {
	value     string
	expecting X
}

type parse func(s State) PStep

type PStep interface {
	pstep() _PStep
}

// Variants
type Good[T any] struct {
	_PStep
	State
	value T
}

type Bad[X any] struct {
	_PStep
	problem DeadEnd[X]
}

type _PStep struct{}

type State struct {
	src    string
	offset int
	row    int
	col    int
}

type DeadEnd[P any] struct {
	row     int
	col     int
	problem P
}

// Constructors
func NewGood[T any](s State, v T) PStep {
	return Good[T]{State: s, value: v}
}

func NewState(src string, offset int, row int, col int) State {
	return State{src: src, offset: offset, row: row, col: col}
}

// Methods
func (ps _PStep) pstep() _PStep {
	return ps
}

/* Parse exactly the given string, without any regard to what comes next.
 */
func (t *Token[X]) token() Parser[string, X] {
	return Parser[string, X]{
		parse: func(s State) PStep {
			nOffset, nRow, nCol := IsSubString(t.value, s.offset, s.row, s.col, s.src)
			nState := NewState(s.src, nOffset, nRow, nCol)

			if nOffset == -1 {
				return Bad[X]{problem: DeadEnd[X]{row: s.row, col: s.col, problem: t.expecting}}
			} else {
				return NewGood(nState, t.value)
			}
		},
	}
}

/*
A parser on it's own doesn't do anything, we need to run it.
*/
func Run[T any, P any](p Parser[T, P], source string) core.Result[T, []DeadEnd[P]] {
	init := p.parse(NewState(source, 0, 1, 1))

	return WithPStep(
		init,
		func(g *Good[T]) core.Result[T, []DeadEnd[P]] { return core.Ok[T, []DeadEnd[P]]{Value: g.value} },
		func(b *Bad[P]) core.Result[T, []DeadEnd[P]] {
			de := []DeadEnd[P]{b.problem}
			return core.Err[T, []DeadEnd[P]]{Value: de}
		},
	)
}

// Matcher
func WithPStep[R any, V any, X any](
	pstep PStep,
	good func(*Good[V]) R,
	bad func(*Bad[X]) R,
) R {
	switch d := pstep.(type) {
	case Good[V]:
		return good(&d)

	case Bad[X]:
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
