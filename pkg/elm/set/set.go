package set

type set[C comparable] struct {
	s[C]
}

type Set[C comparable] *set[C]

type s[C comparable] map[C]struct{}

func FromList[C comparable](s []C) Set[C] {
	newSet := make(map[C]struct{})
	for _, v := range s {
		newSet[v] = struct{}{}
	}
	return &set[C]{s: newSet}
}

func ToList[C comparable](s Set[C]) []C {
	newList := make([]C, len(s.s))
	i := 0
	for k := range s.s {
		newList[i] = k
		i++
	}
	return newList
}
