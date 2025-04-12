package ioc

import (
	"github.com/golobby/container/v3"
	"github.com/xlgmokha/x/pkg/x"
)

var Default container.Container = container.Global

type Resolver[T any] func() T

func Register[T any](c container.Container, factory Resolver[T]) error {
	return c.Transient(func() T {
		return factory()
	})
}

func RegisterSingleton[T any](c container.Container, factory Resolver[T]) error {
	return c.Singleton(func() T {
		return factory()
	})
}

func Resolve[T any](c container.Container) (T, error) {
	var item T
	err := c.Call(func(i T) {
		item = i
	})
	return item, err
}

func MustResolve[T any](c container.Container) T {
	return x.Must(Resolve[T](c))
}
