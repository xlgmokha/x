package convert

func ToPtr[T any](item T) *T {
	return &item
}

func FromPtr[T any](p *T) T {
	return *p
}
