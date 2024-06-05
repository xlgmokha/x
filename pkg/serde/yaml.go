package serde

import (
	"io"

	"gopkg.in/yaml.v2"
)

func ToYAML[T any](w io.Writer, item T) error {
	return yaml.NewEncoder(w).Encode(item)
}
