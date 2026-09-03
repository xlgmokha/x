package peg

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	tt := []struct {
		stream string
		re     string

		position int
		ok       bool
		result   any
	}{
		{stream: "   abcd ", re: "", position: 0, ok: true, result: ""},
		{stream: "abcd ", re: "[a-z]", position: 1, ok: true, result: "a"},
		{stream: "   abcd ", re: "\\s*[a-z]+", position: 7, ok: true, result: "   abcd"},
		{stream: "   abcd ", re: `\s`, position: 1, ok: true, result: " "},
		{stream: "  ab", re: `\s*`, position: 2, ok: true, result: "  "},
	}
	for _, expected := range tt {
		t.Run(fmt.Sprintf("%s %d", expected.stream, expected.position), func(t *testing.T) {
			atom := Match(regexp.MustCompile(expected.re))
			ctx := NewContext(expected.stream)
			result, ok := atom(ctx)

			assert.Equal(t, expected.ok, ok)
			assert.Equal(t, expected.position, ctx.position)
			assert.Equal(t, expected.result, result)
		})
	}
}

func TestMatchNoMatch(t *testing.T) {
	tt := []struct {
		stream string
		re     string
	}{
		{stream: "123", re: "[a-z]"},
		{stream: "123", re: "xyz"},
		{stream: "", re: "[a-z]+"},
		{stream: "abc", re: "^[d-z]+$"},
	}
	for _, expected := range tt {
		t.Run(fmt.Sprintf("%s %s", expected.stream, expected.re), func(t *testing.T) {
			atom := Match(regexp.MustCompile(expected.re))
			ctx := &Context{stream: expected.stream}
			result, ok := atom(ctx)

			assert.False(t, ok)
			assert.Nil(t, result)
			assert.Equal(t, 0, ctx.position)
		})
	}
}
