package dict

import (
	"github.com/stretchr/testify/assert"
	"scion/pkg/elm/maybe"
	"testing"
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

func TestNoRotationLeft(t *testing.T) {
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
func TestLeftRotation(t *testing.T) {
	asserts := assert.New(t)

	t.Run("RR Right rotation", func(t *testing.T) {
		d := Singleton(50, 1)
		d.Insert(60, 2)
		d.Insert(70, 3)

		SUT := d.getNode(60)

		asserts.NotNil(SUT)
		asserts.Equal(BLACK, SUT.color)
		asserts.Nil(SUT.parent)
		asserts.Equal(60, SUT.key)
		asserts.Equal(70, SUT.right.key)
		asserts.Equal(50, SUT.left.key)
		asserts.Equal(RED, SUT.right.color)
		asserts.Equal(RED, SUT.left.color)
	})
}

func TestLeftRightRotation(t *testing.T) {
	asserts := assert.New(t)

	t.Run("LR rotation and balance", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(40, 2)
		SUT.Insert(45, 3)

		asserts.Equal(BLACK, SUT.rbt.color)
		asserts.Equal(45, SUT.rbt.key)
		asserts.Nil(SUT.rbt.parent)
		asserts.Equal(50, SUT.rbt.right.key)
		asserts.Equal(40, SUT.rbt.left.key)
		asserts.Equal(RED, SUT.rbt.left.color)
		asserts.Equal(RED, SUT.rbt.right.color)
	})
}

func TestRightLeftRotation(t *testing.T) {
	asserts := assert.New(t)

	t.Run("RL rotation and balance", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(60, 2)
		SUT.Insert(55, 3)

		asserts.Equal(BLACK, SUT.rbt.color)
		asserts.Equal(55, SUT.rbt.key)
		asserts.Nil(SUT.rbt.parent)
		asserts.Equal(60, SUT.rbt.right.key)
		asserts.Equal(RED, SUT.rbt.right.color)
		asserts.Equal(50, SUT.rbt.left.key)
		asserts.Equal(RED, SUT.rbt.right.color)
	})
}

func TestRemove(t *testing.T) {
	asserts := assert.New(t)

	t.Run("Removes root node with 2 children", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(60, 2)
		SUT.Insert(40, 3)
		SUT.Remove(50)

		asserts.Equal(60, SUT.rbt.key)
		asserts.Equal(2, SUT.rbt.value)
		asserts.Equal(40, SUT.rbt.left.key)
		asserts.Nil(SUT.rbt.right)
	})

	t.Run("Removes root node with no children", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Remove(50)

		asserts.Nil(SUT.rbt)
	})

	t.Run("Removes a red right leaf node", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(60, 2)
		SUT.Insert(40, 3)
		SUT.Remove(60)

		asserts.Nil(SUT.rbt.right)
	})
	t.Run("Removes a red left leaf node", func(t *testing.T) {
		SUT := Singleton(50, 1)
		SUT.Insert(40, 2)
		SUT.Insert(30, 3)
		SUT.Remove(30)

		asserts.Nil(SUT.rbt.left)
	})
	t.Run("CASE 3 - Double black node with right sibling with black children and red parent", func(t *testing.T) {
		SUT := Singleton(10, 1)
		SUT.Insert(5, 2)
		SUT.Insert(20, 3)
		SUT.Insert(15, 3)
		SUT.Insert(25, 3)
		SUT.Insert(30, 3)
		SUT.Remove(30)
		SUT.Remove(15)

		asserts.Equal(BLACK, SUT.rbt.color)
		asserts.Equal(BLACK, SUT.rbt.right.color)
		asserts.Equal(20, SUT.rbt.right.key)
		asserts.Equal(RED, SUT.rbt.right.right.color)
		asserts.Equal(25, SUT.rbt.right.right.key)
		asserts.Nil(SUT.rbt.right.left)
	})
	t.Run("CASE 3 - Double black node with left sibling with black children and red parent", func(t *testing.T) {
		SUT := Singleton(10, 1)
		SUT.Insert(5, 2)
		SUT.Insert(20, 3)
		SUT.Insert(15, 3)
		SUT.Insert(25, 3)
		SUT.Insert(30, 3)
		SUT.Remove(30)
		SUT.Remove(25)

		asserts.Equal(BLACK, SUT.rbt.color)
		asserts.Equal(BLACK, SUT.rbt.right.color)
		asserts.Equal(20, SUT.rbt.right.key)
		asserts.Equal(RED, SUT.rbt.right.left.color)
		asserts.Equal(15, SUT.rbt.right.left.key)
		asserts.Nil(SUT.rbt.right.right)
	})

	t.Run("CASE 3 - Double black node with left sibling with black children and black parent", func(t *testing.T) {
		SUT := Singleton(10, 1)
		SUT.Insert(5, 2)
		SUT.Insert(20, 3)
		SUT.Insert(15, 3)
		SUT.Insert(25, 3)
		SUT.Insert(7, 3)
		SUT.Insert(1, 3)

		// Mutate tree to be all black
		SUT.rbt.left.left.color = BLACK
		SUT.rbt.left.right.color = BLACK
		SUT.rbt.right.left.color = BLACK
		SUT.rbt.right.right.color = BLACK

		SUT.Remove(15)

		asserts.Equal(BLACK, SUT.rbt.color)
		asserts.Equal(BLACK, SUT.rbt.right.color)
		asserts.Equal(20, SUT.rbt.right.key)
		asserts.Equal(RED, SUT.rbt.right.right.color)
		asserts.Equal(25, SUT.rbt.right.right.key)
		asserts.Nil(SUT.rbt.right.left)
	})

	t.Run("test the following inserts 7,5,10,20,15", func(t *testing.T) {
		SUT := Singleton(7, 1)
		SUT.Insert(5, 2)
		SUT.Insert(10, 3)
		SUT.Insert(20, 3)
		SUT.Insert(15, 3)

		asserts.Equal(15, SUT.rbt.right.key)
	})
}
