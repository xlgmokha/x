package xcontext

import (
	"context"

	"github.com/xlgmokha/x/pkg/x"
)

type Key[T any] string

func (self Key[T]) With(ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, self, value)
}

func (self Key[T]) From(ctx context.Context) T {
	if value, ok := ctx.Value(self).(T); ok {
		return value
	}
	return x.Zero[T]()
}
