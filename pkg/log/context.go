package log

import (
	"context"

	"github.com/rs/zerolog"
)

func WithFields(ctx context.Context, fields Fields) {
	From(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Fields(fields.ToMap())
	})
}

func From(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
