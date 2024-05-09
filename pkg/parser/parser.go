package parser

import (
	"scion/pkg/elm"
	"scion/pkg/elm/list"
	"scion/pkg/parser/advanced"
)

type Parser[T any] struct {
	advanced.Parser[elm.Never, T, Problem]
}

type DeadEnd struct {
	Row     int
	Col     int
	Problem Problem
}

// Iterfaces
type Problem interface {
	problem() _Problem
}

type _Problem struct{}

// Problem variants
type Expecting struct {
	_Problem
	value string
}

func (p _Problem) problem() _Problem {
	return p
}

/* When you run into a `DeadEnd`, I record some information about why you
got stuck. This data is useful for producing helpful error messages. This is
how [`deadEndsToString`](#deadEndsToString) works!

**Note:** If you feel limited by this type (i.e. having to represent custom
problems as strings) I highly recommend switching to `Parser.Advanced`. It
lets you define your own `Problem` type. It can also track "context" which
can improve error messages a ton! This is how the Elm compiler produces
relatively nice parse errors, and I am excited to see those techniques applied
elsewhere!
*/
// type Problem
//   = Expecting String
//   | ExpectingInt
//   | ExpectingHex
//   | ExpectingOctal
//   | ExpectingBinary
//   | ExpectingFloat
//   | ExpectingNumber
//   | ExpectingVariable
//   | ExpectingSymbol String
//   | ExpectingKeyword String
//   | ExpectingEnd
//   | UnexpectedChar
//   | Problem String
//   | BadRepeat
/*
A parser on it's own doesn't do anything, we need to run it.
*/
func Run[T any](p Parser[T], source string) elm.Result[T, []DeadEnd] {
	result := advanced.Run(p.Parser, source)

	return elm.ResultWith(
		result,
		func(o *elm.Ok[T, []advanced.DeadEnd[elm.Never, Problem]]) elm.Result[T, []DeadEnd] {
			return elm.Ok[T, []DeadEnd]{Value: o.Value}
		},
		func(e *elm.Err[T, []advanced.DeadEnd[elm.Never, Problem]]) elm.Result[T, []DeadEnd] {
			deadends := list.Map[advanced.DeadEnd[elm.Never, Problem], DeadEnd](problemToDeadend, e.Value)
			return elm.Err[T, []DeadEnd]{Value: deadends}
		})
}

func problemToDeadend(ade advanced.DeadEnd[elm.Never, Problem]) DeadEnd {
	return DeadEnd{Row: ade.Row, Col: ade.Col, Problem: ade.Problem}
}

// Parsers
func Token(s string) Parser[string] {
	token := toToken(s)
	return Parser[string]{token.Token()}
}

func toToken(s string) advanced.Token[elm.Never, Problem] {
	return advanced.Token[elm.Never, Problem]{Value: s, Expecting: Expecting{value: s}}
}
