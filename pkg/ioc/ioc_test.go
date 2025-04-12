package ioc

import (
	"testing"
	"time"

	"github.com/golobby/container/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testItem struct {
	num int64
}

func TestIOC(t *testing.T) {
	c := container.New()

	t.Run("Register", func(t *testing.T) {
		var ctr int64
		err := Register[*testItem](c, func() *testItem {
			item := &testItem{num: ctr}
			ctr++
			return item
		})

		require.NoError(t, err)

		assert.Nil(t, c.Call(func(result *testItem) {
			assert.Equal(t, int64(1), result.num)
		}))

		assert.Nil(t, c.Call(func(result *testItem) {
			assert.Equal(t, int64(2), result.num)
		}))
	})

	t.Run("RegisterSingleton", func(t *testing.T) {
		err := RegisterSingleton[*testItem](c, func() *testItem {
			return &testItem{num: time.Now().Unix()}
		})
		require.NoError(t, err)

		var first int64
		assert.Nil(t, c.Call(func(result *testItem) {
			first = result.num
		}))

		assert.Nil(t, c.Call(func(result *testItem) {
			assert.Equal(t, first, result.num)
		}))
	})

	t.Run("Resolve", func(t *testing.T) {
		item := &testItem{}

		require.NoError(t, Register[*testItem](c, func() *testItem { return item }))

		result, err := Resolve[*testItem](c)
		require.NoError(t, err)

		assert.Equal(t, item, result)
	})

	t.Run("MustResolve", func(t *testing.T) {
		item := &testItem{}

		require.NoError(t, Register[*testItem](c, func() *testItem { return item }))

		result := MustResolve[*testItem](c)
		assert.Equal(t, item, result)
	})
}
