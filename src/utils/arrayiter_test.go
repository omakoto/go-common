package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasic(t *testing.T) {
	{
		i := Iter([]string{"a", "b"})
		assert.Equal(t, "a", *i())
		assert.Equal(t, "b", *i())
		assert.Equal(t, (*string)(nil), i())
		assert.Equal(t, (*string)(nil), i())
	}

	{
		i := Iter([]string{})
		assert.Equal(t, (*string)(nil), i())
		assert.Equal(t, (*string)(nil), i())
	}
}
