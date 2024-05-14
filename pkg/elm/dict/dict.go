package dict

import (
	"cmp"
	"scion/pkg/elm"
	"scion/pkg/elm/maybe"
)

const (
	LEFT = iota
	RIGHT
	NONE
)

/*
Elm Dicts under the hood are Red-Black trees.

Insertions are standard BST insertions.

The rules are as follows for Red-Black trees:

    1. Every node is colored red or black.
    2. Every leaf is a NIL node, and is colored black.
    3. If a node is red, then both its children are black.
    4. Every simple path from a node to a descendant leaf contains the same number of black nodes.
*/

type Dict[K cmp.Ordered, V any] struct {
	rbt *node[K, V]
}

// Variants

type node[K cmp.Ordered, V any] struct {
	key    K
	value  V
	color  ncolor
	parent *node[K, V]
	left   *node[K, V]
	right  *node[K, V]
}

// Node color
type ncolor int

const (
	RED ncolor = iota
	BLACK
)

// Builders
func Empty[K cmp.Ordered, V any]() Dict[K, V] {
	return Dict[K, V]{rbt: nil}
}

func Singleton[K cmp.Ordered, V any](key K, value V) *Dict[K, V] {
	// Root nodes are always black
	return &Dict[K, V]{
		rbt: &node[K, V]{
			key:    key,
			value:  value,
			color:  BLACK,
			parent: nil,
			left:   nil,
			right:  nil},
	}
}

// Methods
func (d *Dict[K, V]) Get(targetKey K) maybe.Maybe[V] {
	if d.rbt == nil {
		return maybe.Nothing{}
	} else {
		return getHelp(targetKey, d.rbt)
	}
}

func getHelp[K cmp.Ordered, V any](targetKey K, n *node[K, V]) maybe.Maybe[V] {
	if n != nil {
		switch elm.Compare(targetKey, n.key) {
		case elm.LT:
			return getHelp(targetKey, n.left)
		case elm.EQ:
			return maybe.Just[V]{Value: n.value}
		case elm.GT:
			return getHelp(targetKey, n.right)
		}
	}
	return maybe.Nothing{}
}

func (d *Dict[K, V]) getNode(targetKey K) *node[K, V] {
	if d.rbt != nil {
		return getNodeHelp(targetKey, d.rbt)
	}
	return d.rbt
}

func getNodeHelp[K cmp.Ordered, V any](targetKey K, n *node[K, V]) *node[K, V] {
	if n != nil {
		switch elm.Compare(targetKey, n.key) {
		case elm.LT:
			return getNodeHelp(targetKey, n.left)
		case elm.EQ:
			return n
		case elm.GT:
			return getNodeHelp(targetKey, n.right)
		}
	}
	return nil
}

func (d *Dict[K, V]) Insert(key K, v V) {
	// There is no root, must be an empty tree
	if d.rbt == nil {
		d.rbt = &node[K, V]{key: key, value: v, color: RED, parent: nil, left: nil, right: nil}
		balance(d.rbt)
	} else {
		newRoot := balance(insertHelp(key, v, d.rbt))
		if newRoot != nil {
			d.rbt = newRoot
		}
	}
}

/*
Returns the inserted node and it's inserted direction
*/
func insertHelp[K cmp.Ordered, V any](key K, value V, d *node[K, V]) *node[K, V] {
	nKey := d.key
	switch elm.Compare(key, nKey) {
	case elm.LT:
		if d.left == nil {
			d.left = &node[K, V]{key: key, value: value, color: RED, parent: d, left: nil, right: nil}
			return d.left
		} else {
			return insertHelp(key, value, d.left)
		}
	case elm.EQ:
		d.value = value
		return d
	case elm.GT:
		if d.right == nil {
			d.right = &node[K, V]{key: key, value: value, color: RED, parent: d, left: nil, right: nil}
			return d.right
		} else {
			return insertHelp(key, value, d.right)

		}
	}
	panic("unreachable")
}

func balance[K cmp.Ordered, V any](n *node[K, V]) *node[K, V] {
	// Root node
	if n.parent == nil {
		n.color = BLACK
		return nil
	}
	// parent is Black no more work to do
	if n.parent.color == BLACK {
		return nil
	}
	parent := n.parent
	gDir := grampsSide(parent)
	grandparent := n.parent.parent

	switch gDir {
	case RIGHT:
		uncle := grandparent.left

		// Handle no rotation
		if uncle != nil && uncle.color == RED {
			parent.color = BLACK
			uncle.color = BLACK
			grandparent.color = RED
			return balance(grandparent)
		}

	case LEFT:
		uncle := grandparent.right
		// Handle no rotation case
		if uncle != nil && uncle.color == RED {
			parent.color = BLACK
			uncle.color = BLACK
			grandparent.color = RED
			return balance(grandparent)
		}
		if uncle == nil || uncle.color == BLACK {
			pDir := parentSide(n)

			switch pDir {
			case LEFT:
				// LL case -> Right rotation

				// 1. Parent gets gramps's parent
				parent.parent = grandparent.parent

				// 2. Gramps's parent is now parent
				grandparent.parent = parent

				// 3. Gramps left child is now parents right child
				grandparent.left = parent.right

				// 4. Parents right child is now gramps
				parent.right = grandparent

				// 5. Parent and gramps swap colors
				pColor := parent.color
				gColor := grandparent.color

				parent.color = gColor
				grandparent.color = pColor

				// If root, return new root
				if parent.parent != nil {
					return nil
				} else {
					return parent
				}
			}
		}
	}
	return nil
}

func parentSide[K cmp.Ordered, V any](x *node[K, V]) int {
	parent := x.parent
	if parent.left != nil && x.key == parent.left.key {
		return LEFT
	} else {
		return RIGHT
	}
}

func grampsSide[K cmp.Ordered, V any](parentNode *node[K, V]) int {
	grandparent := parentNode.parent
	if grandparent.left != nil && parentNode.key == grandparent.left.key {
		return LEFT
	} else {
		return RIGHT
	}
}
