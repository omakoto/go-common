package fileutils

import (
	"testing"
)

func TestWrapperPath(t *testing.T) {
	cases := []struct {
		expected      bool
		errorExpected bool
		file1, file2  string
	}{
		{true, false, "/tmp", "/"},
		{false, false, "/", "/tmp"},
		{false, false, "/tmpxxx", "/"},
		{false, true, "/tmp", "/xxxx"},
	}

	for _, i := range cases {
		actual, err := NewerThan(i.file1, i.file2)
		if i.errorExpected && err == nil {
			t.Errorf("NewerThan(%v, %v) expected to return err, but didn't", i.file1, i.file2)
		} else if !i.errorExpected && err != nil {
			t.Errorf("NewerThan(%v, %v) expected to succeed, but returned error ", i.file1, i.file2)
		} else if actual != i.expected {
			t.Errorf("NewerThan(%v, %v) expected to return %v, but %v", i.file1, i.file2, i.expected, actual)
		}
	}
}
