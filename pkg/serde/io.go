package serde

import (
	"io"
)

func From[T any](r io.Reader, mediaType MediaType) (T, error) {
	switch mediaType {
	case JSONAPI:
		return FromJSONAPI[T](r)
	case YAML:
		return FromYAML[T](r)
	default:
		return FromJSON[T](r)
	}
}

func To[T any](w io.Writer, item T, mediaType MediaType) error {
	switch mediaType {
	case JSONAPI:
		return ToJSONAPI[T](w, item)
	case Text:
		return ToPlain[T](w, item)
	case YAML:
		return ToYAML[T](w, item)
	default:
		return ToJSON[T](w, item)
	}
}
