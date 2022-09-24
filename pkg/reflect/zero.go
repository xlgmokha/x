package reflect

func ZeroValue[T any]() T {
	var item T
	return item
}
