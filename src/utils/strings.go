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
	if index >= l {
		return 0
	}
	return s[index]
}

func BytesByteAt(s string, index int) byte {
	l := len(s)
	if index < 0 {
		index = l + index
	}
	if index >= l {
		return 0
	}
	return s[index]
}
