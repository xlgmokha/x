package x

import "slices"

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

func Reject[T any](items []T, predicate Predicate[T]) []T {
	return FindAll(items, func(item T) bool { return !predicate(item) })
}

func Contains[T any](items []T, predicate Predicate[T]) bool {
	return slices.ContainsFunc(items, predicate)
}

func All[T any](items []T, predicate Predicate[T]) bool {
	return !slices.ContainsFunc(items, func(item T) bool { return !predicate(item) })
}

func None[T any](items []T, predicate Predicate[T]) bool {
	return !slices.ContainsFunc(items, predicate)
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

func Inject[TInput any, TOutput any](items []TInput, memo TOutput, f func(TOutput, TInput) TOutput) TOutput {
	for _, item := range items {
		memo = f(memo, item)
	}
	return memo
}

func EachWithObject[TInput any, TOutput any](items []TInput, memo TOutput, f func(TInput, TOutput)) TOutput {
	return Inject(items, memo, func(memo TOutput, item TInput) TOutput {
		f(item, memo)
		return memo
	})
}

func Prepend[T any](rest []T, beginning ...T) []T {
	return slices.Concat(beginning, rest)
}

func Reverse[T any](items []T) []T {
	reversed := slices.Clone(items)
	slices.Reverse(reversed)
	return reversed
}
