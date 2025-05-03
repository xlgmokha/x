package cookie

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReset(t *testing.T) {
	result := Reset(
		"example",
		WithSecure(true),
		WithHttpOnly(true),
		WithSameSite(http.SameSiteDefaultMode),
		WithDomain("example.com"),
	)

	assert.Equal(t, -1, result.MaxAge)
	assert.Equal(t, time.Unix(0, 0), result.Expires)
	assert.Empty(t, result.Value)
	assert.Equal(t, time.Unix(0, 0), result.Expires)
	assert.True(t, result.HttpOnly)
	assert.True(t, result.Secure)
	assert.Equal(t, http.SameSiteDefaultMode, result.SameSite)
	assert.Equal(t, "example.com", result.Domain)
}
