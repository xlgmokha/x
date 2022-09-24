package context

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWith(t *testing.T) {
	t.Run("injects the value into context", func(t *testing.T) {
		key := Key[int]("ticket")

		value := 42
		ctx := key.With(context.Background(), value)

		assert.Equal(t, value, ctx.Value(key))
	})

	t.Run("works like this", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "ticket", 42)
		assert.Equal(t, 42, ctx.Value("ticket"))
	})
}
