package advanced

import (
	"github.com/stretchr/testify/assert"
	"scion/pkg/elm"
	"testing"
)

func TestTokenParser(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Token parser Ok", func(t *testing.T) {
		token := Token[string]{Value: "hello", Expecting: "Expecting hello"}

		SUT := Run[string](token.Token(), "hello world")

		asserts.Equal(elm.Ok[string, []DeadEnd[string]]{Value: "hello"}, SUT)
	})

	t.Run("Token parser Err", func(t *testing.T) {
		token := Token[string]{Value: "foo", Expecting: "Expecting foo"}
		SUT := Run[string](token.Token(), "hello world")

		de := []DeadEnd[string]{DeadEnd[string]{Row: 1, Col: 1, Problem: "Expecting foo"}}

		asserts.Equal(elm.Err[string, []DeadEnd[string]]{Value: de}, SUT)
	})
}
