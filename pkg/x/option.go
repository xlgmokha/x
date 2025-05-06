package x

type Configure[T any] func(T)
type Option[T any] func(T) T
type Factory[T any] func() T

func New[T any](options ...Option[T]) T {
	item := Default[T]()
	for _, option := range options {
		item = option(item)
	}
	return item
}

func With[T any](with Configure[T]) Option[T] {
	return func(item T) T {
		with(item)
		return item
	}
}
