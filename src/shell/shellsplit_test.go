package shell

import (
	"testing"

	"github.com/omakoto/zenlog-go/zenlog/util"
	"github.com/stretchr/testify/assert"
)

func TestShellSplit(t *testing.T) {
	inputs := []struct {
		source   string
		expected []string
	}{
		{"", util.Ar()},
		{"a", util.Ar("a")},
		{"aaa", util.Ar("aaa")},
		{"aaa b  ccc", util.Ar("aaa", "b", "ccc")},
		{"aaa 'b  b'  ccc", util.Ar("aaa", "'b  b'", "ccc")},
		{`aaa 'b  b'\''  d'  ccc`, util.Ar("aaa", `'b  b'\''  d'`, "ccc")},
		{"`ab\"`  ccc", util.Ar("`ab\"`", "ccc")},
		{`a\ \'\ \"`, util.Ar(`a\ \'\ \"`)},
		{`$HOME/abc`, util.Ar(`$HOME/abc`)},
		{`${HOME}/abc`, util.Ar(`${HOME}/abc`)},
		{`  $(cat  ok  "$(next   "de  f")")/abc  xyz`, util.Ar(`$(cat  ok  "$(next   "de  f")")/abc`, "xyz")},
		{`$ \`, util.Ar(`$`, `\`)},
		{`$`, util.Ar(`$`)},
		{`$'xyz' abc`, util.Ar(`$'xyz'`, `abc`)},
		{`$"xyz" abc`, util.Ar(`$"xyz"`, `abc`)},
		{`"\`, util.Ar(`"\`)},
		{`'a x ;' b`, util.Ar(`'a x ;'`, `b`)},
		{`cat|&grep>&ab#def  # commenct;abc`, util.Ar(`cat`, `|&`, `grep`, `>&`, `ab#def`, `# commenct;abc`)},
		{`echo $'a\xffb' # broken utf8`, util.Ar(`echo`, `$'a\xffb'`, `# broken utf8`)},
		{`cat fi\ le.txt|grep -V ^# >'out$$.txt' # Find non-comment lines.`,
			util.Ar(`cat`, `fi\ le.txt`, `|`, `grep`, `-V`, `^#`, `>`, `'out$$.txt'`, `# Find non-comment lines.`)},
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
