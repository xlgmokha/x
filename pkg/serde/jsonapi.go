package serde

import (
	"io"
	"reflect"

	"github.com/google/jsonapi"
	"github.com/xlgmokha/x/pkg/x"
)

func ToJSONAPI[T any](w io.Writer, item T) error {
	if err, ok := any(item).(*jsonapi.ErrorsPayload); ok {
		return jsonapi.MarshalErrors(w, err.Errors)
	}

	if x.IsSlice[T](item) {
		firstItem := reflect.TypeOf(item).Elem()
		if firstItem.Kind() == reflect.Pointer {
			return jsonapi.MarshalPayload(w, item)
		}

		sliceValue := reflect.ValueOf(item)
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.PointerTo(firstItem)), 0, sliceValue.Len())
		for i := 0; i < sliceValue.Len(); i++ {
			slice = reflect.Append(slice, sliceValue.Index(i).Addr())
		}
		return jsonapi.MarshalPayload(w, slice.Interface())
	}

	if x.IsPtr(item) {
		return jsonapi.MarshalPayload(w, item)
	}
	return jsonapi.MarshalPayload(w, &item)
}

func FromJSONAPI[T any](reader io.Reader) (T, error) {
	item := x.Default[T]()
	if _, ok := any(item).(*jsonapi.ErrorsPayload); ok {
		return FromJSON[T](reader)
	}
	if x.IsSlice[T](item) {
		sliceType := reflect.TypeOf(item).Elem()

		items, err := jsonapi.UnmarshalManyPayload(reader, sliceType)
		if err != nil {
			return item, err
		}
		slice := reflect.MakeSlice(reflect.SliceOf(sliceType), 0, len(items))
		for _, item := range items {
			slice = reflect.Append(slice, reflect.ValueOf(item))
		}
		return slice.Interface().(T), err
	}

	if x.IsPtr(item) {
		return item, jsonapi.UnmarshalPayload(reader, item)
	}
	return item, jsonapi.UnmarshalPayload(reader, &item)
}
