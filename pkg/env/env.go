package env

import (
	"os"
	"strings"
)

func Fetch(key string, defaultValue string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return defaultValue
}

func Variables() Vars {
	items := Vars{}
	for _, line := range os.Environ() {
		segments := strings.SplitN(line, "=", 2)
		items[segments[0]] = segments[1]
	}
	return items
}
