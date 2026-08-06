package event

import (
	"slices"
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

func (a *Aggregator) Subscribe(event Event, f Subscription) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.subscriptions[event] = append(a.subscriptions[event], f)
}

func (a *Aggregator) Publish(event Event, message any) {
	for _, subscription := range a.subscriptionsTo(event) {
		subscription(message)
	}
}

func (a *Aggregator) subscriptionsTo(event Event) []Subscription {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return slices.Clone(a.subscriptions[event])
}
