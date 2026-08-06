package env

import (
	"os"
	"strings"
)

func Fetch(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func Variables() Vars {
	items := Vars{}
	for _, line := range os.Environ() {
		if key, value, ok := strings.Cut(line, "="); ok {
			items[key] = value
		}
	}
	return items
}
