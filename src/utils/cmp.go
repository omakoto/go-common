package utils

import "cmp"

func LessToCmp[E cmp.Ordered](less func(E, E) bool) func(E, E) int {
	return func(a, b E) int {
		if less(a, b) {
			return -1
		}
		if less(b, a) {
			return 1
		}
		return 0
	}
}
