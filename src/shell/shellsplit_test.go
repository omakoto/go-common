package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ar(vals ...string) []string {
	return vals
}

func TestShellSplit(t *testing.T) {
	inputs := []struct {
		source   string
		expected []string
	}{
		{"", ar()},
		{"a", ar("a")},
		{"aaa", ar("aaa")},
		{"aaa b  ccc", ar("aaa", "b", "ccc")},
		{"aaa 'b  b'  ccc", ar("aaa", "'b  b'", "ccc")},
		{`aaa 'b  b'\''  d'  ccc`, ar("aaa", `'b  b'\''  d'`, "ccc")},
		{"`ab\"`  ccc", ar("`ab\"`", "ccc")},
		{`a\ \'\ \"`, ar(`a\ \'\ \"`)},
		{`$HOME/abc`, ar(`$HOME/abc`)},
		{`${HOME}/abc`, ar(`${HOME}/abc`)},
		{`  $(cat  ok  "$(next   "de  f")")/abc  xyz`, ar(`$(cat  ok  "$(next   "de  f")")/abc`, "xyz")},
		{`$ \`, ar(`$`, `\`)},
		{`$`, ar(`$`)},
		{`$'xyz' abc`, ar(`$'xyz'`, `abc`)},
		{`$"xyz" abc`, ar(`$"xyz"`, `abc`)},
		{`"\`, ar(`"\`)},
		{`'a x ;' b`, ar(`'a x ;'`, `b`)},
		{`cat|&grep>&ab#def  # commenct;abc`, ar(`cat`, `|&`, `grep`, `>&`, `ab#def`, `# commenct;abc`)},
		{`echo $'a\xffb' # broken utf8`, ar(`echo`, `$'a\xffb'`, `# broken utf8`)},
		{`cat fi\ le.txt|grep -V ^# >'out$$.txt' # Find non-comment lines.`,
			ar(`cat`, `fi\ le.txt`, `|`, `grep`, `-V`, `^#`, `>`, `'out$$.txt'`, `# Find non-comment lines.`)},
	}
	for _, v := range inputs {
		actual := Split(v.source)
		assert.Equal(t, v.expected, actual, "Source="+v.source)
	}
}

func ai(vals ...interface{}) []Token {
	ret := make([]Token, 0)

	start := 0

	for _, v := range vals {
		switch v.(type) {
		case int:
			start += v.(int)
		case string:
			s := v.(string)
			ret = append(ret, Token{s, start})
			start += len(s)
		}
	}
	return ret
}

func TestShellSplitWithIndexes(t *testing.T) {
	inputs := []struct {
		source   string
		expected []Token
	}{
		{"", ai()},
		{"a", ai("a")},
		{"aaa", ai("aaa")},
		{"  aaa", ai(2, "aaa")},
		{"  a", ai(2, "a")},
		{"aaa b  ccc", ai("aaa", 1, "b", 2, "ccc")},
		{"aaa 'b  b'  ccc", ai("aaa", 1, "'b  b'", 2, "ccc")},
		{`aaa 'b  b'\''  d'  ccc`, ai("aaa", 1, `'b  b'\''  d'`, 2, "ccc")},
		{"`ab\"`  ccc", ai("`ab\"`", 2, "ccc")},
		{`a\ \'\ \"`, ai(`a\ \'\ \"`)},
		{`$HOME/abc`, ai(`$HOME/abc`)},
		{`${HOME}/abc`, ai(`${HOME}/abc`)},
		{`  $(cat  ok  "$(next   "de  f")")/abc  xyz`, ai(2, `$(cat  ok  "$(next   "de  f")")/abc`, 2, "xyz")},
		{`$ \`, ai(`$`, 1, `\`)},
		{`$`, ai(`$`)},
		{`$'xyz' abc`, ai(`$'xyz'`, 1, `abc`)},
		{`$"xyz" abc`, ai(`$"xyz"`, 1, `abc`)},
		{`$'xyz  XYZ' abc`, ai(`$'xyz  XYZ'`, 1, `abc`)},
		{`$"xyz  XYZ" abc`, ai(`$"xyz  XYZ"`, 1, `abc`)},
		{`$'' abc`, ai(`$''`, 1, `abc`)},
		{`$"" abc`, ai(`$""`, 1, `abc`)},
		{`"\`, ai(`"\`)},
		{`'a x ;'     b`, ai(`'a x ;'`, 5, `b`)},
		{`cat|&grep>&ab#def  # commenct;abc`, ai(`cat`, `|&`, `grep`, `>&`, `ab#def`, 2, `# commenct;abc`)},
		{`echo $'a\xffb' # broken utf8`, ai(`echo`, 1, `$'a\xffb'`, 1, `# broken utf8`)},
		{`cat fi\ le.txt|grep -V ^# >'out$$.txt' # Find non-comment lines.`,
			ai(`cat`, 1, `fi\ le.txt`, `|`, `grep`, 1, `-V`, 1, `^#`, 1, `>`, `'out$$.txt'`, 1, `# Find non-comment lines.`)},
	}
	for _, v := range inputs {
		actual := SplitToTokens(v.source)
		assert.Equal(t, v.expected, actual, "Source="+v.source)
	}
}
