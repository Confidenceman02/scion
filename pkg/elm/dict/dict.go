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

type Dict[K comparable] interface {
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

type node[K comparable] struct {
	_dict
	key    K
	color  ncolor
	parent Dict[K]
	left   Dict[K]
	right  Dict[K]
}

// Node color
type ncolor int

const (
	RED ncolor = iota
	BLACK
)

// Methods
func Empty[K comparable]() Dict[K] {
	e := empty[K]{}
	return e
}

func Singleton[K comparable](key K) Dict[K] {
	// Root nodes are always black
	return node[K]{key: key, color: BLACK, parent: empty[K]{}, left: empty[K]{}, right: empty[K]{}}
}

func IsEmpty[K comparable](d Dict[K]) bool {
	return dictWith(
		d,
		func(*empty[K]) bool { return true },
		func(*node[K]) bool { return false },
	)
}

// Matcher
func dictWith[K comparable, R any](
	d Dict[K],
	e func(*empty[K]) R,
	n func(*node[K]) R,
) R {
	switch d := d.(type) {
	case empty[K]:
		return e(&d)
	case node[K]:
		return n(&d)
	}
	panic("unreachable")
}
