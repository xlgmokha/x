package peg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpace(t *testing.T) {
	atom := Space()

	tt := []struct {
		input    string
		position int
	}{
		{input: "   abcd ", position: 3},
		{input: "   abcd", position: 3},
		{input: " ", position: 1},
		{input: "", position: 0},
		{input: "A", position: 0},
		{input: "ABCD", position: 0},
		{input: "abcd", position: 0},
	}
	for _, expected := range tt {
		t.Run("returns true", func(t *testing.T) {
			ctx := &Context{src: expected.input}
			item, ok := atom(ctx)

			require.Nil(t, item)
			assert.True(t, ok)
			assert.Equal(t, expected.position, ctx.pos)
		})
	}
}
