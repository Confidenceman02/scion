package elm

func Always[A any, B any](a A, _ B) A {
	return a
}
