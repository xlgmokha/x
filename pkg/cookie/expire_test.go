package cookie

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpire(t *testing.T) {
	w := httptest.NewRecorder()

	Expire(w, "example", WithDomain("example.com"))

	cookie, err := http.ParseSetCookie(w.Header().Get("Set-Cookie"))
	require.NoError(t, err)

	assert.Empty(t, cookie.Value)
	assert.Equal(t, "example", cookie.Name)
	assert.Equal(t, "example.com", cookie.Domain)
	assert.Equal(t, -1, cookie.MaxAge)
	assert.Equal(t, time.Unix(0, 0).Unix(), cookie.Expires.Unix())
}
