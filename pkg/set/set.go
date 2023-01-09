package set

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](vals ...T) Set[T] {
	s := Set[T]{}
	for _, val := range vals {
		s.Add(val)
	}
	return s
}

func (s Set[T]) Has(val T) bool {
	_, exists := s[val]
	return exists
}

func (s *Set[T]) Add(val T) {
	if s == nil {
		newMap := map[T]struct{}{}
		*s = newMap
	}

	(*s)[val] = struct{}{}
}
