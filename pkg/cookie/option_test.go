package cookie

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {
	t.Run("WithPath", func(t *testing.T) {
		assert.Equal(t, "/blah", New("name", WithPath("/blah")).Path)
	})

	t.Run("WithHttpOnly", func(t *testing.T) {
		assert.False(t, New("x", WithHttpOnly(false)).HttpOnly)
		assert.True(t, New("x", WithHttpOnly(true)).HttpOnly)
	})

	t.Run("WithSecure", func(t *testing.T) {
		assert.False(t, New("x", WithSecure(false)).Secure)
		assert.True(t, New("x", WithSecure(true)).Secure)
	})

	t.Run("WithDomain", func(t *testing.T) {
		assert.Equal(t, "example.com", New("x", WithDomain("example.com")).Domain)
	})

	t.Run("WithSameSite", func(t *testing.T) {
		assert.Equal(t, http.SameSiteLaxMode, New("x", WithSameSite(http.SameSiteLaxMode)).SameSite)
		assert.Equal(t, http.SameSiteStrictMode, New("x", WithSameSite(http.SameSiteStrictMode)).SameSite)
		assert.Equal(t, http.SameSiteNoneMode, New("x", WithSameSite(http.SameSiteNoneMode)).SameSite)
	})

	t.Run("WithExpiration", func(t *testing.T) {
		now := time.Now()

		t.Run("with future time", func(t *testing.T) {
			expires := now.Add(1 * time.Second)
			cookie := New("x", WithExpiration(expires))
			assert.Equal(t, expires, cookie.Expires)
			assert.Equal(t, 1, cookie.MaxAge)
		})

		t.Run("with past time", func(t *testing.T) {
			expires := now.Add(-1 * time.Second)
			cookie := New("x", WithExpiration(expires))
			assert.Equal(t, expires, cookie.Expires)
			assert.Equal(t, -1, cookie.MaxAge)
		})
	})
}
