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

func TestNoRotation(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Insert with no rotations left", func(t *testing.T) {
		SUT := Singleton(10, 233)
		SUT.Insert(5, 23)
		SUT.Insert(15, 23)
		SUT.Insert(2, 23)

		asserts.NotNil(SUT)
		asserts.Equal(10, SUT.rbt.key)
		asserts.Equal(BLACK, SUT.rbt.color)
		asserts.Equal(BLACK, SUT.rbt.left.color)
		asserts.Equal(BLACK, SUT.rbt.right.color)
		asserts.Equal(RED, SUT.rbt.left.left.color)
	})
}

func TestInsertNoRotationsRight(t *testing.T) {
	asserts := assert.New(t)
	t.Run("Insert with no rotations right", func(t *testing.T) {
		d := Singleton(10, 233)
		d.Insert(5, 23)
		d.Insert(15, 23)
		d.Insert(16, 23)

		SUT := d.getNode(10)

		asserts.NotNil(SUT)
		asserts.Equal(10, SUT.key)
		asserts.Equal(BLACK, SUT.color)
		asserts.Equal(BLACK, SUT.left.color)
		asserts.Equal(BLACK, SUT.right.color)
		asserts.Equal(RED, SUT.right.right.color)
	})

}

func TestRightRotation(t *testing.T) {
	asserts := assert.New(t)

	t.Run("LL Right rotation", func(t *testing.T) {
		d := Singleton(50, 1)
		d.Insert(40, 2)
		d.Insert(30, 3)

		SUT := d.getNode(40)

		asserts.NotNil(SUT)
		asserts.Equal(40, SUT.key)
		asserts.Nil(SUT.parent)
		asserts.Equal(50, SUT.right.key)
		asserts.Nil(SUT.right.left)
		asserts.Equal(BLACK, SUT.color)
		asserts.Equal(RED, SUT.right.color)
	})
}
