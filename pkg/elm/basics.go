package elm

import "cmp"

const (
	EQ = 0
	LT = -1
	GT = +1
)

func Always[A any, B any](a A, _ B) A {
	return a
}

func Identity[A any](a A) A {
	return a
}

type Never struct{}

func Compare[A cmp.Ordered](a, b A) int {
	return cmp.Compare(a, b)
}
