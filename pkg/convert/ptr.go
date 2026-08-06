package convert

import "github.com/xlgmokha/x/pkg/x"

func ToPtr[T any](item T) *T {
	return &item
}

func FromPtr[T any](p *T) T {
	if p == nil {
		return x.Zero[T]()
	}
	return *p
}
