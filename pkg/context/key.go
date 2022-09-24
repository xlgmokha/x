package context

import "context"

type Key[T any] string

func (self Key[T]) With(ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, self, value)
}
