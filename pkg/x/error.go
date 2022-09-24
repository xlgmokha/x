package x

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Must[T any](item T, err error) T {
	Check(err)
	return item
}
