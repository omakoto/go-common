package textio

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExpandTab(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"    a", "    a"},
		{"\ta", "    a"},
		{" \ta", "    a"},
		{"  \ta", "    a"},
		{"   \ta", "    a"},
		{"     \ta", "        a"},
		{"\ta\tb", "    a   b"},
		{" \ta\tb", "    a   b"},
		{" \ta \tb", "    a   b"},
		{" \ta  \tb", "    a   b"},
		{" \ta   \tb", "    a       b"},
		{" \ta    \tb", "    a       b"},
	}
	for i, v := range tests {
		assert.Equal(t, v.expected, ExpandTab(v.input, 4), "#%d Input=\"%s\"", i, v.input)
	}
}
