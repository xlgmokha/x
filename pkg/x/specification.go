package x

type Specification[T any] interface {
	IsSatisfiedBy(T) bool
}

func (f Predicate[T]) IsSatisfiedBy(item T) bool {
	return f(item)
}

func (s Predicate[T]) Or(other Specification[T]) Predicate[T] {
	return Predicate[T](func(item T) bool {
		return s.IsSatisfiedBy(item) || other.IsSatisfiedBy(item)
	})
}

func (s Predicate[T]) And(other Specification[T]) Predicate[T] {
	return Predicate[T](func(item T) bool {
		return s.IsSatisfiedBy(item) && other.IsSatisfiedBy(item)
	})
}
