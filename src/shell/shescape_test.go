package shell

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShescape(t *testing.T) {
	inputs := []struct {
		expected string
		source   string
	}{
		{"", ""},
		{"abc", "abc"},
		{"'abc '", "abc "},
		{"'abc def \" '\\'' xyz '\\'''", "abc def \" ' xyz '"},
	}
	for _, v := range inputs {
		actual := Escape(v.source)
		assert.Equal(t, v.expected, actual)
	}
}

func TestShescapeNoQuotes(t *testing.T) {
	inputs := []struct {
		expected string
		source   string
	}{
		{``, ""},
		{`abc`, "abc"},
		{`abc\ `, "abc "},
		{`abc\ def\ \"\ \'\ xyz\ \'あいう\ `, `abc def " ' xyz 'あいう `},
	}
	for _, v := range inputs {
		actual := EscapeNoQuotes(v.source)
		assert.Equal(t, v.expected, actual)
	}
}
