package peg

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTag(t *testing.T) {
	t.Run("returns true", func(t *testing.T) {
		atom := Tag("filter", Sequence(
			Str("userName"),
			Space(),
			Str("eq"),
			Space(),
			Match(regexp.MustCompile(`"[a-z]+"`)),
		))
		ctx := NewContext(`userName eq "bjensen"`)
		result, ok := atom(ctx)

		assert.Equal(t, true, ok)
		assert.Equal(t, 21, ctx.position)
		assert.Equal(t, Token{"filter": []ASTNode{"userName", "eq", `"bjensen"`}}, result)
	})

	t.Run("returns false when inner parser fails", func(t *testing.T) {
		atom := Tag("filter", Str("doesnotexist"))
		ctx := &Context{stream: `userName eq "bjensen"`}
		result, ok := atom(ctx)

		assert.Equal(t, false, ok)
		assert.Nil(t, result)
		assert.Equal(t, 0, ctx.position)
	})
}
