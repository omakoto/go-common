package utils

import (
	"github.com/omakoto/go-common/src/common"
	"strconv"
)

func MustParseInt(val string, base int) int64 {
	v, err := strconv.ParseInt(val, base, 64)
	common.Checkf(err, "invalid string \"%s\"", val)
	return v
}

func ParseInt(val string, base int, defaultValue int64) int64 {
	v, err := strconv.ParseInt(val, base, 64)
	if err != nil {
		return defaultValue
	}
	return v
}

func MustParseFloat(val string) float64 {
	v, err := strconv.ParseFloat(val, 64)
	common.Checkf(err, "invalid string \"%s\"", val)
	return v
}

func ParseFloat(val string, defaultValue float64) float64 {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultValue
	}
	return v
}
