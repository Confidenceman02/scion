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

		SUT := Run(token.Token(), "hello world")

		asserts.Equal(elm.Ok[struct{}, []DeadEnd[elm.Never, string]]{Value: struct{}{}}, SUT)
	})

	t.Run("Token parser Err", func(t *testing.T) {
		token := Token[elm.Never, string]{Value: "foo", Expecting: "Expecting foo"}
		SUT := Run(token.Token(), "hello world")

		de := []DeadEnd[elm.Never, string]{DeadEnd[elm.Never, string]{Row: 1, Col: 1, Problem: "Expecting foo", ContextStack: nil}}

		asserts.Equal(elm.Err[struct{}, []DeadEnd[elm.Never, string]]{Value: de}, SUT)
	})
}

func TestInContext(t *testing.T) {
	asserts := assert.New(t)

	t.Run("DeadEnd with context", func(t *testing.T) {
		token := Token[int, string]{Value: "foo", Expecting: "Expecting foo"}
		parserWithContext := InContext(1, token.Token())
		SUT := Run(parserWithContext, "hello world")

		de := []DeadEnd[int, string]{DeadEnd[int, string]{
			Row:          1,
			Col:          1,
			Problem:      "Expecting foo",
			ContextStack: []Located[int]{Located[int]{Row: 1, Col: 1, Context: 1}},
		}}

		asserts.Equal(elm.Err[struct{}, []DeadEnd[int, string]]{Value: de}, SUT)
	})
}
