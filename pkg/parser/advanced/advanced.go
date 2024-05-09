package advanced

import (
	"scion/pkg/elm"
	"scion/pkg/parser/internal"
)

type Parser[C any, Value any, Problem any] struct {
	Parse
}

type Token[C any, X any] struct {
	Value     string
	Expecting X
}

type Parse func(s State) PStep

type PStep interface {
	pstep() _PStep
}

// Variants
type Good[T any] struct {
	_PStep
	State
	value T
}

type Bad[C any, X any] struct {
	_PStep
	problem DeadEnd[C, X]
}

type _PStep struct{}

type State struct {
	src    string
	offset int
	row    int
	col    int
}

// TODO add a context
type DeadEnd[C any, P any] struct {
	Row          int
	Col          int
	Problem      P
	ContextStack []ContextStack[C]
}

type ContextStack[C any] struct {
	Row     int
	Col     int
	Context C
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
func (t *Token[C, X]) Token() Parser[C, string, X] {
	return Parser[C, string, X]{
		Parse: func(s State) PStep {
			nOffset, nRow, nCol := internal.IsSubString(t.Value, s.offset, s.row, s.col, s.src)
			nState := NewState(s.src, nOffset, nRow, nCol)

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
	init := p.Parse(NewState(source, 0, 1, 1))

	return PStepWith(
		init,
		func(g *Good[T]) elm.Result[T, []DeadEnd[C, P]] { return elm.Ok[T, []DeadEnd[C, P]]{Value: g.value} },
		func(b *Bad[C, P]) elm.Result[T, []DeadEnd[C, P]] {
			de := []DeadEnd[C, P]{b.problem}
			return elm.Err[T, []DeadEnd[C, P]]{Value: de}
		},
	)
}

// Matcher
func PStepWith[C any, R any, V any, X any](
	pstep PStep,
	good func(*Good[V]) R,
	bad func(*Bad[C, X]) R,
) R {
	switch d := pstep.(type) {
	case Good[V]:
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
