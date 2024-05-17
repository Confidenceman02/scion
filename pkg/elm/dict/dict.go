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

func balance[K cmp.Ordered, V any](n *node[K, V]) *node[K, V] {
	// Root node
	if n.parent == nil {
		n.color = BLACK
		return n
	}
	// parent is Black no more work to do
	if n.parent.color == BLACK {
		return nil
	}
	parent := n.parent
	gDir := grampsSide(parent)
	pDir := parentSide(n)
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
		if uncle == nil || uncle.color == BLACK {
		}
		switch pDir {
		case RIGHT:
			// RR case -> Left rotation
			n.rrRotation()

			// If root, return new root
			if parent.parent != nil {
				return nil
			} else {
				return parent
			}
		case LEFT:
			// RL -> Single right rotation -> balance
			newTarget := n.parent.srRotation()
			return balance(newTarget)
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

			switch pDir {
			case LEFT:
				// LL rotation
				n.llRotation()

				// If root, return new root
				if parent.parent != nil {
					return nil
				} else {
					return parent
				}
			case RIGHT:
				// LR case -> single right rotation -> balance
				newTarget := n.parent.slRotation()
				return balance(newTarget)
			}
		}
	}
	return nil
}
func (x *node[K, V]) srRotation() *node[K, V] {
	grandparent := x.parent
	left := x.left

	// 1. left becomes new parent
	left.parent = x.parent

	// 2. x's parent is now left
	x.parent = left

	// 3. x's left is now lefts right
	x.left = left.right

	// 4. left's right is x
	left.right = x

	// 5. Grandparent's right is now left
	grandparent.right = left

	return x
}

func (x *node[K, V]) slRotation() *node[K, V] {
	grandparent := x.parent
	right := x.right
	// 1. right becomes new parent
	right.parent = x.parent

	// 2. x's parent is now right
	x.parent = right

	// 3. x's right is right's left
	x.right = right.left

	// 4. right's left is x
	right.left = x

	// 5. Grandparent's left now points to right
	grandparent.left = right

	return x
}

// func (n *node[K, V]) lRotation() {
// 	parent := n.parent
// 	grandparent := parent.parent
//
// 	// 1. n becomes new parent
// 	n.parent = grandparent
//
// 	// 2. Parents parent is now n
// 	parent.parent = n
//
// 	// 3. Gramps's right is n
// 	grandparent.right = n
//
// 	// 4. Parents left is n's right
// 	parent.left = n.right
//
// 	// 5. n's right is now parent
// 	n.right = parent
// }

// func (n *node[K, V]) rRotation() {
// 	parent := n.parent
// 	grandparent := parent.parent
//
// 	// 1. n becomes new parent
// 	n.parent = grandparent
//
// 	// 2. Parents parent is now n
// 	parent.parent = n
//
// 	// 3. Gramps's left is n
// 	grandparent.left = n
//
// 	// 4. Parents right is n's left
// 	parent.right = n.left
//
// 	// 5. n's left is now parent
// 	n.left = parent
// }

func (n *node[K, V]) rrRotation() {
	parent := n.parent
	grandparent := parent.parent

	// 1. Parent gets gramps's parent
	parent.parent = grandparent.parent

	// 2. Gramps's parent is now parent
	grandparent.parent = parent

	// 3. Gramps left child is now parents right child
	grandparent.right = parent.left

	// 4. Parents right child is now gramps
	parent.left = grandparent

	// 5. Parent and gramps swap colors
	pColor := parent.color
	gColor := grandparent.color

	parent.color = gColor
	grandparent.color = pColor
}

func (x *node[K, V]) llRotation() {
	parent := x.parent
	grandparent := parent.parent

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
