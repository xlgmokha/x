package event

type Event string
type Subscription func(any)

type Aggregator struct {
	subscriptions map[Event]Subscription
}

func New() *Aggregator {
	return &Aggregator{
		subscriptions: map[Event]Subscription{},
	}
}

func (a *Aggregator) Subscribe(event Event, f Subscription) {
	a.subscriptions[event] = f
}

func (a *Aggregator) Publish(event Event, message any) {
	subscription := a.subscriptions[event]
	subscription(message)
}
