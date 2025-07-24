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

func (a *Aggregator) SubscribeTo[T any](event Event, f func(T)) {
	wrapper := func(message any) {
		if typedMessage, ok := message.(T); ok {
			f(typedMessage)
		}
	}
	a.Subscribe(event, wrapper)
}

func (a *Aggregator) Publish(event Event, message any) {
	for _, subscription := range a.subscriptions[event] {
		subscription(message)
	}
}
