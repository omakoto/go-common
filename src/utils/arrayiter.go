package utils

// Iter returns a function that returns elements in a given array, or return nil at the end.
func Iter[T any](array []T) func() *T {
	i := 0

	return func() *T {
		if i >= len(array) {
			return nil
		}
		ret := &array[i]
		i++
		return ret
	}
}
