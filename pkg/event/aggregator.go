package event

import (
	"sync"

	"github.com/xlgmokha/x/pkg/x"
)

type Aggregator struct {
	mu            sync.RWMutex
	subscriptions map[Event][]Subscription
}

func WithDefaults() x.Option[*Aggregator] {
	return x.With(func(item *Aggregator) {
		item.mu = sync.RWMutex{}
		item.subscriptions = map[Event][]Subscription{}
	})
}

func WithSubscriptions(subscriptions map[Event][]Subscription) x.Option[*Aggregator] {
	return x.With(func(item *Aggregator) {
		item.subscriptions = subscriptions
	})
}

func (a *Aggregator) Subscribe(event Event, f Subscription) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.subscriptions[event] = append(a.subscriptions[event], f)
}

func (a *Aggregator) Publish(event Event, message any) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, subscription := range a.subscriptions[event] {
		subscription(message)
	}
}
