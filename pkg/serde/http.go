package serde

import (
	"net/http"
)

func FromHTTP[T any](r *http.Request) (T, error) {
	return From[T](r.Body, MediaTypeFor(r.Header.Get("Content-Type")))
}

func ToHTTP[T any](w http.ResponseWriter, r *http.Request, item T) error {
	mediaType := MediaTypeFor(r.Header.Get("Accept"))
	w.Header().Set("Content-Type", string(mediaType))
	return To[T](w, item, mediaType)
}
