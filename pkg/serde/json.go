package serde

import (
	"encoding/json"
	"io"
)

func ToJSON[T any](w io.Writer, item T) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(item)
}

func FromJSON[T any](reader io.Reader) (T, error) {
	var item T
	return item, json.NewDecoder(reader).Decode(&item)
}
