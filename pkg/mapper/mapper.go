package mapper

import (
	"fmt"
	"reflect"
)

type Mapping[TInput any, TOutput any] func(TInput) TOutput

var mappings map[string]interface{}

func init() {
	mappings = map[string]interface{}{}
}

func Register[Input any, Output any](mapping Mapping[Input, Output]) {
	mappings[keyFor[Input, Output]()] = mapping
}

func MapFrom[Input any, Output any](input Input) Output {
	if mapping, ok := mappings[keyFor[Input, Output]()]; ok {
		return mapping.(Mapping[Input, Output])(input)
	}
	var output Output
	return output
}

func MapEachFrom[Input any, Output any](input []Input) []Output {
	var zero Output
	zeroValue := reflect.Zero(reflect.TypeOf(zero))

	results := []Output{}
	for _, item := range input {
		tmp := MapFrom[Input, Output](item)
		if zeroValue != reflect.ValueOf(tmp) {
			results = append(results, tmp)
		}
	}
	return results
}

func keyFor[Input any, Output any]() string {
	var input Input
	var output Output
	return fmt.Sprintf("%v-%v", reflect.TypeOf(input), reflect.TypeOf(output))
}
