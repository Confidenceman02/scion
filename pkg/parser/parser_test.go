package parser

import (
	"github.com/stretchr/testify/assert"
	"scion/pkg/elm"
	"testing"
)

func TestTokenParser(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Token parser Ok", func(t *testing.T) {
		token := Token("hello")

		SUT := Run(token, "hello world")

		asserts.Equal(elm.Ok[struct{}, []DeadEnd]{Value: struct{}{}}, SUT)
	})

	t.Run("Token parser Err", func(t *testing.T) {
		token := Token("foo")
		SUT := Run(token, "hello world")

		de := []DeadEnd{DeadEnd{Row: 1, Col: 1, Problem: Expecting{value: "foo"}}}

		asserts.Equal(elm.Err[struct{}, []DeadEnd]{Value: de}, SUT)
	})
}

func TestSymbolParser(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Symbol Ok", func(t *testing.T) {
		symbol := Symbol("(")

		SUT := Run(symbol, "(..)")

		asserts.Equal(elm.Ok[struct{}, []DeadEnd]{Value: struct{}{}}, SUT)
	})

	t.Run("Symbol Err", func(t *testing.T) {
		symbol := Symbol("(")

		SUT := Run(symbol, "..)")
		de := []DeadEnd{DeadEnd{Row: 1, Col: 1, Problem: ExpectingSymbol{value: "("}}}

		asserts.Equal(elm.Err[struct{}, []DeadEnd]{Value: de}, SUT)
	})
}
