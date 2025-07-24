package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventAggregator(t *testing.T) {
	t.Run("Publish", func(t *testing.T) {
		events := New()
		called := false

		events.Subscribe("announcement", func(message any) {
			called = true
			assert.Equal(t, "Hello", message)
		})

		events.Publish("announcement", "Hello")

		assert.True(t, called)
	})
}
