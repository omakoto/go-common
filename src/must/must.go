package must

import "github.com/omakoto/go-common/src/common"

func Must(err error) {
	common.Checke(err)
}

func Must2[V any](value V, err error) V {
	common.Checke(err)
	return value
}

func Must3[V1 any, V2 any](value1 V1, value2 V2, err error) (V1, V2) {
	common.Checke(err)
	return value1, value2
}

type WithContext struct {
	err error
}

func With(err error) *WithContext {
	return &WithContext{err}
}

func (c *WithContext) Checkf(fmt string, args ...interface{}) {
	common.Checkf(c.err, fmt, args...)
}

type With2Context[V any] struct {
	v   V
	err error
}

func With2[V any](value V, err error) *With2Context[V] {
	return &With2Context[V]{value, err}
}

func (c *With2Context[V]) Checkf(fmt string, args ...interface{}) V {
	common.Checkf(c.err, fmt, args...)
	return c.v
}

type With3Context[V1, V2 any] struct {
	v1  V1
	v2  V2
	err error
}

func With3[V1, V2 any](value1 V1, value2 V2, err error) *With3Context[V1, V2] {
	return &With3Context[V1, V2]{value1, value2, err}
}

func (c *With3Context[V1, V2]) Checkf(fmt string, args ...interface{}) (V1, V2) {
	common.Checkf(c.err, fmt, args...)
	return c.v1, c.v2
}
