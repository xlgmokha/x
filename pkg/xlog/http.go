package xlog

import (
	"log/slog"
	"net/http"

	"github.com/xlgmokha/x/pkg/mapper"
)

func HTTP(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := Into(r.Context(), logger)

			defer func() {
				if !logger.Enabled(ctx, slog.LevelInfo) {
					return
				}
				logger.LogAttrs(ctx, slog.LevelInfo, "request", key.From(ctx).collect(AttrsFor(r))...)
			}()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var AttrsFor mapper.Mapping[*http.Request, []slog.Attr] = func(r *http.Request) []slog.Attr {
	return []slog.Attr{
		slog.String("host", r.URL.Host),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote_host", r.RemoteAddr),
	}
}
