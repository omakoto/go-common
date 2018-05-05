package textio

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSanitizer(t *testing.T) {
	s := NewSanitizer()
	tests := []struct {
		input    string
		expected string
	}{
		{``, ``},
		{"ab\ncd\r\ndef\r\rxyz\x08\a\x08", "ab\ncd\ndef\n\nxyz^H^H"},
	}
	for _, v := range tests {
		a := s.Sanitize([]byte(v.input))
		assert.Equal(t, v.expected, string(a), "Input="+v.input)
	}
}
