package shell

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeSingleQuote(t *testing.T) {
	inputs := []struct {
		source   string
		expected string
	}{
		{``, ``},
		{`'`, ``},
		{`''`, ``},
		{`'''`, ``},
		{`''''`, ``},
		{`'abc'`, `abc`},
	}
	for _, v := range inputs {
		actual := Unescape(v.source)
		assert.Equal(t, v.expected, actual, v.source)
	}
}
