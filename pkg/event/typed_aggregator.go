package event

import "github.com/xlgmokha/x/pkg/x"

type TypedAggregator[T any] struct {
	aggregator *Aggregator
}

func NewAggregator[T any]() *TypedAggregator[T] {
	return NewWith[T](x.New(WithoutSubscriptions()))
}

func NewWith[T any](aggregator *Aggregator) *TypedAggregator[T] {
	return x.New[*TypedAggregator[T]](WithAggregator[T](aggregator))
}

func WithAggregator[T any](aggregator *Aggregator) x.Option[*TypedAggregator[T]] {
	return x.With(func(item *TypedAggregator[T]) {
		item.aggregator = aggregator
	})
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
