package x

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	t.Run("returns the item", func(t *testing.T) {
		items := []int{1, 2, 3, 42, 5}

		result := Find(items, func(item int) bool {
			return item == 42
		})

		assert.Equal(t, 42, result)
	})

	t.Run("returns the zero value", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}

		result := Find(items, func(item int) bool {
			return item == 42
		})

		assert.Equal(t, 0, result)
	})
}

func TestFindAll(t *testing.T) {
	t.Run("returns all matching items", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}

		results := FindAll(items, func(item int) bool {
			return item%2 == 0
		})

		assert.Equal(t, []int{2, 4}, results)
	})

	t.Run("returns an empty slice", func(t *testing.T) {
		items := []int{1, 3, 5, 7}

		results := FindAll(items, func(item int) bool {
			return item%2 == 0
		})

		assert.Empty(t, results)
	})
}

func TestContains(t *testing.T) {
	t.Run("returns true", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}

		result := Contains(items, func(item int) bool {
			return item == 3
		})

		assert.True(t, result)
	})

	t.Run("returns false", func(t *testing.T) {
		items := []int{1, 2, 3, 4, 5}

		result := Contains(items, func(item int) bool {
			return item == 7
		})

		assert.False(t, result)
	})
}

func TestMap(t *testing.T) {
	t.Run("maps each item", func(t *testing.T) {
		items := []int{1, 2, 3}
		result := Map(items, func(item int) string {
			return strconv.Itoa(item)
		})
		assert.Equal(t, []string{"1", "2", "3"}, result)
	})
}
