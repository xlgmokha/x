package event

import "github.com/xlgmokha/x/pkg/x"

type Aggregator struct {
	subscriptions map[Event][]Subscription
}

func WithoutSubscriptions() x.Option[*Aggregator] {
	return WithSubscriptions(map[Event][]Subscription{})
}

func WithSubscriptions(subscriptions map[Event][]Subscription) x.Option[*Aggregator] {
	return x.With(func(item *Aggregator) {
		item.subscriptions = subscriptions
	})
}

func (a *Aggregator) Subscribe(event Event, f Subscription) {
	a.subscriptions[event] = append(a.subscriptions[event], f)
}

func (a *Aggregator) Publish(event Event, message any) {
	for _, subscription := range a.subscriptions[event] {
		subscription(message)
	}
}
