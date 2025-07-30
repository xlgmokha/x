package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type person struct {
	name string
	age  int
}

func withName(name string) Option[*person] {
	return With(func(person *person) {
		person.name = name
	})
}

func withAge(age int) Option[*person] {
	return With(func(person *person) {
		person.age = age
	})
}

func TestOption(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		t.Run("without options", func(t *testing.T) {
			item := New[*person]()

			require.NotNil(t, item)
			assert.Equal(t, "", item.name)
			assert.Zero(t, item.age)
		})

		t.Run("with name", func(t *testing.T) {
			item := New[*person](withName("mo"))

			require.NotNil(t, item)
			assert.Equal(t, "mo", item.name)
		})

		t.Run("with age", func(t *testing.T) {
			item := New[*person](withAge(42))

			require.NotNil(t, item)
			assert.Equal(t, 42, item.age)
		})
	})

	t.Run("NewWith", func(t *testing.T) {
		t.Run("without options", func(t *testing.T) {
			p := &person{}
			item := NewWith[*person](p)

			require.NotNil(t, item)
			require.Same(t, p, item)
			assert.Equal(t, "", item.name)
			assert.Zero(t, item.age)
		})

		t.Run("with name", func(t *testing.T) {
			p := &person{}
			item := NewWith[*person](p, withName("mo"))

			require.NotNil(t, item)
			require.Same(t, p, item)
			assert.Equal(t, "mo", item.name)
		})

		t.Run("with age", func(t *testing.T) {
			p := &person{}
			item := NewWith[*person](p, withAge(42))

			require.NotNil(t, item)
			require.Same(t, p, item)
			assert.Equal(t, 42, item.age)
		})
	})
}
