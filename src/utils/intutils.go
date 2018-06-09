package utils

import "fmt"

func Clip64(n, start, last int64) int64 {
	if start <= n {
		if n < last {
			return n
		}
		return last - 1
	}
	return start
}

func Wrap64(n, len int64) int64 {
	if n < 0 {
		n += len
	}
	if 0 <= n && n < len {
		return n
	}
	panic(fmt.Errorf("index out of range: n=%d, len=%d", n, len))
}

func Clip(n, start, last int) int {
	return int(Clip64(int64(n), int64(start), int64(last)))
}

func Clip32(n, start, last int32) int32 {
	return int32(Clip64(int64(n), int64(start), int64(last)))
}

func Wrap(n, len int) int {
	return int(Wrap64(int64(n), int64(len)))
}

func Wrap32(n, len int32) int32 {
	return int32(Wrap64(int64(n), int64(len)))
}
