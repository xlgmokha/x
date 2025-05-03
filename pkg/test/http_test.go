package test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xlgmokha/x/pkg/context"
	"github.com/xlgmokha/x/pkg/serde"
)

var exampleHeader RequestOption = RequestOption(func(r *http.Request) *http.Request {
	r.Header.Add("X-Example", "example")
	return r
})

var withHost RequestOption = RequestOption(func(r *http.Request) *http.Request {
	r.Host = "example.com"
	return r
})

func TestRequest(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		r := Request("GET", "/example")

		require.NotNil(t, r)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/example", r.URL.Path)
		assert.Zero(t, r.Body)
	})

	t.Run("with an option", func(t *testing.T) {
		r := Request("GET", "/example", exampleHeader)

		require.NotNil(t, r)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/example", r.URL.Path)
		assert.Zero(t, r.Body)
		assert.Equal(t, "example", r.Header.Get("X-Example"))
	})

	t.Run("with options", func(t *testing.T) {
		r := Request("GET", "/example", exampleHeader, withHost)

		require.NotNil(t, r)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/example", r.URL.Path)
		assert.Zero(t, r.Body)
		assert.Equal(t, "example", r.Header.Get("X-Example"))
		assert.Equal(t, "example.com", r.Host)
	})
}

func TestRequestResponse(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		r, w := RequestResponse("GET", "/health")

		require.NotNil(t, r)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/health", r.URL.Path)

		require.NotNil(t, w)
	})

	t.Run("with options", func(t *testing.T) {
		r, w := RequestResponse("GET", "/example", exampleHeader, withHost)

		require.NotNil(t, r)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/example", r.URL.Path)
		assert.Zero(t, r.Body)
		assert.Equal(t, "example", r.Header.Get("X-Example"))
		assert.Equal(t, "example.com", r.Host)

		require.NotNil(t, w)
	})
}

func TestWithAcceptHeader(t *testing.T) {
	t.Run("applies the Accept header", func(t *testing.T) {
		r := Request("GET", "/example", WithAcceptHeader(serde.JSON))

		require.NotNil(t, r)
		assert.Equal(t, serde.JSON.String(), r.Header.Get("Accept"))
	})
}

func TestWithRequestHeader(t *testing.T) {
	t.Run("applies a header to the request", func(t *testing.T) {
		r := Request("GET", "/example", WithRequestHeader("Via", "gtwy"))

		require.NotNil(t, r)
		assert.Equal(t, "gtwy", r.Header.Get("Via"))
	})
}

func TestWithContentType(t *testing.T) {
	type example struct {
		ID   int    `json: "id" jsonapi:"primary,examples"`
		Name string `json: "name" jsonapi:"attr,name"`
	}
	item := &example{ID: 1, Name: "example"}

	tt := []serde.MediaType{serde.JSON, serde.YAML}
	for _, mediaType := range tt {
		t.Run(fmt.Sprintf("generates a %v request body", mediaType), func(t *testing.T) {
			r := Request("GET", "/example", WithContentType(item, mediaType))
			require.NotNil(t, r)

			result, err := serde.From[example](r.Body, mediaType)
			require.NoError(t, err)
			assert.Equal(t, 1, result.ID)
			assert.Equal(t, "example", result.Name)
		})
	}
}

func TestWithRequestBody(t *testing.T) {
	t.Run("applies the io to the request body", func(t *testing.T) {
		body := io.NopCloser(strings.NewReader("example"))
		r := Request("GET", "/example", WithRequestBody(body))

		require.NotNil(t, r)

		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Equal(t, "example", string(b))
	})
}

func TestWithContext(t *testing.T) {
	t.Run("returns a request with a new context", func(t *testing.T) {
		key := context.Key[string]("x")

		ctx := key.With(t.Context(), "example")
		r := Request("GET", "/example", WithContext(ctx))

		require.NotNil(t, r)
		assert.Equal(t, "example", key.From(r.Context()))
	})
}

func TestWithContextKeyValue(t *testing.T) {
	t.Run("returns a request with a new context", func(t *testing.T) {
		key := context.Key[string]("x")

		r := Request("GET", "/example", WithContextKeyValue(t.Context(), key, "example"))

		require.NotNil(t, r)
		assert.Equal(t, "example", key.From(r.Context()))
	})
}

func TestWithCookie(t *testing.T) {
	t.Run("adds a cookie to the request", func(t *testing.T) {
		r := Request("GET", "/example", WithCookie(&http.Cookie{Name: "example", Value: "value"}))

		require.NotNil(t, r)
		result, err := r.Cookie("example")
		require.NoError(t, err)

		assert.Equal(t, "example", result.Name)
		assert.Equal(t, "value", result.Value)
	})
}
