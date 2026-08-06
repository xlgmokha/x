package ioc

import (
	"reflect"
	"slices"
	"sync"

	"github.com/xlgmokha/x/pkg/x"
)

type Resolver[T any] func(*Container) T

type Container struct {
	registry *registry
	pending  []reflect.Type
}

type registry struct {
	mu       sync.RWMutex
	bindings map[reflect.Type]*binding
}

type binding struct {
	factory   func(*Container) any
	singleton bool
	once      sync.Once
	instance  any
}

var Default = New()

func New() *Container {
	return &Container{registry: &registry{bindings: map[reflect.Type]*binding{}}}
}

func Register[T any](c *Container, factory Resolver[T]) {
	c.bind(reflect.TypeFor[T](), &binding{factory: erase(factory)})
}

// RegisterSingleton binds T to a factory whose result is built once and reused.
//
// A cycle is reported as a *CycleError when it occurs within a single
// resolution chain. Two singletons that depend on each other and are resolved
// from different goroutines will block on each other's construction instead.
func RegisterSingleton[T any](c *Container, factory Resolver[T]) {
	c.bind(reflect.TypeFor[T](), &binding{factory: erase(factory), singleton: true})
}

func Resolve[T any](c *Container) (T, error) {
	item := reflect.TypeFor[T]()

	instance, err := c.resolve(item)
	if err != nil {
		return x.Zero[T](), err
	}
	if instance == nil {
		return x.Zero[T](), nil
	}

	value, ok := instance.(T)
	if !ok {
		return x.Zero[T](), &TypeMismatchError{Type: item, Instance: instance}
	}
	return value, nil
}

func MustResolve[T any](c *Container) T {
	return x.Must(Resolve[T](c))
}

func erase[T any](factory Resolver[T]) func(*Container) any {
	return func(c *Container) any { return factory(c) }
}

func (c *Container) bind(item reflect.Type, b *binding) {
	c.registry.mu.Lock()
	defer c.registry.mu.Unlock()

	c.registry.bindings[item] = b
}

func (c *Container) bindingFor(item reflect.Type) (*binding, bool) {
	c.registry.mu.RLock()
	defer c.registry.mu.RUnlock()

	b, ok := c.registry.bindings[item]
	return b, ok
}

func (c *Container) resolve(item reflect.Type) (any, error) {
	if slices.Contains(c.pending, item) {
		return nil, &CycleError{Type: item, Chain: c.pending}
	}

	b, ok := c.bindingFor(item)
	if !ok {
		return nil, &NotRegisteredError{Type: item}
	}

	child := c.descend(item)
	if !b.singleton {
		return b.factory(child), nil
	}

	b.once.Do(func() { b.instance = b.factory(child) })
	return b.instance, nil
}

func (c *Container) descend(item reflect.Type) *Container {
	return &Container{
		registry: c.registry,
		pending:  append(slices.Clone(c.pending), item),
	}
}
