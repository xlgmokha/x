package x

var Panic = func(err error) {
	panic(err)
}

func Check(err error) {
	if err != nil {
		Panic(err)
	}
}

func Must[T any](item T, err error) T {
	Check(err)
	return item
}
