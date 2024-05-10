package elm

func Always[A any, B any](a A, _ B) A {
	return a
}

func Identity[A any](a A) A {
	return a
}

type Never struct{}
