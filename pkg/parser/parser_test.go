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

		SUT := Run[string](token, "hello world")

		asserts.Equal(elm.Ok[string, []DeadEnd]{Value: "hello"}, SUT)
	})

	t.Run("Token parser Err", func(t *testing.T) {
		token := Token("foo")
		SUT := Run[string](token, "hello world")

		de := []DeadEnd{DeadEnd{Row: 1, Col: 1, Problem: Expecting{value: "foo"}}}

		asserts.Equal(elm.Err[string, []DeadEnd]{Value: de}, SUT)
	})
}
