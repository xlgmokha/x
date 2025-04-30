package x

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {
	t.Run("Zero", func(t *testing.T) {
		t.Run("returns nil", func(t *testing.T) {
			assert.Nil(t, Zero[[]string]())
			assert.Nil(t, Zero[*http.Client]())
		})

		t.Run("returns 0", func(t *testing.T) {
			assert.Zero(t, Zero[int]())
		})

		t.Run("returns ''", func(t *testing.T) {
			assert.Equal(t, "", Zero[string]())
		})

		t.Run("returns false", func(t *testing.T) {
			assert.False(t, Zero[bool]())
		})
	})

	t.Run("Default", func(t *testing.T) {
		t.Run("returns a new instance", func(t *testing.T) {
			result := Default[*http.Client]()

			assert.NotNil(t, result)
			assert.True(t, IsPtr(result))
		})

		t.Run("returns 0", func(t *testing.T) {
			assert.Equal(t, 0, Default[int]())
		})

		t.Run("returns ''", func(t *testing.T) {
			assert.Equal(t, "", Default[string]())
		})

		t.Run("returns false", func(t *testing.T) {
			assert.False(t, Default[bool]())
		})

		t.Run("returns an empty slice", func(t *testing.T) {
			assert.Equal(t, 0, len(Default[[]string]()))
		})
	})

	t.Run("New", func(t *testing.T) {
		t.Run("returns a new instance", func(t *testing.T) {
			result := New[*http.Client]()

			assert.NotNil(t, result)
			assert.True(t, IsPtr(result))
		})

		t.Run("configures the new instance", func(t *testing.T) {
			option := func(r *http.Request) *http.Request {
				r.Method = "GET"
				return r
			}

			result := New[*http.Request](option)

			assert.NotNil(t, result)
			assert.True(t, IsPtr(result))
			assert.Equal(t, "GET", result.Method)
		})

		t.Run("configures a new instance with multiple options", func(t *testing.T) {
			option := func(r *http.Request) *http.Request {
				r.Method = "GET"
				return r
			}
			otherOption := func(r *http.Request) *http.Request {
				r.URL = Must(url.Parse("/example"))
				return r
			}

			result := New[*http.Request](option, otherOption)

			assert.NotNil(t, result)
			assert.True(t, IsPtr(result))
			assert.Equal(t, "GET", result.Method)
			assert.Equal(t, "/example", result.URL.Path)
		})
	})

	t.Run("IsSlice", func(t *testing.T) {
		t.Run("returns true", func(t *testing.T) {
			assert.True(t, IsSlice(Default[[]string]()))
		})

		t.Run("returns false", func(t *testing.T) {
			assert.False(t, IsSlice(""))
			assert.False(t, IsSlice[string](""))
		})
	})

	t.Run("IsPtr", func(t *testing.T) {
		t.Run("returns true", func(t *testing.T) {
			assert.True(t, IsPtr[*http.Client](&http.Client{}))
			assert.True(t, IsPtr[*http.Client](nil))
		})

		t.Run("returns false", func(t *testing.T) {
			assert.False(t, IsPtr[http.Client](http.Client{}))
			assert.False(t, IsPtr[int](0))
			assert.False(t, IsPtr[string](""))
			assert.False(t, IsPtr[bool](true))
			assert.False(t, IsPtr[[]string]([]string{}))
		})
	})

	t.Run("IsZero", func(t *testing.T) {
		t.Run("returns true", func(t *testing.T) {
			assert.True(t, IsZero[int](0))
			assert.True(t, IsZero[bool](false))
			assert.True(t, IsZero[string](""))
			assert.True(t, IsZero[*http.Client](nil))
		})

		t.Run("returns false", func(t *testing.T) {
			assert.False(t, IsZero[int](1))
			assert.False(t, IsZero[bool](true))
			assert.False(t, IsZero[string]("hello"))
			assert.False(t, IsZero[*http.Client](&http.Client{}))
		})
	})

	t.Run("IsPresent", func(t *testing.T) {
		t.Run("returns false", func(t *testing.T) {
			assert.False(t, IsPresent[int](0))
			assert.False(t, IsPresent[bool](false))
			assert.False(t, IsPresent[string](""))
			assert.False(t, IsPresent[*http.Client](nil))
		})

		t.Run("returns true", func(t *testing.T) {
			assert.True(t, IsPresent[int](1))
			assert.True(t, IsPresent[bool](true))
			assert.True(t, IsPresent[string]("hello"))
			assert.True(t, IsPresent[*http.Client](&http.Client{}))
		})
	})
}
