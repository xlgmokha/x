package ioc

import (
	"fmt"
	"reflect"
	"strings"
)

type NotRegisteredError struct {
	Type reflect.Type
}

func (e *NotRegisteredError) Error() string {
	return fmt.Sprintf("ioc: no binding registered for %v", e.Type)
}

type TypeMismatchError struct {
	Type     reflect.Type
	Instance any
}

func (e *TypeMismatchError) Error() string {
	return fmt.Sprintf("ioc: binding for %v produced %T", e.Type, e.Instance)
}

type CycleError struct {
	Type  reflect.Type
	Chain []reflect.Type
}

func (e *CycleError) Error() string {
	names := make([]string, 0, len(e.Chain)+1)
	for _, item := range e.Chain {
		names = append(names, fmt.Sprintf("%v", item))
	}
	names = append(names, fmt.Sprintf("%v", e.Type))

	return fmt.Sprintf("ioc: cycle resolving %v: %v", e.Type, strings.Join(names, " -> "))
}
