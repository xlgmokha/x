package env

import (
	"os"
)

func With(env Vars, callback func()) {
	original := Vars{}

	for key, value := range env {
		if val, ok := os.LookupEnv(key); ok {
			original[key] = val
		}
		os.Setenv(key, value)
	}

	defer func() {
		for key, _ := range env {
			os.Unsetenv(key)
		}
		for key, value := range original {
			os.Setenv(key, value)
		}
	}()

	callback()
}
