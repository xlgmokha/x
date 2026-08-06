package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnv(t *testing.T) {
	t.Run("Fetch", func(t *testing.T) {
		t.Run("returns the environment variable value", func(t *testing.T) {
			With(Vars{"SECRET": "42"}, func() {
				assert.Equal(t, "42", Fetch("SECRET", "default"))
			})
		})

		t.Run("returns an explicitly empty value", func(t *testing.T) {
			With(Vars{"X_VAR": ""}, func() {
				assert.Equal(t, "", Fetch("X_VAR", "default"))
			})
		})

		t.Run("returns the default value when unset", func(t *testing.T) {
			assert.Equal(t, "default", Fetch("X_DEFINITELY_UNSET_VAR", "default"))
		})
	})

	t.Run("Variables", func(t *testing.T) {
		for key, value := range Variables() {
			assert.NotEmpty(t, key)
			assert.NotNil(t, value)
		}

		t.Run("does not panic on an entry without an equals sign", func(t *testing.T) {
			assert.NotPanics(t, func() {
				_ = Variables()
			})
		})
	})
}
