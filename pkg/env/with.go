package env

import "os"

func With(env Vars, callback func()) {
	original := Vars{}

	for key, value := range env {
		original[key] = os.Getenv(key)
		os.Setenv(key, value)
	}

	defer func(o Vars) {
		for key, value := range o {
			os.Setenv(key, value)
		}
	}(original)

	callback()
}
