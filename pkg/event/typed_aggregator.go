package event

type TypedAggregator[T any] struct {
	aggregator *Aggregator
}

func NewAggregator[T any]() *TypedAggregator[T] {
	return NewWith[T](New())
}

func NewWith[T any](aggregator *Aggregator) *TypedAggregator[T] {
	return &TypedAggregator[T]{
		aggregator: aggregator,
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
