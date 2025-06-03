package log

import (
	"context"

	"github.com/xlgmokha/x/pkg/x"
)

func Debug[T any](ctx context.Context, item T) T {
	WithFields(ctx, Fields(x.ToMap(item)))
	return item
}
