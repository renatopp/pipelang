package internal

// Used for general stack data structure
type Stack[T any] struct {
	elements []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		elements: make([]T, 0),
	}
}

func (s *Stack[T]) Push(e T) {
	s.elements = append(s.elements, e)
}

func (s *Stack[T]) Peek() T {
	return s.elements[len(s.elements)-1]
}
func (s *Stack[T]) PeekOr(v T) T {
	if len(s.elements) == 0 {
		return v
	}

	return s.elements[len(s.elements)-1]
}

func (s *Stack[T]) Pop() T {
	e := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return e
}

func (s *Stack[T]) Len() int {
	return len(s.elements)
}

func (s *Stack[T]) Set(e T) {
	s.elements[len(s.elements)-1] = e
}
