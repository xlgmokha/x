package event

type Event any
type Subscription func(any)

type Aggregator struct {
	subscriptions map[Event][]Subscription
}

func New() *Aggregator {
	return &Aggregator{
		subscriptions: map[Event][]Subscription{},
	}
}

func (a *Aggregator) Subscribe(event Event, f Subscription) {
	a.subscriptions[event] = append(a.subscriptions[event], f)
}

func (a *Aggregator) Publish(event Event, message any) {
	for _, subscription := range a.subscriptions[event] {
		subscription(message)
	}
}
