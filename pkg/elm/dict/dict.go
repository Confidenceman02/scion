package dict

/*
Elm Dicts under the hood are Red-Black trees.

Insertions are standard BST insertions.

The rules are as follows for Red-Black trees:

    1. Every node is colored red or black.
    2. Every leaf is a NIL node, and is colored black.
    3. If a node is red, then both its children are black.
    4. Every simple path from a node to a descendant leaf contains the same number of black nodes.
*/

type Dict[K comparable, V any] interface {
	dict() _dict
}
type _dict struct{}

func (r _dict) dict() _dict {
	return r
}

// Variants
type empty[K comparable] struct {
	_dict
}

type node[K comparable, V any] struct {
	_dict
	key    K
	value  V
	color  ncolor
	parent Dict[K, V]
	left   Dict[K, V]
	right  Dict[K, V]
}

// Node color
type ncolor int

const (
	RED ncolor = iota
	BLACK
)

// Methods
func Empty[K comparable, V any]() Dict[K, V] {
	e := empty[K]{}
	return e
}

func Singleton[K comparable, V any](key K, value V) Dict[K, V] {
	// Root nodes are always black
	return node[K, V]{key: key, value: value, color: BLACK, parent: empty[K]{}, left: empty[K]{}, right: empty[K]{}}
}

func IsEmpty[K comparable, V any](d Dict[K, V]) bool {
	return dictWith(
		d,
		func(*empty[K]) bool { return true },
		func(*node[K, V]) bool { return false },
	)
}

func Insert[K comparable, V any](key K, v V, d Dict[K, V]) Dict[K, V] {
	switch n := insertHelp(key, v, d).(type) {
	case node[K, V]:
		if n.color == RED {
			n.color = BLACK
			return n
		}
	default:
		return n
	}
	panic("unreachable")
}

func insertHelp[K comparable, V any](key K, value V, dict Dict[K, V]) Dict[K, V] {
	switch dict.(type) {
	case empty[K]:
		return node[K, V]{key: key, value: value, color: RED, parent: empty[K]{}, left: empty[K]{}, right: empty[K]{}}
	case node[K, V]:
		return node[K, V]{key: key, value: value, color: RED, parent: empty[K]{}, left: empty[K]{}, right: empty[K]{}}
	}
	panic("unreachable")
}

// Matcher
func dictWith[K comparable, V any, R any](
	d Dict[K, V],
	e func(*empty[K]) R,
	n func(*node[K, V]) R,
) R {
	switch d := d.(type) {
	case empty[K]:
		return e(&d)
	case node[K, V]:
		return n(&d)
	}
	panic("unreachable")
}
