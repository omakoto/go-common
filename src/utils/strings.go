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
