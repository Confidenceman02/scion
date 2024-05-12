package dict

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDict(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		asserts.Equal(empty[int]{}, Empty[int, struct{}]())
	})

	t.Run("Singleton", func(t *testing.T) {
		SUT := Singleton[int, struct{}](1, struct{}{})
		asserts.Equal(node[int, struct{}]{
			key:    1,
			value:  struct{}{},
			color:  BLACK,
			parent: empty[int]{},
			left:   empty[int]{},
			right:  empty[int]{}},
			SUT,
		)
	})

	t.Run("IsEmpty false", func(t *testing.T) {
		d := Singleton[int, struct{}](1, struct{}{})
		SUT := IsEmpty[int, struct{}](d)
		asserts.Equal(false, SUT)
	})

	t.Run("IsEmpty true", func(t *testing.T) {
		d := Empty[int, struct{}]()
		SUT := IsEmpty[int, struct{}](d)
		asserts.Equal(true, SUT)
	})

	t.Run("Insert into empty dict", func(t *testing.T) {
		d := Empty[int, int]()
		SUT := Insert(1, 233, d)

		asserts.Equal(node[int, int]{
			key:    1,
			value:  233,
			color:  BLACK,
			parent: empty[int]{},
			left:   empty[int]{},
			right:  empty[int]{}},
			SUT,
		)
	})
}
