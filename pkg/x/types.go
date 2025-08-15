package x

import "reflect"

func Default[T any]() T {
	item := Zero[T]()

	if IsPtr[T](item) {
		return reflect.New(reflect.TypeOf(item).Elem()).Interface().(T)
	}

	return item
}

func Zero[T any]() T {
	var item T
	return item
}

func IsZero[T comparable](item T) bool {
	return item == Zero[T]()
}

func IsPresent[T comparable](item T) bool {
	return !IsZero[T](item)
}

func IsPtr[T any](item T) bool {
	return Is[T](item, reflect.Pointer)
}

func IsSlice[T any](item T) bool {
	return Is[T](item, reflect.Slice)
}

func Is[T any](item T, kind reflect.Kind) bool {
	return reflect.TypeOf(item).Kind() == kind
}

func ToMap[T any](item T) map[string]interface{} {
	return map[string]interface{}{
		"item": item,
	}
}
