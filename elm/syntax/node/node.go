package node

import "scion/elm/syntax/position"

type Node[T any] struct {
	position.Range
	Type T
}
