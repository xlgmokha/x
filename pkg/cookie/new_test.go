package cookie

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cookie := New("name", WithValue("value"))

	require.NotNil(t, cookie)
	assert.Equal(t, "name", cookie.Name)
	assert.Equal(t, "value", cookie.Value)
}
