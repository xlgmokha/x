package event

import "github.com/xlgmokha/x/pkg/x"

type TypedAggregator[T any] struct {
	aggregator *Aggregator
}

func New[T any]() *TypedAggregator[T] {
	return x.New[*TypedAggregator[T]](
		WithAggregator[T](
			x.New(
				WithDefaults(),
			),
		),
	)
}

func WithAggregator[T any](aggregator *Aggregator) x.Option[*TypedAggregator[T]] {
	return x.With(func(item *TypedAggregator[T]) {
		item.aggregator = aggregator
	})
}

func (a *TypedAggregator[T]) SubscribeTo(event Event, f func(T)) {
	a.aggregator.Subscribe(event, a.mapFrom(f))
}

func (a *TypedAggregator[T]) Publish(event Event, message T) {
	a.aggregator.Publish(event, message)
}

func (a *TypedAggregator[T]) mapFrom(f func(T)) Subscription {
	return func(message any) {
		if item, ok := message.(T); ok {
			f(item)
		}
	}
}
