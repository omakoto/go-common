package runner

import (
	"testing"
)

func TestWrapperPath(t *testing.T) {
	cases := []struct {
		expected   string
		mainSource string
		o          Options
	}{
		{"a", "a.go", Options{}},
		{"wrapper", "a.go", Options{WrapperPath: "wrapper"}},
		{"../../wrapper", "a.go", Options{WrapperPath: "../../wrapper"}},
	}

	for _, i := range cases {
		actual := i.o.wrapperPath(i.mainSource)
		if actual != i.expected {
			t.Errorf("%v.wrapperPath(%v) expected to be %s but was %s", i.o, i.mainSource, i.expected, actual)
		}
	}
}

func TestRelPath(t *testing.T) {
	cases := []struct {
		expected    string
		mainSource  string
		wrapperPath string
	}{
		{"a.go", "a.go", "a"},
		{"a.go", "a.go", "wrapper"},
		{"a.go", "a/b/c/a.go", "wrapper"},
		{"c/d/a.go", "/a/b/c/d/a.go", "../../wrapper"},
	}

	for _, i := range cases {
		actual := getRelPath(i.mainSource, i.wrapperPath)
		if actual != i.expected {
			t.Errorf("getRelPath(%v, %v) expected to be %s but was %s", i.mainSource, i.wrapperPath, i.expected, actual)
		}
	}
}
