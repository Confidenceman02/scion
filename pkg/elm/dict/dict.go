package dict

import (
	"cmp"
	"scion/pkg/elm"
	"scion/pkg/elm/maybe"
)

const (
	LEFT = iota
	RIGHT
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

type node[K cmp.Ordered, V any] struct {
	key    K
	value  V
	color  int
	parent *node[K, V]
	left   *node[K, V]
	right  *node[K, V]
}

const (
	RED int = iota
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
			if n.left != nil {
				return getNodeHelp(targetKey, n.left)
			}
			return nil
		case elm.EQ:
			return n
		case elm.GT:
			if n.right != nil {
				return getNodeHelp(targetKey, n.right)
			}
			return nil
		}
	}
	return nil
}

func (d *Dict[K, V]) Insert(key K, v V) {
	balance(d, insertHelp(key, v, d, d.rbt))
}

func insertHelp[K cmp.Ordered, V any](key K, value V, dict *Dict[K, V], n *node[K, V]) *node[K, V] {
	if dict.rbt == nil {
		dict.rbt = &node[K, V]{key: key, value: value, color: RED, parent: nil, left: nil, right: nil}
		return dict.rbt
	} else {
		nKey := n.key
		switch elm.Compare(key, nKey) {
		case elm.LT:
			if n.left == nil {
				n.left = &node[K, V]{key: key, value: value, color: RED, parent: n, left: nil, right: nil}
				return n.left
			} else {
				return insertHelp(key, value, dict, n.left)
			}
		case elm.EQ:
			n.value = value
			return n
		case elm.GT:
			if n.right == nil {
				n.right = &node[K, V]{key: key, value: value, color: RED, parent: n, left: nil, right: nil}
				return n.right
			} else {
				return insertHelp(key, value, dict, n.right)
			}
		}
	}
	panic("unreachable")
}

/*
Removal is a bit more of a process to that of insertion

Case 1 - Node is a red leaf
    1.1
        Delete node and exit
Case 2 - Double Black (DB) is root
    2.2
        Remove DB
Case 3 - DB sibling is black and both nephews are black
    3.1
        Remove DB node
    3.2
        Make sibling red
    3.3
        Add black to parent. If parent was red, make black
        otherwise make it a DB and find appropriate CASE
Case 4 - DB sibling is red
    4.1 Swap colors of DB parent & sibling
    4.2 Rotate parent in DB's direction
    4.3 Find next case for DB

Case 5 - DB sibling is black, far nephew is black and ner nephew is red
    5.1 Swap colors of the DB sibling and near nephew
    5.2 Rotate sibling of DB node in opposite direction of DB node
    5.3 Apply case 6
Case 6 - DB sibling is black and far nephew is red
    6.1 Swap the colors of the DB parent and sibling
    6.2 Rotate DB parent in DB direction
    6.3 Turn far nephews color to black
    6.4 Remove DB node to single black
*/

func (d *Dict[K, V]) Remove(key K) {
	if d.rbt == nil {
		// Empty tree
		return
	} else {
		// Find node to delete
		dn := d.getNode(key)
		if dn == nil {
			// Node not found
			return
		}
		if dn.parent == nil && dn.blackChildren() {
			// Root is leaf
			d.rbt = nil
			return
		}
		removeHelp(dn)
	}
}

func removeHelp[K cmp.Ordered, V any](n *node[K, V]) {
	if n.left == nil && n.right == nil {
		if n.color == RED {
			// CASE 1
			switch parentSide(n) {
			case RIGHT:
				n.parent.right = nil
			case LEFT:
				n.parent.left = nil
			}
			return
		} else {
			// Black leaf (DB)
			sibling, siblingSide := findSiblingWithSide(n)
			switch siblingSide {
			case LEFT:
				n.parent.right = nil
				fixDB(n, sibling)
			case RIGHT:
				n.parent.left = nil
				fixDB(n, sibling)
			}
			return
		}
	}
	if n.left != nil && n.right != nil {
		// Internal node || root - find successor -> swap -> delete successor
		successor := findMin(n.right)
		n.key = successor.key
		n.value = successor.value
		removeHelp(successor)
		return
	}
}

