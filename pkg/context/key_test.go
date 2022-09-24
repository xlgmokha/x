package context

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWith(t *testing.T) {
	t.Run("With", func(t *testing.T) {
		t.Run("injects the value into context", func(t *testing.T) {
			key := Key[int]("ticket")

			value := 42
			ctx := key.With(context.Background(), value)

			assert.Equal(t, value, ctx.Value(key))
		})
	})

	t.Run("From", func(t *testing.T) {
		t.Run("returns the value for the key", func(t *testing.T) {
			key := Key[time.Time]("secret")
			now := time.Now()

			ctx := key.With(context.Background(), now)

			assert.Equal(t, now, key.From(ctx))
		})

		t.Run("returns the zero value", func(t *testing.T) {
			key := Key[int]("not-found")

			assert.Equal(t, 0, key.From(context.Background()))
		})
	})
}
