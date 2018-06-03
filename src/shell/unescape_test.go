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
		{`a`, `a`},
		{`a b c`, `a b c`},
		{`a\ `, `a `},
		{`a\ b\  "''" 'xx\"' \'`, `a b  '' xx\" '`},
		{`a\ \b \b`, `a b b`},

		{`$''`, ``},

		{`$'abc'def`, `abcdef`},
		{`$'\"\'\\\a\b\e\E\f\n\r\t\v\q\ca\xbX'`, "\"'\\\a\b\x1b\x1b\f\n\r\t\v\\q\x01\x0bX"},

		{`$'\c@\cA\cZ\c[\ca\cz\c0'`, "\x00\x01\x1a\x1b\x01\x1a\x10"},
		{`$'\c'`, `\c`},
		{`$'\\c@'`, `\c@`},

		{`$'\x0 \xff \xfff'`, "\x00 \xff \xfff"},

		{`$'\u1\uz1\u7e\u56fdfX'`, "\x01\\uz1~\u56fdfX"},
		{`$'\U1\Uz1\U7e\U56fd\U1F466X'`, "\x01\\Uz1\x7e\u56fd\U0001F466X"},
		{`$'\U0001F466X'`, "\U0001F466X"},
		{`$'\U0001F466f'`, "\U0001F466f"},

		{`$""`, ``},
		{`$"aaa bbb ccc"`, `aaa bbb ccc`},
	}
	for _, v := range inputs {
		assert.Equal(t, v.expected, Unescape(v.source), v.source)
	}
}
