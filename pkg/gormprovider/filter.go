package gormprovider

/* OptionalValue */
func OptionalValue[T any](value T) *T {
	return &value
}
