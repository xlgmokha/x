package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregator(t *testing.T) {
	t.Run("Publish", func(t *testing.T) {
		t.Run("without any subscribers", func(t *testing.T) {
			aggregator := New()

			aggregator.Publish("announcements.engineering", "Business, Business, Business... Numbers!")
		})

		t.Run("with a single subscriber", func(t *testing.T) {
			aggregator := New()
			called := false

			aggregator.Subscribe("announcement", func(message any) {
				called = true
				assert.Equal(t, "Hello", message)
			})

			aggregator.Publish("announcement", "Hello")

			assert.True(t, called)
		})

		t.Run("with multiple subscribers", func(t *testing.T) {
			aggregator := New()
			called := map[int]bool{}

			aggregator.Subscribe("announcement", func(message any) {
				called[0] = true
				assert.Equal(t, "Greetings", message)
			})

			aggregator.Subscribe("announcement", func(message any) {
				called[1] = true
				assert.Equal(t, "Greetings", message)
			})

			aggregator.Publish("announcement", "Greetings")

			assert.Equal(t, 2, len(called))
			assert.True(t, called[0])
			assert.True(t, called[1])
		})
	})
}
