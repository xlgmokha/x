package x

type Mapper[T any, K any] func(T) K
type Predicate[T any] func(T) bool
type Visitor[T any] func(T)

func Find[T any](items []T, predicate Predicate[T]) T {
	for _, item := range items {
		if predicate(item) {
			return item
		}
	}
	return Zero[T]()
}

func FindAll[T any](items []T, predicate Predicate[T]) []T {
	results := []T{}
	Each[T](items, func(item T) {
		if predicate(item) {
			results = append(results, item)
		}
	})
	return results
}

func Contains[T comparable](items []T, predicate Predicate[T]) bool {
	return Find[T](items, predicate) != Zero[T]()
}

func Map[TInput any, TOutput any](items []TInput, mapFrom Mapper[TInput, TOutput]) []TOutput {
	results := []TOutput{}
	Each[TInput](items, func(item TInput) {
		results = append(results, mapFrom(item))
	})
	return results
}

func Each[T any](items []T, v Visitor[T]) {
	for _, item := range items {
		v(item)
	}
}

func Inject[TInput any, TOutput any](items []TInput, memo TOutput, f func(TOutput, TInput)) TOutput {
	for _, item := range items {
		f(memo, item)
	}
	return memo
}
