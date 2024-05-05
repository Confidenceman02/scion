package elm

/*
A `Result` is the result of a computation that may fail. This is a great way to manage errors in Scion.
*/

type Result[T any, E any] interface {
	result() _Result[T, E]
}

type _Result[T any, E any] struct{}

func (r _Result[T, E]) result() _Result[T, E] {
	return r
}

/*
A `Result` is `Ok` meaning the computation succeeded
*/
type Ok[T any, E any] struct {
	_Result[T, E]
	Value T
}

/*
A `Result` is `Err` meaning the computation failed
*/
type Err[T any, E any] struct {
	_Result[T, E]
	Value E
}

/*
A handy way to pattern match Result types
*/
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
