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
		stream   string
		position int
	}{
		{stream: "   abcd ", position: 3},
		{stream: "   abcd", position: 3},
		{stream: " ", position: 1},
		{stream: "", position: 0},
		{stream: "A", position: 0},
		{stream: "ABCD", position: 0},
		{stream: "abcd", position: 0},
	}
	for _, expected := range tt {
		t.Run(fmt.Sprintf("%s %d", expected.stream, expected.position), func(t *testing.T) {
			ctx := &Context{stream: expected.stream}
			item, ok := atom(ctx)

			require.Nil(t, item)
			assert.Equal(t, expected.position, ctx.pos)
			assert.True(t, ok)
		})
	}
}

func TestStr(t *testing.T) {
	tt := []struct {
		stream string
		input  string

		position int
		ok       bool
		result   any
	}{
		{stream: "   abcd ", input: "abcd", position: 7, ok: true, result: "abcd"},
		{stream: "   abcd ", input: "ab", position: 5, ok: true, result: "ab"},
		{stream: " cd ", input: "ab", position: 0, ok: false, result: nil},
		{stream: " ABCD ", input: "abcd", position: 5, ok: true, result: "abcd"},
	}
	for _, expected := range tt {
		t.Run(fmt.Sprintf("%s %d", expected.stream, expected.position), func(t *testing.T) {
			atom := Str(expected.input)
			ctx := &Context{stream: expected.stream}
			result, ok := atom(ctx)

			assert.Equal(t, expected.ok, ok)
			assert.Equal(t, expected.position, ctx.pos)
			assert.Equal(t, expected.result, result)
		})
	}
}

func TestMatch(t *testing.T) {
	tt := []struct {
		stream string
		re     string

		position int
		ok       bool
		result   any
	}{
		{stream: "   abcd ", re: "", position: 3, ok: true, result: ""},
		{stream: "   abcd ", re: "[a-z]", position: 4, ok: true, result: "a"},
		{stream: "   abcd ", re: "[a-z]+", position: 7, ok: true, result: "abcd"},
		{stream: "   abcd ", re: `\s`, position: 0, ok: false, result: nil},
		{stream: "  ab", re: `\s*`, position: 2, ok: true, result: ""},
	}
	for _, expected := range tt {
		t.Run(fmt.Sprintf("%s %d", expected.stream, expected.position), func(t *testing.T) {
			atom := Match(expected.re)
			ctx := &Context{stream: expected.stream}
			result, ok := atom(ctx)

			assert.Equal(t, expected.ok, ok)
			assert.Equal(t, expected.position, ctx.pos)
			assert.Equal(t, expected.result, result)
		})
	}
}

func TestSequence(t *testing.T) {
	t.Run("returns true", func(t *testing.T) {
		atom := Sequence(
			Str("username"),
			Str("eq"),
			Match(`"[a-z]+"`),
		)
		ctx := &Context{stream: `username eq "bjensen"`}
		result, ok := atom(ctx)

		assert.Equal(t, true, ok)
		assert.Equal(t, 21, ctx.pos)
		assert.Equal(t, []any{"username", "eq", `"bjensen"`}, result)
	})
}
