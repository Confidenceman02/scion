package dict

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDict(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		asserts.Equal(empty[int]{}, Empty[int]())
	})

	t.Run("Singleton", func(t *testing.T) {
		SUT := Singleton(1)
		asserts.Equal(node[int]{
			key:    1,
			color:  BLACK,
			parent: empty[int]{},
			left:   empty[int]{},
			right:  empty[int]{}},
			SUT)
	})

	t.Run("IsEmpty false", func(t *testing.T) {
		d := Singleton(1)
		SUT := IsEmpty[int](d)
		asserts.Equal(false, SUT)
	})

	t.Run("IsEmpty true", func(t *testing.T) {
		d := Empty[int]()
		SUT := IsEmpty[int](d)
		asserts.Equal(true, SUT)
	})
}
