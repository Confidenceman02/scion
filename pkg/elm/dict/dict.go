package dict

import (
	"cmp"
	"scion/pkg/elm"
	"scion/pkg/elm/maybe"
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

func (d *Dict[K, V]) Insert(key K, v V) {
	// There is no root, must be an empty tree
	if d.rbt == nil {
		d.rbt = &node[K, V]{key: key, value: v, color: BLACK, parent: nil, left: nil, right: nil}
	} else {
		balancedTree := insertHelp(key, v, d.rbt)
		d.rbt = balancedTree
	}
}

func insertHelp[K cmp.Ordered, V any](key K, value V, d *node[K, V]) *node[K, V] {
	nKey := d.key
	switch elm.Compare(key, nKey) {
	case elm.LT:
		if d.left == nil {
			d.left = &node[K, V]{key: key, value: value, color: RED, parent: d, left: nil, right: nil}
			return d
		}
	case elm.EQ:
		d.value = value
		return d
	case elm.GT:
		d.right = &node[K, V]{key: key, value: value, color: RED, parent: d, left: nil, right: nil}
		return d

	}
	panic("unreachable")
}
