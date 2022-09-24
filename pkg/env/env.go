package env

import "os"

func Fetch(key string, defaultValue string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return defaultValue
}
