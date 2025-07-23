package x

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	t.Run("panics when given an error", func(t *testing.T) {
		assert.Panics(t, func() {
			Check(errors.New("Ooops..."))
		})
	})

	t.Run("does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Check(nil)
		})
	})
}

func TestMust(t *testing.T) {
	t.Run("panics when given an error", func(t *testing.T) {
		assert.Panics(t, func() {
			Must(42, errors.New("Ooops..."))
		})
	})

	t.Run("returns the value", func(t *testing.T) {
		assert.NotPanics(t, func() {
			assert.Equal(t, 42, Must(42, nil))
		})
	})
}

func TestTry(t *testing.T) {
	t.Run("does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Try(0, errors.New("Ooops..."))
		})
	})

	t.Run("returns the value", func(t *testing.T) {
		assert.NotPanics(t, func() {
			assert.Equal(t, 42, Try(42, nil))
			assert.Equal(t, "hello", Try("hello", errors.New("Error")))
		})
	})
}
