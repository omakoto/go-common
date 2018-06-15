package utils

import (
	"sync"
)

var reRegexCleaner = NewLazyRegexp(`(?:\s+|\s*#[^\n]*\n\s*)`)

// Remove whitespace and comments from a regex pattern.
func CleanUpRegexp(pattern string) string {
	return reRegexCleaner.Pattern().ReplaceAllLiteralString(pattern, "")
}

func GoMulti(n int, f func()) func() {
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			f()
			wg.Done()
		}()
	}

	return func() {
		wg.Wait()
	}
}

func DoAndEnsure(fun func(), ensure func()) {
	defer ensure()

	fun()
}
