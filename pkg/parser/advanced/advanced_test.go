package advanced

import (
	"github.com/stretchr/testify/assert"
	"scion/pkg/elm"
	"testing"
)

func TestTokenParser(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Token parser Ok", func(t *testing.T) {
		token := Token[elm.Never, string]{Value: "hello", Expecting: "Expecting hello"}

		SUT := Run[elm.Never, string](token.Token(), "hello world")

		asserts.Equal(elm.Ok[string, []DeadEnd[elm.Never, string]]{Value: "hello"}, SUT)
	})

	t.Run("Token parser Err", func(t *testing.T) {
		token := Token[elm.Never, string]{Value: "foo", Expecting: "Expecting foo"}
		SUT := Run[elm.Never, string](token.Token(), "hello world")

		de := []DeadEnd[elm.Never, string]{DeadEnd[elm.Never, string]{Row: 1, Col: 1, Problem: "Expecting foo", ContextStack: nil}}

		asserts.Equal(elm.Err[string, []DeadEnd[elm.Never, string]]{Value: de}, SUT)
	})
}
