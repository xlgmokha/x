package env

import (
	"strings"
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

		t.Run("returns the default value", func(t *testing.T) {
			With(Vars{"X_VAR": ""}, func() {
				assert.Equal(t, "default", Fetch("X_VAR", "default"))
			})
		})
	})

	t.Run("Variables", func(t *testing.T) {
		for key, value := range Variables() {
			if strings.HasPrefix(key, "GITHUB_") {
				continue
			}

			assert.False(t, key == "", "key: '%v'", key)
			assert.False(t, value == "", "key: '%v', value: '%v'", key, value)
		}
	})
}