func fixDB[K cmp.Ordered, V any](db *node[K, V], sibling *node[K, V]) {
	// CASE 2
	if db != nil && db.parent == nil {
		// DB is root, nothing more to do
		return
	}
	// CASE 3
	if sibling.color == BLACK && sibling.blackChildren() {
		sibling.color = RED
		if sibling.parent.color == BLACK && sibling.parent.parent != nil {
			sib1, _ := findSiblingWithSide(sibling.parent)
			// Another (DB)
			fixDB(sibling.parent, sib1)
			return
		}
		sibling.parent.color = BLACK
		return
	}
	// TODO Sibling with non black children
}

func (n *node[K, V]) blackChildren() bool {
	if n != nil {
		return (n.left == nil || n.left.color == BLACK) && (n.right == nil || n.right.color == BLACK)
	}
	return false
}

func findSiblingWithSide[K cmp.Ordered, V any](n *node[K, V]) (*node[K, V], int) {
	parent := n.parent
	if parent.left.key == n.key {
		return parent.right, RIGHT
	}
	return parent.left, LEFT
}

func findMin[K cmp.Ordered, V any](n *node[K, V]) *node[K, V] {
	if n.left == nil {
		return n
	}
	return findMin(n.left)
}

func balance[K cmp.Ordered, V any](dict *Dict[K, V], n *node[K, V]) {
	// Root case
	if n.parent == nil {
		n.color = BLACK
		dict.rbt = n
		return
	}
	pColor := n.parent.color
	if pColor == BLACK {
		// Nothing more to do
		return
	}
	// Parent and n are red
	nDir := parentSide(n)
	pDir := parentSide(n.parent)

	switch nDir {
	case LEFT:
		switch pDir {
		case LEFT:
			// LL Red - right rotate on grandparent
			newRoot := n.parent.parent.srRotation()
			// Push down newRoot color
			newRoot.right.color = newRoot.color
			// balance newRoot
			balance(dict, newRoot)
			return
		}
	case RIGHT:
		switch pDir {
		case RIGHT:
			// RR Red - left rotate on grandparent
			newRoot := n.parent.parent.slRotation()
			// Push down newRoot color
			newRoot.left.color = newRoot.color
			// balance newRoot
			balance(dict, newRoot)
			return
		}
	}
	// We have either a LR or RL
	uncle := getUncle(n)
	grandparent := n.parent.parent

	if uncle != nil && uncle.color == RED {
		// Red uncle - push down blackness from root - balance root
		uncle.color = grandparent.color
		n.parent.color = grandparent.color
		grandparent.color = RED
		balance(dict, grandparent)
		return
	}

	// Black uncle - rotate - balance
	switch nDir {
	case RIGHT:
		// Rotate parent left
		newRoot := n.parent.slRotation()
		// balance outer grandchild
		balance(dict, newRoot.left)
		return

	case LEFT:
		// Rotate parent right
		newRoot := n.parent.srRotation()
		// balance outer grandchild
		balance(dict, newRoot.right)
		return
	}
}

func (x *node[K, V]) srRotation() *node[K, V] {
	left := x.left

	// Handle x's parent
	if x.parent != nil {
		pSide := parentSide(x)

		switch pSide {
		case LEFT:
			x.parent.left = left
		case RIGHT:
			x.parent.right = left
		}
	}

	// 1. left becomes new parent
	left.parent = x.parent

	// 2. x's parent is now left
	x.parent = left

	// 3. x's left is now lefts right
	x.left = left.right

	// 4. left's right is x
	left.right = x

	return left
}

func (x *node[K, V]) slRotation() *node[K, V] {
	right := x.right

	// Handle x's parent
	if x.parent != nil {
		pSide := parentSide(x)

		switch pSide {
		case LEFT:
			x.parent.left = right
		case RIGHT:
			x.parent.right = right
		}
	}
	// 1. right becomes new parent
	right.parent = x.parent

	// 2. x's parent is now right
	x.parent = right

	// 3. x's right is right's left
	x.right = right.left

	// 4. right's left is x
	right.left = x

	return right
}

func parentSide[K cmp.Ordered, V any](n *node[K, V]) int {
	parent := n.parent
	if parent.left != nil && n.key == parent.left.key {
		return LEFT
	} else {
		return RIGHT
	}
}

func getUncle[K cmp.Ordered, V any](n *node[K, V]) *node[K, V] {
	grandparent := n.parent.parent
	parent := n.parent

	switch parentSide(parent) {
	case LEFT:
		// Uncle is right side
		if grandparent.right != nil {
			return grandparent.right
		} else {
			return nil
		}
	case RIGHT:
		if grandparent.left != nil {
			return grandparent.left
		} else {
			return nil
		}
	}
	panic("unreachable")
}
