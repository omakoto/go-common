package textio

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

// ExpandTab expands \ts in a string. Multiline strings are allowed.
func ExpandTab(text string, tabWidth int) string {
	return string(ExpandTabBytes([]byte(text), tabWidth))
}

// ExpandTabBytes expands \ts in a []byte. Multiline strings are allowed.
func ExpandTabBytes(text []byte, tabWidth int) []byte {
	in := bufio.NewReader(bytes.NewBuffer(text))
	out := bytes.NewBuffer(make([]byte, 0, len(text)))

	for {
		line, err := in.ReadBytes(byte('\n'))
		if len(line) > 0 {
			expandTabSingle(line, out, tabWidth)
		}
		if err == io.EOF {
			break
		} else if err != nil {
			panic(fmt.Sprintf("shouldn't get errors: %s", err.Error()))
		}
	}
	return out.Bytes()
}

func expandTabSingle(line []byte, out *bytes.Buffer, tabWidth int) {
	outIndex := 0
	for _, ch := range line {
		if ch != '\t' {
			out.WriteByte(ch)
			outIndex++
			continue
		}
		fill := tabWidth - (outIndex % tabWidth)
		for i := 0; i < fill; i++ {
			out.WriteByte(' ')
		}
		outIndex += fill
	}
}
