package xlog

import (
	"context"
	"log/slog"
	"sync"

	"github.com/xlgmokha/x/pkg/xcontext"
)

var key = xcontext.Key[*scope]("xlog")

type scope struct {
	logger *slog.Logger
	mu     sync.Mutex
	attrs  []slog.Attr
}

func Into(ctx context.Context, logger *slog.Logger) context.Context {
	return key.With(ctx, &scope{logger: logger})
}

func From(ctx context.Context) *slog.Logger {
	if s := key.From(ctx); s != nil {
		return s.logger
	}
	return slog.Default()
}

func WithFields(ctx context.Context, fields Fields) {
	WithAttrs(ctx, fields.ToAttrs()...)
}

func WithAttrs(ctx context.Context, attrs ...slog.Attr) {
	s := key.From(ctx)
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.attrs = append(s.attrs, attrs...)
}

// collect builds a fresh slice rather than appending to base, so that a base
// with spare capacity, which an AttrsFor override is free to return, is never
// written to.
func (s *scope) collect(base []slog.Attr) []slog.Attr {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.attrs) == 0 {
		return base
	}

	attrs := make([]slog.Attr, 0, len(base)+len(s.attrs))
	attrs = append(attrs, base...)

	return append(attrs, s.attrs...)
}
