package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPtr(t *testing.T) {
	t.Run("returns a pointer", func(t *testing.T) {
		var ptr *string
		value := "hello"

		ptr = ToPtr(value)

		assert.Equal(t, *ptr, value)
	})
}

func TestFromPtr(t *testing.T) {
	t.Run("returns the value", func(t *testing.T) {
		var value string

		tmp := "hello"
		ptr := &tmp

		value = FromPtr(ptr)

		assert.Equal(t, *ptr, value)
	})
}
