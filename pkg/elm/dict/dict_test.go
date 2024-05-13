package dict

import (
	"scion/pkg/elm/maybe"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDict(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		asserts.Equal(Dict[int, struct{}]{rbt: nil}, Empty[int, struct{}]())
	})

	t.Run("Singleton", func(t *testing.T) {
		SUT := Singleton[int, struct{}](1, struct{}{})
		asserts.Equal(&Dict[int, struct{}]{
			rbt: &node[int, struct{}]{
				key:    1,
				value:  struct{}{},
				color:  BLACK,
				parent: nil,
				left:   nil,
				right:  nil},
		},
			SUT,
		)
	})

	t.Run("Insert into empty dict", func(t *testing.T) {
		SUT := Empty[int, int]()
		SUT.Insert(1, 233)

		asserts.Equal(Dict[int, int]{rbt: &node[int, int]{
			key:    1,
			value:  233,
			color:  BLACK,
			parent: nil,
			left:   nil,
			right:  nil}},
			SUT,
		)
	})

	t.Run("Insert into existing entry", func(t *testing.T) {
		SUT := Singleton(10, 233)
		SUT.Insert(10, 100)

		asserts.Equal(&Dict[int, int]{rbt: &node[int, int]{
			key:    10,
			value:  100,
			color:  BLACK,
			parent: nil,
			left:   nil,
			right:  nil},
		}, SUT)
	})

	t.Run("Insert left in to tree", func(t *testing.T) {
		SUT := Singleton(10, 233)
		SUT.Insert(5, 23)

		asserts.Equal(&Dict[int, int]{rbt: &node[int, int]{
			key:    10,
			value:  233,
			color:  BLACK,
			parent: nil,
			left: &node[int, int]{
				key:    5,
				value:  23,
				color:  RED,
				parent: SUT.rbt,
				left:   nil,
				right:  nil,
			},
			right: nil,
		},
		}, SUT)

	})
	t.Run("Insert right in to tree", func(t *testing.T) {
		SUT := Singleton(10, 233)
		SUT.Insert(15, 23)

		asserts.Equal(&Dict[int, int]{rbt: &node[int, int]{
			key:    10,
			value:  233,
			color:  BLACK,
			parent: nil,
			right: &node[int, int]{
				key:    15,
				value:  23,
				color:  RED,
				parent: SUT.rbt,
				left:   nil,
				right:  nil,
			},
			left: nil,
		},
		}, SUT)

	})
	t.Run("Insert right/left in to tree", func(t *testing.T) {
		SUT := Singleton(10, 233)
		SUT.Insert(15, 23)
		SUT.Insert(5, 23)

		asserts.Equal(&Dict[int, int]{rbt: &node[int, int]{
			key:    10,
			value:  233,
			color:  BLACK,
			parent: nil,
			right: &node[int, int]{
				key:    15,
				value:  23,
				color:  RED,
				parent: SUT.rbt,
				left:   nil,
				right:  nil,
			},
			left: &node[int, int]{
				key:    5,
				value:  23,
				color:  RED,
				parent: SUT.rbt,
				left:   nil,
				right:  nil,
			},
		},
		}, SUT)

	})

	t.Run("Get existing entry", func(t *testing.T) {
		d := Singleton(10, 23)
		SUT := d.Get(10)

		asserts.Equal(maybe.Just[int]{Value: 23}, SUT)
	})

	t.Run("Get non-existing entry", func(t *testing.T) {
		d := Empty[int, int]()
		SUT := d.Get(10)

		asserts.Equal(maybe.Nothing{}, SUT)
	})
}
