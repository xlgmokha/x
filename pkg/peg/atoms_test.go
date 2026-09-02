package peg

import (
	"fmt"
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
		t.Run(fmt.Sprintf("%s %d", expected.input, expected.position), func(t *testing.T) {
			ctx := &Context{src: expected.input}
			item, ok := atom(ctx)

			require.Nil(t, item)
			assert.Equal(t, expected.position, ctx.pos)
			assert.True(t, ok)
		})
	}
}

func TestStr(t *testing.T) {

	tt := []struct {
		input string

		position int
		ok       bool
		result   any
	}{
		{input: "   abcd ", position: 0, ok: false, result: nil},
		{input: "   abcd", position: 0, ok: false, result: nil},
		{input: " ", position: 0, ok: false, result: nil},
		{input: "", position: 0, ok: true, result: ""},
		{input: "A", position: 1, ok: true, result: "A"},
		{input: "ABCD", position: 4, ok: true, result: "ABCD"},
		{input: "abcd", position: 4, ok: true, result: "abcd"},
	}
	for _, expected := range tt {
		t.Run(fmt.Sprintf("%s %d", expected.input, expected.position), func(t *testing.T) {
			atom := Str(expected.input)
			ctx := &Context{src: expected.input}
			result, ok := atom(ctx)

			assert.Equal(t, expected.ok, ok)
			assert.Equal(t, expected.position, ctx.pos)
			assert.Equal(t, expected.result, result)
		})
	}
}
