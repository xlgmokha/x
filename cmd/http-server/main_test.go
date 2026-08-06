package main

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xlgmokha/x/pkg/xtest"
)

func TestBuildHttpHandlerFor(t *testing.T) {
	t.Run("can be built more than once", func(t *testing.T) {
		assert.NotPanics(t, func() {
			buildHttpHandlerFor(t.TempDir())
			buildHttpHandlerFor(t.TempDir())
		})
	})

	t.Run("serves the index from its root", func(t *testing.T) {
		root := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(root, "index.html"), []byte("hello"), 0o600))

		r, w := xtest.RequestResponse("GET", "/")
		buildHttpHandlerFor(root).ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "hello", w.Body.String())
	})

	t.Run("each handler serves its own root", func(t *testing.T) {
		first, second := t.TempDir(), t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(first, "a.txt"), []byte("first"), 0o600))
		require.NoError(t, os.WriteFile(filepath.Join(second, "a.txt"), []byte("second"), 0o600))

		firstHandler := buildHttpHandlerFor(first)
		secondHandler := buildHttpHandlerFor(second)

		r, w := xtest.RequestResponse("GET", "/a.txt")
		secondHandler.ServeHTTP(w, r)
		assert.Equal(t, "second", w.Body.String())

		r, w = xtest.RequestResponse("GET", "/a.txt")
		firstHandler.ServeHTTP(w, r)
		assert.Equal(t, "first", w.Body.String())
	})
}
