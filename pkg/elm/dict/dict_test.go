package dict

import (
	"github.com/stretchr/testify/assert"
	"scion/pkg/elm/maybe"
	"testing"
)

func TestInsert(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Empty", func(t *testing.T) {
		asserts.Equal(Dict[int, struct{}]{root: nil}, Empty[int, struct{}]())
	})

	t.Run("Singleton", func(t *testing.T) {
		SUT := Singleton[int, struct{}](1, struct{}{})
		asserts.Equal(&Dict[int, struct{}]{
			root: &node[int, struct{}]{
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

	t.Run("Insert nil root", func(t *testing.T) {
		SUT := Empty[int, int]()
		SUT.Insert(1, 233)

		asserts.Equal(Dict[int, int]{root: &node[int, int]{
			key:    1,
			value:  233,
			color:  BLACK,
			parent: nil,
			left:   nil,
			right:  nil}},
			SUT,
		)
	})

	t.Run("Insert root right side", func(t *testing.T) {
		SUT := Singleton[int, int](1, 1)
		SUT.Insert(2, 2)

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(RED, SUT.root.right.color)
	})

	t.Run("Insert into existing entry", func(t *testing.T) {
		SUT := Singleton(10, 233)
		SUT.Insert(10, 100)

		asserts.Equal(&Dict[int, int]{root: &node[int, int]{
			key:    10,
			value:  100,
			color:  BLACK,
			parent: nil,
			left:   nil,
			right:  nil},
		}, SUT)
	})

	t.Run("LL Single right rotation", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(40, 2)
		SUT.Insert(30, 3)

		asserts.Equal(40, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(50, SUT.root.right.key)
		asserts.Equal(RED, SUT.root.right.color)
		asserts.Equal(30, SUT.root.left.key)
		asserts.Equal(RED, SUT.root.left.color)
	})

	t.Run("RR Single right rotation", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(60, 2)
		SUT.Insert(70, 3)

		asserts.Equal(60, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(50, SUT.root.left.key)
		asserts.Equal(RED, SUT.root.left.color)
		asserts.Equal(70, SUT.root.right.key)
		asserts.Equal(RED, SUT.root.right.color)
	})

	t.Run("LR double red, red uncle", func(t *testing.T) {
		SUT := Singleton(50, 1)
		// Left
		SUT.Insert(40, 2)
		SUT.Insert(60, 3)
		SUT.Insert(45, 4)

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(BLACK, SUT.root.left.color)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(RED, SUT.root.left.right.color)
	})

	t.Run("LR double red, black uncle", func(t *testing.T) {
		SUT := Singleton(50, 1)
		// Left
		SUT.Insert(40, 2)
		SUT.Insert(45, 3)

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(45, SUT.root.key)
		asserts.Equal(RED, SUT.root.left.color)
		asserts.Equal(40, SUT.root.left.key)
		asserts.Equal(RED, SUT.root.right.color)
		asserts.Equal(50, SUT.root.right.key)
	})

	t.Run("RL double red, red uncle", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(60, 2)
		SUT.Insert(40, 3)
		SUT.Insert(55, 4)

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(BLACK, SUT.root.left.color)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(RED, SUT.root.right.left.color)
	})
	t.Run("RL double red, black uncle", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(60, 2)
		SUT.Insert(55, 4)

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(55, SUT.root.key)
		asserts.Equal(RED, SUT.root.left.color)
		asserts.Equal(50, SUT.root.left.key)
		asserts.Equal(RED, SUT.root.right.color)
		asserts.Equal(60, SUT.root.right.key)
	})

	t.Run("test the following inserts 7,5,10,20,15", func(t *testing.T) {
		SUT := Singleton(7, 1)
		SUT.Insert(5, 2)
		SUT.Insert(10, 3)
		SUT.Insert(20, 3)
		SUT.Insert(15, 3)

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(7, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(15, SUT.root.right.key)
		asserts.Equal(RED, SUT.root.right.right.color)
		asserts.Equal(20, SUT.root.right.right.key)
		asserts.Equal(RED, SUT.root.right.left.color)
		asserts.Equal(10, SUT.root.right.left.key)

	})

	t.Run("test the following inserts 10,15,5,0,2", func(t *testing.T) {
		SUT := Singleton(10, 1)
		SUT.Insert(15, 2)
		SUT.Insert(5, 3)
		SUT.Insert(0, 3)
		SUT.Insert(2, 3)

		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(10, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.left.color)
		asserts.Equal(2, SUT.root.left.key)
		asserts.Equal(RED, SUT.root.left.left.color)
		asserts.Equal(0, SUT.root.left.left.key)
		asserts.Equal(RED, SUT.root.left.right.color)
		asserts.Equal(5, SUT.root.left.right.key)

	})
}

func TestGet(t *testing.T) {
	asserts := assert.New(t)
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

func TestRemove(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Removes root node with no children", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Remove(50)

		asserts.Nil(SUT.root)
	})

	t.Run("Removes root node with 2 red children", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(60, 2)
		SUT.Insert(40, 3)
		SUT.Remove(50)

		asserts.Equal(40, SUT.root.key)
		asserts.Equal(3, SUT.root.value)
		asserts.Nil(SUT.root.left)
		asserts.Equal(60, SUT.root.right.key)
	})

	t.Run("Removes red right leaf node with no children", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(40, 2)
		SUT.Insert(60, 3)
		SUT.Remove(60)

		asserts.Nil(SUT.root.right)
	})

	t.Run("Removes a red left node with no children", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(40, 2)
		SUT.Insert(60, 3)
		SUT.Remove(40)

		asserts.Nil(SUT.root.left)
	})

	t.Run("Removes a black leaf node with 1 child | Left", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(40, 2)
		SUT.Insert(60, 3)
		SUT.Insert(45, 4)

		SUT.Remove(40)

		asserts.Equal(50, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(45, SUT.root.left.key)
		asserts.Equal(BLACK, SUT.root.left.color)
		asserts.Nil(SUT.root.left.right)
	})

	t.Run("Removes a black leaf node with 1 child | Right", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(40, 2)
		SUT.Insert(60, 3)
		SUT.Insert(55, 4)

		SUT.Remove(60)

		asserts.Equal(50, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(55, SUT.root.right.key)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Nil(SUT.root.right.left)
	})

	t.Run("Removes a black leaf node with no children | p = RED | s = BLACK with no children", func(t *testing.T) {
		SUT := Singleton(10, 1)
		SUT.Insert(5, 2)
		SUT.Insert(20, 3)
		SUT.Insert(15, 4)
		SUT.Insert(30, 5)

		// Mutate tree to for testing
		SUT.root.right.color = RED
		SUT.root.right.right.color = BLACK
		SUT.root.right.left.color = BLACK
		SUT.root.left.color = BLACK

		SUT.Remove(15)

		asserts.Nil(SUT.root.right.left)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(RED, SUT.root.right.right.color)
	})

	t.Run("Removes a black leaf node with no children | p = BLACK | s = BLACK with no children", func(t *testing.T) {
		SUT := Singleton(10, 1)
		SUT.Insert(5, 2)
		SUT.Insert(20, 3)
		SUT.Insert(1, 2)
		SUT.Insert(7, 2)
		SUT.Insert(15, 4)
		SUT.Insert(30, 5)

		// Manually balance for testing scenario
		// RIGHT
		SUT.root.right.color = BLACK
		SUT.root.right.right.color = BLACK
		SUT.root.right.left.color = BLACK
		// LEFT
		SUT.root.left.color = BLACK
		SUT.root.left.left.color = BLACK
		SUT.root.left.right.color = BLACK

		SUT.Remove(15)

		asserts.Nil(SUT.root.right.left)
		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(RED, SUT.root.left.color)
		asserts.Equal(RED, SUT.root.right.right.color)
	})

	t.Run("Removes a black leaf node with no children | p = BLACK | s = RED | right branch", func(t *testing.T) {
		SUT := Singleton(10, 1)
		SUT.Insert(5, 2)
		SUT.Insert(20, 3)
		SUT.Insert(1, 2)
		SUT.Insert(7, 2)
		SUT.Insert(15, 4)
		SUT.Insert(30, 5)

		// Mutate tree
		// RIGHT
		SUT.root.right.color = BLACK
		SUT.root.right.right.color = RED
		SUT.root.right.left.color = BLACK
		// LEFT
		SUT.root.left.color = BLACK
		SUT.root.left.left.color = BLACK
		SUT.root.left.right.color = BLACK

		// Balance
		SUT.root.right.right.right = &node[int, int]{parent: SUT.root.right.right, key: 40, value: 6, color: BLACK}
		SUT.root.right.right.left = &node[int, int]{parent: SUT.root.right.right, key: 25, value: 7, color: BLACK}

		SUT.Remove(15)

		asserts.Equal(30, SUT.root.right.key)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(40, SUT.root.right.right.key)
		asserts.Equal(BLACK, SUT.root.right.right.color)
		asserts.Equal(20, SUT.root.right.left.key)
		asserts.Equal(BLACK, SUT.root.right.left.color)
		asserts.Equal(25, SUT.root.right.left.right.key)
		asserts.Equal(RED, SUT.root.right.left.right.color)
		asserts.Nil(SUT.root.right.left.left)
	})

	t.Run("Removes a black leaf node with no children | p = BLACK | s = RED | left branch", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(40, 2)
		SUT.Insert(60, 3)
		SUT.Insert(70, 2)
		SUT.Insert(55, 2)
		SUT.Insert(45, 4)
		SUT.Insert(35, 5)

		// Mutate tree
		// LEFT
		SUT.root.left.color = BLACK
		SUT.root.left.left.color = RED
		SUT.root.left.right.color = BLACK
		// RIGHT
		SUT.root.right.color = BLACK
		SUT.root.right.right.color = BLACK
		SUT.root.right.left.color = BLACK

		// Balance
		SUT.root.left.left.left = &node[int, int]{parent: SUT.root.left.left, key: 20, value: 6, color: BLACK}
		SUT.root.left.left.right = &node[int, int]{parent: SUT.root.left.left, key: 37, value: 7, color: BLACK}

		SUT.Remove(45)

		asserts.Equal(35, SUT.root.left.key)
		asserts.Equal(BLACK, SUT.root.left.color)
		asserts.Equal(20, SUT.root.left.left.key)
		asserts.Equal(BLACK, SUT.root.left.left.color)
		asserts.Equal(40, SUT.root.left.right.key)
		asserts.Equal(BLACK, SUT.root.left.right.color)
		asserts.Equal(37, SUT.root.left.right.left.key)
		asserts.Equal(RED, SUT.root.left.right.left.color)
		asserts.Nil(SUT.root.left.right.right)
	})

	t.Run("DB | s = BLACK with red and black child | Left subtree", func(t *testing.T) {
		// From example https://www.youtube.com/watch?v=4KDovab_OS8&list=PLmp4WtRF6yg0_07IUb2eOsS0k1jIa2IgP&index=5&t=1819s
		SUT := Singleton(10, 1)
		SUT.Insert(5, 2)
		SUT.Insert(30, 3)
		SUT.Insert(25, 2)
		SUT.Insert(40, 2)
		SUT.Insert(7, 4)
		SUT.Insert(1, 5)

		// Mutate tree for example
		// LEFT
		SUT.root.left.left.color = BLACK
		SUT.root.left.right.color = BLACK
		// RIGHT
		SUT.root.right.right.color = BLACK
		SUT.root.right.left.color = RED

		// Manually Balance
		SUT.root.right.left.left = &node[int, int]{parent: SUT.root.right.left, key: 20, value: 6, color: BLACK}
		SUT.root.right.left.right = &node[int, int]{parent: SUT.root.right.left, key: 28, value: 7, color: BLACK}

		SUT.Remove(1)

		asserts.Equal(25, SUT.root.key)
		asserts.Equal(10, SUT.root.left.key)
		asserts.Equal(30, SUT.root.right.key)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(5, SUT.root.left.left.key)
		asserts.Equal(BLACK, SUT.root.left.left.color)
		asserts.Equal(7, SUT.root.left.left.right.key)
		asserts.Equal(RED, SUT.root.left.left.right.color)
		asserts.Equal(20, SUT.root.left.right.key)
	})

	t.Run("DB | s = BLACK with red and black child | Right subtree", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(60, 2)
		SUT.Insert(40, 3)
		SUT.Insert(45, 2)
		SUT.Insert(30, 2)
		SUT.Insert(55, 4)
		SUT.Insert(70, 5)

		// Mutate tree for testing
		// LEFT
		SUT.root.left.left.color = BLACK
		// RIGHT
		SUT.root.right.right.color = BLACK
		SUT.root.right.left.color = BLACK

		// Manually Balance
		SUT.root.left.right.right = &node[int, int]{parent: SUT.root.left.right, key: 47, value: 6, color: BLACK}
		SUT.root.left.right.left = &node[int, int]{parent: SUT.root.left.right, key: 41, value: 7, color: BLACK}

		SUT.Remove(70)

		asserts.Equal(45, SUT.root.key)
		asserts.Equal(BLACK, SUT.root.color)
		asserts.Equal(40, SUT.root.left.key)
		asserts.Equal(BLACK, SUT.root.left.color)
		asserts.Equal(50, SUT.root.right.key)
		asserts.Equal(BLACK, SUT.root.right.color)
		asserts.Equal(30, SUT.root.left.left.key)
		asserts.Equal(BLACK, SUT.root.left.left.color)
		asserts.Equal(60, SUT.root.right.right.key)
		asserts.Equal(BLACK, SUT.root.right.right.color)
		asserts.Equal(47, SUT.root.right.left.key)
		asserts.Equal(BLACK, SUT.root.right.left.color)
		asserts.Equal(55, SUT.root.right.right.left.key)
		asserts.Equal(RED, SUT.root.right.right.left.color)
		asserts.Nil(SUT.root.right.right.right)
	})
}
