package log

import (
	"bufio"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xlgmokha/x/pkg/serde"
)

func TestLog(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		log := New(writer, Fields{"env": "test"})
		log.Print()

		require.NoError(t, writer.Flush())

		items, err := serde.FromJSON[map[string]string](bufio.NewReader(&b))
		require.NoError(t, err)

		assert.Equal(t, "test", items["env"])
		assert.Equal(t, "debug", items["level"])
		assert.NotEmpty(t, items["time"])
	})

	t.Run("WithFields", func(t *testing.T) {
		var b bytes.Buffer
		writer := bufio.NewWriter(&b)

		log := New(writer, Fields{"env": "test"})
		ctx := log.WithContext(context.Background())
		WithFields(ctx, Fields{"ip": "127.0.0.1"})
		zerolog.Ctx(ctx).Print()
		log.Print()

		require.NoError(t, writer.Flush())

		items, err := serde.FromJSON[map[string]string](bufio.NewReader(&b))
		require.NoError(t, err)

		assert.Equal(t, "test", items["env"])
		assert.Equal(t, "debug", items["level"])
		assert.Equal(t, "127.0.0.1", items["ip"])
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
	})

	t.Run("Debug", func(t *testing.T) {
		t.Run("logs the item and returns it", func(t *testing.T) {
			var b bytes.Buffer
			writer := bufio.NewWriter(&b)
			log := New(writer, Fields{})

			ctx := log.WithContext(t.Context())
			item := "subject"
			result := Debug(ctx, item)

			zerolog.Ctx(ctx).Print()

			require.NoError(t, writer.Flush())

			items, err := serde.FromJSON[map[string]string](bufio.NewReader(&b))
			require.NoError(t, err)

			require.Contains(t, items, "item")
			assert.Equal(t, "subject", items["item"])

			assert.Equal(t, item, result)
		})
	})
}
