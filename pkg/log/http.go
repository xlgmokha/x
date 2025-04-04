package log

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/xlgmokha/x/pkg/mapper"
)

func HTTP(logger *zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.WithContext(r.Context())

			defer func() {
				WithFields(ctx, mapper.MapFrom[*http.Request, Fields](r))
				zerolog.Ctx(ctx).Print()
			}()

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
