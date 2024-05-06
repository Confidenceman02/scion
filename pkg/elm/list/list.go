package list

type mapper[A any, B any] func(a A) B

func Map[A any, B any](f mapper[A, B], xs []A) []B {
	xslen := len(xs)
	xs1 := make([]B, xslen)

	for i, x := range xs {
		xs1[i] = f(x)
	}
	return xs1
}
