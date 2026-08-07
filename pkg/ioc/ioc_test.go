package ioc

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testItem struct {
	num int64
}

type Greeter interface {
	Greet() string
}

type englishGreeter struct{}

func (englishGreeter) Greet() string { return "hello" }

type frenchGreeter struct{}

func (frenchGreeter) Greet() string { return "bonjour" }

func TestRegister(t *testing.T) {
	t.Run("resolves a new instance each time", func(t *testing.T) {
		c := New()
		var ctr int64

		c.Register(func(*Container) *testItem {
			ctr++
			return &testItem{num: ctr}
		})

		first, err := c.Resolve[*testItem]()
		require.NoError(t, err)
		second, err := c.Resolve[*testItem]()
		require.NoError(t, err)

		assert.Equal(t, int64(1), first.num)
		assert.Equal(t, int64(2), second.num)
		assert.NotSame(t, first, second)
	})

	t.Run("binds a concrete type to an interface", func(t *testing.T) {
		c := New()

		c.Register(func(*Container) Greeter { return englishGreeter{} })

		greeter, err := c.Resolve[Greeter]()
		require.NoError(t, err)
		assert.Equal(t, "hello", greeter.Greet())
	})

	t.Run("rebinding an interface replaces the implementation", func(t *testing.T) {
		c := New()

		c.Register(func(*Container) Greeter { return englishGreeter{} })
		c.Register(func(*Container) Greeter { return frenchGreeter{} })

		greeter, err := c.Resolve[Greeter]()
		require.NoError(t, err)
		assert.Equal(t, "bonjour", greeter.Greet())
	})

	t.Run("keeps interface and concrete bindings distinct", func(t *testing.T) {
		c := New()

		c.Register(func(*Container) Greeter { return englishGreeter{} })
		c.Register(func(*Container) englishGreeter { return englishGreeter{} })

		greeter, err := c.Resolve[Greeter]()
		require.NoError(t, err)
		assert.Equal(t, "hello", greeter.Greet())

		concrete, err := c.Resolve[englishGreeter]()
		require.NoError(t, err)
		assert.Equal(t, "hello", concrete.Greet())
	})

	t.Run("resolves a dependency through the container", func(t *testing.T) {
		c := New()

		c.Register(func(*Container) Greeter { return frenchGreeter{} })
		c.Register(func(c *Container) *testItem {
			greeter := c.MustResolve[Greeter]()
			return &testItem{num: int64(len(greeter.Greet()))}
		})

		item, err := c.Resolve[*testItem]()
		require.NoError(t, err)
		assert.Equal(t, int64(len("bonjour")), item.num)
	})
}

func TestRegisterSingleton(t *testing.T) {
	t.Run("resolves the same instance every time", func(t *testing.T) {
		c := New()
		var ctr int64

		c.RegisterSingleton(func(*Container) *testItem {
			ctr++
			return &testItem{num: ctr}
		})

		first, err := c.Resolve[*testItem]()
		require.NoError(t, err)
		second, err := c.Resolve[*testItem]()
		require.NoError(t, err)

		assert.Same(t, first, second)
		assert.Equal(t, int64(1), ctr)
	})

	t.Run("builds the instance only once under concurrency", func(t *testing.T) {
		c := New()
		var mu sync.Mutex
		calls := 0

		c.RegisterSingleton(func(*Container) Greeter {
			mu.Lock()
			defer mu.Unlock()
			calls++
			return englishGreeter{}
		})

		var wg sync.WaitGroup
		for range 8 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = c.Resolve[Greeter]()
			}()
		}
		wg.Wait()

		assert.Equal(t, 1, calls)
	})
}

func TestResolveUnregistered(t *testing.T) {
	t.Run("returns a typed error", func(t *testing.T) {
		c := New()

		result, err := c.Resolve[Greeter]()

		require.Error(t, err)
		assert.Nil(t, result)

		notRegistered, ok := errors.AsType[*NotRegisteredError](err)
		require.True(t, ok)
		assert.Equal(t, reflect.TypeFor[Greeter](), notRegistered.Type)
		assert.Contains(t, err.Error(), "Greeter")
	})

	t.Run("MustResolve panics", func(t *testing.T) {
		c := New()

		assert.Panics(t, func() { c.MustResolve[Greeter]() })
	})

	t.Run("MustResolve returns the instance when registered", func(t *testing.T) {
		c := New()
		c.Register(func(*Container) Greeter { return englishGreeter{} })

		assert.Equal(t, "hello", c.MustResolve[Greeter]().Greet())
	})
}

func TestResolveTypeMismatch(t *testing.T) {
	t.Run("reports a binding that produces the wrong type", func(t *testing.T) {
		c := New()
		c.bind(reflect.TypeFor[Greeter](), &binding{factory: func(*Container) any { return 42 }})

		result, err := c.Resolve[Greeter]()

		require.Error(t, err)
		assert.Nil(t, result)

		mismatch, ok := errors.AsType[*TypeMismatchError](err)
		require.True(t, ok)
		assert.Equal(t, reflect.TypeFor[Greeter](), mismatch.Type)
		assert.Contains(t, err.Error(), "produced int")
	})

	t.Run("a factory that returns nil is not a mismatch", func(t *testing.T) {
		c := New()
		c.Register(func(*Container) Greeter { return nil })

		result, err := c.Resolve[Greeter]()

		require.NoError(t, err)
		assert.Nil(t, result)
	})
}

func TestResolveCycle(t *testing.T) {
	t.Run("a singleton that resolves itself returns an error", func(t *testing.T) {
		c := New()
		c.RegisterSingleton(func(c *Container) Greeter {
			_, err := c.Resolve[Greeter]()
			require.Error(t, err)
			_, ok := errors.AsType[*CycleError](err)
			assert.True(t, ok, "want a CycleError, got %v", err)
			return englishGreeter{}
		})

		done := make(chan error, 1)
		go func() { _, err := c.Resolve[Greeter](); done <- err }()

		select {
		case err := <-done:
			require.NoError(t, err)
		case <-time.After(2 * time.Second):
			t.Fatal("Resolve deadlocked on a self-referential singleton")
		}
	})

	t.Run("a two-step cycle returns an error", func(t *testing.T) {
		c := New()
		c.Register(func(c *Container) Greeter {
			_, _ = c.Resolve[*testItem]()
			return englishGreeter{}
		})
		c.Register(func(c *Container) *testItem {
			_, err := c.Resolve[Greeter]()
			require.Error(t, err)
			_, ok := errors.AsType[*CycleError](err)
			assert.True(t, ok, "want a CycleError, got %v", err)
			return &testItem{}
		})

		_, err := c.Resolve[Greeter]()
		require.NoError(t, err)
	})
}
