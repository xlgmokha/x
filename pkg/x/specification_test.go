package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncSpecification(t *testing.T) {
	even := Predicate[int](func(item int) bool {
		return item%2 == 0
	})

	odd := Predicate[int](func(item int) bool {
		return !even.IsSatisfiedBy(item)
	})

	greaterThan := func(max int) Specification[int] {
		return Predicate[int](func(i int) bool {
			return i > max
		})
	}

	t.Run("IsSatisfiedBy", func(t *testing.T) {
		assert.True(t, even.IsSatisfiedBy(2))
		assert.False(t, even.IsSatisfiedBy(1))

		assert.True(t, odd.IsSatisfiedBy(1))
		assert.False(t, odd.IsSatisfiedBy(2))
	})

	t.Run("Or", func(t *testing.T) {
		assert.True(t, even.Or(greaterThan(3)).IsSatisfiedBy(2))
		assert.True(t, even.Or(greaterThan(3)).IsSatisfiedBy(5))
		assert.False(t, even.Or(greaterThan(3)).IsSatisfiedBy(3))
		assert.False(t, even.Or(greaterThan(3)).IsSatisfiedBy(1))
	})

	t.Run("And", func(t *testing.T) {
		assert.True(t, even.And(greaterThan(10)).IsSatisfiedBy(12))
		assert.False(t, even.And(greaterThan(10)).IsSatisfiedBy(8))
	})

	p := Predicate[int](func(item int) bool {
		return item == 1
	}).Or(Predicate[int](func(item int) bool {
		return item == 2
	})).Or(Predicate[int](func(item int) bool {
		return item == 3
	}))

	assert.False(t, p.IsSatisfiedBy(0))
	assert.True(t, p.IsSatisfiedBy(1))
	assert.True(t, p.IsSatisfiedBy(2))
	assert.True(t, p.IsSatisfiedBy(3))
	assert.False(t, p.IsSatisfiedBy(4))
}
