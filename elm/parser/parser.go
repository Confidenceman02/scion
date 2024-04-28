package parser

import "scion/elm/core"

type Parser[T any] struct {
	parse
}

type parse func(s State) PStep

type PStep interface {
	pstep() _PStep
}

// Variants
type Good struct {
	_PStep
	State
}

type Bad struct {
	_PStep
}

type _PStep struct{}

type State struct {
	src    string
	offset int
	row    int
	col    int
}

type DeadEnd struct {
	row int
	col int
}

// Methods
func (ps _PStep) pstep() _PStep {
	return ps
}

func (p *Parser[T]) run(source string) core.Result[T, []DeadEnd] {
	init := p.parse(State{src: source, offset: 0, row: 1, col: 1})

	return WithPStep(
		init,
		func(g *Good) core.Result[T, []DeadEnd] { return core.Ok[T, []DeadEnd]{} },
		func(b *Bad) core.Result[T, []DeadEnd] { return core.Err[T, []DeadEnd]{} },
	)
}

// Matcher
func WithPStep[R any](
	pstep PStep,
	good func(*Good) R,
	bad func(*Bad) R,
) R {
	switch d := pstep.(type) {
	case Good:
		return good(&d)

	case Bad:
		return bad(&d)
	}
	panic("unreachable")
}

// Parsers
func Token(s string) Parser[string] {
	return Parser[string]{
		parse: func(s State) PStep {
			// TODO Check for substring
			// Mutate State
			return Good{}
		}}
}
