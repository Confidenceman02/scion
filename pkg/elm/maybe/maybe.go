package maybe

type Maybe[V any] interface {
	maybe() _maybe
}

type _maybe struct{}

func (m _maybe) maybe() _maybe {
	return m
}

// Variants
type Just[V any] struct {
	_maybe
	Value V
}
type Nothing struct {
	_maybe
}

// Matcher
func Match[V any, R any](
	m Maybe[V],
	j func(*Just[V]) R,
	n func(*Nothing) R,
) R {
	switch m := m.(type) {
	case Just[V]:
		return j(&m)
	case Nothing:
		return n(&m)
	}
	panic("unreachable")
}
