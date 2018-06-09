package utils

import (
	"github.com/omakoto/go-common/src/common"
	"strconv"
)

func MustParseInt64(val string, base int) int64 {
	v, err := strconv.ParseInt(val, base, 64)
	common.Checkf(err, "invalid string \"%s\"", val)
	return v
}

func ParseInt64(val string, base int, defaultValue int64) int64 {
	v, err := strconv.ParseInt(val, base, 64)
	if err != nil {
		return defaultValue
	}
	return v
}

func MustParseInt(val string, base int) int {
	return int(MustParseInt64(val, base))
}

func ParseInt(val string, base int, defaultValue int) int {
	return int(ParseInt64(val, base, int64(defaultValue)))
}

func MustParseInt32(val string, base int) int32 {
	return int32(MustParseInt64(val, base))
}

func ParseInt32(val string, base int, defaultValue int32) int32 {
	return int32(ParseInt64(val, base, int64(defaultValue)))
}

func MustParseFloat64(val string) float64 {
	v, err := strconv.ParseFloat(val, 64)
	common.Checkf(err, "invalid string \"%s\"", val)
	return v
}

func ParseFloat64(val string, defaultValue float64) float64 {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultValue
	}
	return v
}

func MustParseFloat32(val string) float32 {
	return float32(MustParseFloat64(val))
}

func ParseFloat32(val string, defaultValue float32) float32 {
	return float32(ParseFloat64(val, float64(defaultValue)))
}
