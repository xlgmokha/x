package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xlgmokha/x/pkg/x"
)

func TestAggregator(t *testing.T) {
	t.Run("Publish", func(t *testing.T) {
		t.Run("without any subscribers", func(t *testing.T) {
			aggregator := x.New(WithDefaults())

			aggregator.Publish("announcements.engineering", "Business, Business, Business... Numbers!")
		})

		t.Run("with a single subscriber", func(t *testing.T) {
			aggregator := x.New(WithDefaults())
			called := false

			aggregator.Subscribe("announcement", func(message any) {
				called = true
				assert.Equal(t, "Hello", message)
			})

			aggregator.Publish("announcement", "Hello")

			assert.True(t, called)
		})

		t.Run("with multiple subscribers", func(t *testing.T) {
			aggregator := x.New(WithDefaults())
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

func TestPublishDoesNotDeadlock(t *testing.T) {
	t.Run("when a subscriber subscribes during publish", func(t *testing.T) {
		aggregator := x.New(WithDefaults())

		aggregator.Subscribe("announcement", func(message any) {
			aggregator.Subscribe("other", func(any) {})
		})

		done := make(chan struct{})
		go func() {
			defer close(done)
			aggregator.Publish("announcement", "Hello")
		}()

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("Publish deadlocked when a subscriber subscribed")
		}
	})
}
