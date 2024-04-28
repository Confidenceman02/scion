package core

type Result[T any, E any] interface {
	result() _Result[T, E]
}

type _Result[T any, E any] struct{}

// Methods
func (r _Result[T, E]) result() _Result[T, E] {
	return r
}

type Ok[T any, E any] struct {
	_Result[T, E]
	Value T
}

type Err[T any, E any] struct {
	_Result[T, E]
	Value T
}

func ResultWith[T any, E any, R any](
	result Result[T, E],
	ok func(*Ok[T, E]) R,
	err func(*Err[T, E]) R,
) R {
	switch d := result.(type) {
	case Ok[T, E]:
		return ok(&d)

	case Err[T, E]:
		return err(&d)
	}
	panic("unreachable")
}
