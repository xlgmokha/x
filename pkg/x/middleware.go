package x

func Middleware[T any](item T, items ...Option[T]) T {
	for _, middleware := range Reverse(items) {
		item = middleware(item)
	}
	return item
}
