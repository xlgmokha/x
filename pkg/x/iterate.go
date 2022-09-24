package x

type Predicate[T any] func(T) bool

func Find[T any](items []T, predicate Predicate[T]) T {
	for _, item := range items {
		if predicate(item) {
			return item
		}
	}
	return Default[T]()
}

func FindAll[T any](items []T, predicate Predicate[T]) []T {
	results := []T{}
	for _, item := range items {
		if predicate(item) {
			results = append(results, item)
		}
	}
	return results
}

func Contains[T comparable](items []T, predicate Predicate[T]) bool {
	item := Find[T](items, predicate)
	return item != Default[T]()
}
