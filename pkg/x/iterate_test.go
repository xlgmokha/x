package x

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	t.Run("returns nil", func(t *testing.T) {
		items := []*http.Client{http.DefaultClient}

		result := Find(items, func(item *http.Client) bool {
			return false
		})

		assert.Nil(t, result)
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

	t.Run("returns false on nil", func(t *testing.T) {
		items := []*http.Client{http.DefaultClient}

		result := Contains(items, func(item *http.Client) bool {
			return false
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

func TestEach(t *testing.T) {
	t.Run("visits each item", func(t *testing.T) {
		results := []int{}

		Each([]int{1, 2, 3}, func(i int) {
			results = append(results, i)
		})

		assert.Equal(t, []int{1, 2, 3}, results)
	})
}

func TestInject(t *testing.T) {
	t.Run("collects each item", func(t *testing.T) {
		items := []int{1, 3, 6}

		result := Inject[int, map[int]bool](items,
			map[int]bool{},
			func(memo map[int]bool, item int) {
				memo[item] = true
			})

		require.NotNil(t, result)
		assert.True(t, result[1])
		assert.True(t, result[3])
		assert.True(t, result[6])
	})
}
