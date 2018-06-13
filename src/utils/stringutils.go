package utils

import (
	"bytes"
	"strings"
)

func IndexByteOrLen(s string, c byte) int {
	ret := strings.IndexByte(s, c)
	if ret >= 0 {
		return ret
	}
	return len(s)
}

func BytesIndexByteOrLen(s []byte, c byte) int {
	ret := bytes.IndexByte(s, c)
	if ret >= 0 {
		return ret
	}
	return len(s)
}

func StringByteAt(s string, index int) byte {
	l := len(s)
	if index < 0 {
		index = l + index
	}
	if index < 0 || index >= l {
		return 0
	}
	return s[index]
}

func BytesByteAt(s string, index int) byte {
	l := len(s)
	if index < 0 {
		index = l + index
	}
	if index < 0 || index >= l {
		return 0
	}
	return s[index]
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

func StringsToBytes(s []string) [][]byte {
	ret := make([][]byte, 0)
	for _, v := range s {
		ret = append(ret, []byte(v))
	}
	return ret
}

func BytesToStrings(b [][]byte) []string {
	ret := make([]string, 0)
	for _, v := range b {
		ret = append(ret, string(v))
	}
	return ret
}
