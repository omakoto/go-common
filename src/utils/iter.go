package utils

import "sync"

// Iterator is a versatile iterator.
type Iterator[T any] struct {
	nextFetcher func() (*T, bool)
	cleaner     func()
	cleaned     bool

	lock sync.Mutex
}

// NewIterable creates a new Iterable, with an optional cleaner, which is called at the end.
func NewIterable[T any](nextFetcher func() (*T, bool), cleaner func()) *Iterator[T] {
	return &Iterator[T]{
		nextFetcher: nextFetcher,
		cleaner:     cleaner,
	}
}

// Next returns a next element.
func (i *Iterator[T]) Next() (element *T, ok bool) {
	i.lock.Lock()
	defer i.lock.Unlock()

	ret, ok := i.nextFetcher()

	if !ok && !i.cleaned {
		i.cleanUp()
	}

	return ret, ok
}

func (i *Iterator[T]) Close() {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.cleanUp()
}

func (i *Iterator[T]) cleanUp() {
	if i.cleaned {
		return
	}
	i.cleaned = true

	if i.cleaner != nil {
		i.cleaner()
	}
}
