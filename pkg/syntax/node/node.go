package node

import "scion/pkg/syntax/position"

type Node[T any] struct {
	position.Range
	Type T
}
