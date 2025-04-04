package log

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func Middleware(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.WithContext(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				WithFields(ctx, Fields{
					"status":      ww.Status(),
					"method":      r.Method,
					"path":        r.URL.Path,
					"remote_host": r.RemoteAddr,
				})
				zerolog.Ctx(ctx).Print()
			}()

			next.ServeHTTP(ww, r.WithContext(ctx))
		})
	}
}
