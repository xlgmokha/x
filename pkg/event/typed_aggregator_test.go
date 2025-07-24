package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypedAggregator(t *testing.T) {
	type announcement struct {
		message string
	}

	t.Run("Publish", func(t *testing.T) {
		t.Run("without any subscribers", func(t *testing.T) {
			aggregator := NewAggregator[announcement]()

			aggregator.Publish("announcement", announcement{
				message: "Business, Business, Business... Numbers!",
			})
		})

		t.Run("with a single subscription", func(t *testing.T) {
			called := false
			aggregator := NewAggregator[announcement]()

			aggregator.SubscribeTo("announcement", func(payload announcement) {
				called = true
				assert.Equal(t, "Hello", payload.message)
			})

			aggregator.Publish("announcement", announcement{message: "Hello"})

			assert.True(t, called)
		})

		t.Run("with multiple subscribers", func(t *testing.T) {
			aggregator := NewAggregator[announcement]()
			called := map[int]bool{}

			aggregator.SubscribeTo("announcement", func(payload announcement) {
				called[0] = true
				assert.Equal(t, "Greetings", payload.message)
			})

			aggregator.SubscribeTo("announcement", func(payload announcement) {
				called[1] = true
				assert.Equal(t, "Greetings", payload.message)
			})

			aggregator.Publish("announcement", announcement{
				message: "Greetings",
			})

			assert.Equal(t, 2, len(called))
			assert.True(t, called[0])
			assert.True(t, called[1])
		})
	})
}
