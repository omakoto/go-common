package utils

import (
	"runtime"
	"sync"
)

var reRegexCleaner = NewLazyRegexp(`(?:\s+|\s*#[^\n]*\n\s*)`)

// Remove whitespace and comments from a regex pattern.
func CleanUpRegexp(pattern string) string {
	return reRegexCleaner.Pattern().ReplaceAllLiteralString(pattern, "")
}

func FirstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// StringSlice is a convenient way to build a string slice.
func StringSlice(arr ...string) []string {
	return arr
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

func SourceLineNo() int {
	_, _, line, ok := runtime.Caller(1)

	if ok {
		return line
	}
	return 0
}

func SourceFileName() string {
	_, source, _, ok := runtime.Caller(1)

	if ok {
		return source
	}
	return ""
}
