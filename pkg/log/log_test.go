package log

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
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

	t.Run("WithMiddleware", func(t *testing.T) {
		var b bytes.Buffer
		writer := bufio.NewWriter(&b)
		log := New(writer, Fields{"env": "test"})

		server := httptest.NewServer(Middleware(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		})))
		defer server.Close()

		response, err := http.Get(server.URL)
		require.NoError(t, err)
		assert.Equal(t, http.StatusTeapot, response.StatusCode)

		require.NoError(t, writer.Flush())
		items, err := serde.FromJSON[map[string]interface{}](bufio.NewReader(&b))
		require.NoError(t, err)

		fmt.Printf("%v\n", items)
		assert.Equal(t, "test", items["env"])
		assert.Equal(t, "debug", items["level"])
		assert.NotEmpty(t, items["time"])
		assert.Equal(t, float64(http.StatusTeapot), items["status"])
		assert.Equal(t, "GET", items["method"])
		assert.Equal(t, "/", items["path"])
		assert.Contains(t, items["remote_host"], "127.0.0.1")
	})
}
