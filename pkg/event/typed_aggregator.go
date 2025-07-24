package event

type TypedAggregator[T any] struct {
	aggregator *Aggregator
}

func NewTypedAggregator[T any]() *TypedAggregator[T] {
	return &TypedAggregator[T]{
		aggregator: New(),
	}
}

func (a *TypedAggregator[T]) SubscribeTo(event Event, f func(T)) {
	a.aggregator.Subscribe(event, func(message any) {
		if item, ok := message.(T); ok {
			f(item)
		}
	})

}

func (a *TypedAggregator[T]) Publish(event Event, message T) {
	a.aggregator.Publish(event, message)
}
