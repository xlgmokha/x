package x

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/example", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	passThrough := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("x-pass-through", "true")
			next.ServeHTTP(w, r)
		})
	}

	unauthorized := func(http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("x-unauthorized", "true")
			w.WriteHeader(http.StatusUnauthorized)
		})
	}

	t.Run("executes a single middleware", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/example", nil)
		w := httptest.NewRecorder()

		Middleware[http.Handler](mux, passThrough).ServeHTTP(w, r)

		require.Equal(t, http.StatusTeapot, w.Code)
		assert.Equal(t, "true", w.HeaderMap.Get("x-pass-through"))
	})

	t.Run("excutes multiple middleware", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/example", nil)
		w := httptest.NewRecorder()

		Middleware[http.Handler](mux, passThrough, unauthorized).ServeHTTP(w, r)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Equal(t, "true", w.HeaderMap.Get("x-pass-through"))
		assert.Equal(t, "true", w.HeaderMap.Get("x-unauthorized"))
	})
}
