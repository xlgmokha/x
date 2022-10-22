package serde

import (
	"fmt"
	"io"

	"github.com/google/jsonapi"
)

func ToPlain[T any](w io.Writer, item T) error {
	if err, ok := any(item).(*jsonapi.ErrorsPayload); ok {
		if len(err.Errors) == 1 {
			_, err := w.Write([]byte(err.Errors[0].Title))
			return err
		}
	}
	_, err := w.Write([]byte(fmt.Sprintf("%v", item)))
	return err
}
