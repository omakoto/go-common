package utils

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

func SortedMapFunc[K cmp.Ordered, V any](m map[K]V, compare func(K, K) int) iter.Seq2[K, V] {
	keys := slices.SortedFunc(maps.Keys(m), compare)
	return func(yield func(key K, value V) bool) {
		for _, key := range keys {
			if !yield(key, m[key]) {
				return
			}
		}
	}
}

func SortedMap[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	return SortedMapFunc(m, cmp.Compare)
}
