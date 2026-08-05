package mapper

import (
	"reflect"
	"sync"

	"github.com/xlgmokha/x/pkg/x"
)

type Mapping[TInput any, TOutput any] func(TInput) TOutput

type key struct {
	input  reflect.Type
	output reflect.Type
}

var (
	mu       sync.RWMutex
	mappings map[key]interface{}
)

func init() {
	mappings = map[key]interface{}{}
}

func Register[Input any, Output any](mapping Mapping[Input, Output]) {
	mu.Lock()
	defer mu.Unlock()

	mappings[keyFor[Input, Output]()] = mapping
}

func MapFrom[Input any, Output any](input Input) Output {
	if mapping, ok := mappingFor[Input, Output](); ok {
		return mapping(input)
	}
	return x.Zero[Output]()
}

func mappingFor[Input any, Output any]() (Mapping[Input, Output], bool) {
	mu.RLock()
	defer mu.RUnlock()

	mapping, ok := mappings[keyFor[Input, Output]()]
	if !ok {
		return nil, false
	}
	return mapping.(Mapping[Input, Output]), true
}

func MapEachFrom[Input any, Output any](input []Input) []Output {
	results := []Output{}

	mapping, ok := mappingFor[Input, Output]()
	if !ok {
		return results
	}

	for _, item := range input {
		results = append(results, mapping(item))
	}
	return results
}

func keyFor[Input any, Output any]() key {
	return key{input: reflect.TypeFor[Input](), output: reflect.TypeFor[Output]()}
}
