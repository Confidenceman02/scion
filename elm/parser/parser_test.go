package parser

import (
	"github.com/stretchr/testify/assert"
	"scion/elm/core"
	"testing"
)

func TestTokenParser(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Token parser Ok", func(t *testing.T) {
		token := Token[string]{value: "hello", expecting: "Expecting hello"}

		SUT := Run[string](token.token(), "hello world")

		asserts.Equal(core.Ok[string, []DeadEnd[string]]{Value: "hello"}, SUT)
	})

	t.Run("Token parser Err", func(t *testing.T) {
		token := Token[string]{value: "foo", expecting: "Expecting foo"}
		SUT := Run[string](token.token(), "hello world")

		de := []DeadEnd[string]{DeadEnd[string]{row: 1, col: 1, problem: "Expecting foo"}}

		asserts.Equal(core.Err[string, []DeadEnd[string]]{Value: de}, SUT)
	})
}
