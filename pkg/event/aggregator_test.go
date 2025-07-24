package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventAggregator(t *testing.T) {
	t.Run("Publish", func(t *testing.T) {
		t.Run("without any subscribers", func(t *testing.T) {
			events := New()

			events.Publish("announcements.engineering", "Business, Business, Business... Numbers!")
		})

		t.Run("with a single subscriber", func(t *testing.T) {
			events := New()
			called := false

			events.Subscribe("announcement", func(message any) {
				called = true
				assert.Equal(t, "Hello", message)
			})

			events.Publish("announcement", "Hello")

			assert.True(t, called)
		})

		t.Run("with multiple subscribers", func(t *testing.T) {
			events := New()
			called := map[int]bool{}

			events.Subscribe("announcement", func(message any) {
				called[0] = true
				assert.Equal(t, "Greetings", message)
			})

			events.Subscribe("announcement", func(message any) {
				called[1] = true
				assert.Equal(t, "Greetings", message)
			})

			events.Publish("announcement", "Greetings")

			assert.Equal(t, 2, len(called))
			assert.True(t, called[0])
			assert.True(t, called[1])
		})
	})
}
