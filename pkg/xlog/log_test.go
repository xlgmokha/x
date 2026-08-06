package xlog

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xlgmokha/x/pkg/serde"
)

func TestLog(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		log := New(writer, Fields{"env": "test"})
		log.Info("hello")

		require.NoError(t, writer.Flush())

		items, err := serde.FromJSON[map[string]string](bufio.NewReader(&b))
		require.NoError(t, err)

		assert.Equal(t, "test", items["env"])
		assert.Equal(t, "INFO", items["level"])
		assert.Equal(t, "hello", items["msg"])
		assert.NotEmpty(t, items["time"])
	})

	t.Run("From", func(t *testing.T) {
		t.Run("returns the logger in the context", func(t *testing.T) {
			log := New(bufio.NewWriter(&bytes.Buffer{}), Fields{})

			assert.Same(t, log, From(Into(context.Background(), log)))
		})

		t.Run("falls back to the default logger", func(t *testing.T) {
			assert.NotNil(t, From(context.Background()))
		})
	})

	t.Run("WithFields", func(t *testing.T) {
		t.Run("is collected into the request record", func(t *testing.T) {
			var b bytes.Buffer
			writer := bufio.NewWriter(&b)
			log := New(writer, Fields{"env": "test"})

			handler := HTTP(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				WithFields(r.Context(), Fields{"ip": "127.0.0.1"})
			}))
			r, w := httptest.NewRequest("GET", "/example", nil), httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			require.NoError(t, writer.Flush())

			items, err := serde.FromJSON[map[string]interface{}](bufio.NewReader(&b))
			require.NoError(t, err)

			assert.Equal(t, "test", items["env"])
			assert.Equal(t, "127.0.0.1", items["ip"])
		})

		t.Run("is a no-op outside a request scope", func(t *testing.T) {
			assert.NotPanics(t, func() {
				WithFields(context.Background(), Fields{"ip": "127.0.0.1"})
			})
		})
	})

	t.Run("HTTP", func(t *testing.T) {
		var b bytes.Buffer
		writer := bufio.NewWriter(&b)
		log := New(writer, Fields{})
		server := httptest.NewServer(HTTP(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		})))
		defer server.Close()

		response, err := http.Get(server.URL)
		require.NoError(t, err)
		assert.Equal(t, http.StatusTeapot, response.StatusCode)

		require.NoError(t, writer.Flush())
		items, err := serde.FromJSON[map[string]interface{}](bufio.NewReader(&b))
		require.NoError(t, err)

		require.Contains(t, items, "method")
		assert.Equal(t, "GET", items["method"])

		require.Contains(t, items, "path")
		assert.Equal(t, "/", items["path"])

		require.Contains(t, items, "remote_host")
		assert.Contains(t, items["remote_host"], "127.0.0.1")

		assert.Equal(t, "request", items["msg"])
	})
}

func TestAttrsFor(t *testing.T) {
	t.Run("can be overridden by the caller", func(t *testing.T) {
		original := AttrsFor
		defer func() { AttrsFor = original }()

		AttrsFor = func(r *http.Request) []slog.Attr {
			return []slog.Attr{slog.String("verb", r.Method)}
		}

		var b bytes.Buffer
		writer := bufio.NewWriter(&b)
		handler := HTTP(New(writer, Fields{}))(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/example", nil))

		require.NoError(t, writer.Flush())
		items, err := serde.FromJSON[map[string]interface{}](bufio.NewReader(&b))
		require.NoError(t, err)

		assert.Equal(t, "GET", items["verb"])
		assert.NotContains(t, items, "method")
		assert.NotContains(t, items, "remote_host")
	})

	t.Run("never writes into the slice it returns", func(t *testing.T) {
		original := AttrsFor
		defer func() { AttrsFor = original }()

		shared := append(make([]slog.Attr, 0, 8), slog.String("app", "example"))
		AttrsFor = func(*http.Request) []slog.Attr { return shared }

		handler := HTTP(New(io.Discard, Fields{}))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WithAttrs(r.Context(), slog.String("id", r.URL.Path))
		}))

		var wg sync.WaitGroup
		for range 32 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/example", nil))
			}()
		}
		wg.Wait()

		assert.Equal(t, []slog.Attr{slog.String("app", "example")}, shared)
	})
}
