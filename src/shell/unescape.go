package shell

import "bytes"

func hasQuote(text string) bool {
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '\'', '"':
			return true
		}
	}
	return false
}

func nextByte(text string, pos *int, nextByte *byte) bool {
	if *pos < len(text) {
		*nextByte = text[*pos]
		*pos++
		return true
	}
	return false
}

func Unescape(text string) string {
	if !hasQuote(text) {
		return text
	}
	buffer := bytes.NewBuffer(make([]byte, 0, len(text)))
	pos := 0
	var b byte
	for nextByte(text, &pos, &b) {
		if b == '\'' {
			for nextByte(text, &pos, &b) {
				if b == '\'' {
					break
				}
				buffer.WriteByte(b)
			}
			continue
		}
		if b == '"' {
			for nextByte(text, &pos, &b) {
				if b == '"' {
					break
				}
				if b == '\\' {
					if nextByte(text, &pos, &b) {
						b = text[pos]
					}
				}
				buffer.WriteByte(b)
			}
			continue
		}
		if b == '$' && (pos < len(text)) && text[pos] == '\'' {
			// C-like string
			pos++
			pos = UnescapeCLike(text, buffer, pos)
			continue
		}
		buffer.WriteByte(b)
	}
	return buffer.String()
}

func UnescapeCLike(text string, buffer *bytes.Buffer, pos int) int {
	var b byte

	for nextByte(text, &pos, &b) {
		if b == '\'' {
			break
		}
		if b == '\\' {
			if nextByte(text, &pos, &b) {
				switch b {
				case '"':
					buffer.WriteByte('"')
				case '\'':
					buffer.WriteByte('\'')
				case '\\':
					buffer.WriteByte('\\')
				case '?':
					buffer.WriteByte('?')
				case 'a':
					buffer.WriteByte('\a')
				case 'b':
					buffer.WriteByte('\b')
				case 'e', 'E':
					buffer.WriteByte('\x1b')
				case 'f':
					buffer.WriteByte('\f')
				case 'n':
					buffer.WriteByte('\n')
				case 'r':
					buffer.WriteByte('\r')
				case 't':
					buffer.WriteByte('\t')
				case 'v':
					buffer.WriteByte('\v')
				case 'c': // TODO -- Ctrl+?
					if nextByte(text, &pos, &b) {
						buffer.WriteByte(b - 'a')
					}
				case 'x': // TODO -- hex 00-ff
				case 'u': // TODO -- utf8 0000-ffff
				case 'U': // TODO -- utf8 00000000-ffffffff
				default: // unrecognized escape char.
					buffer.WriteByte('\\')
					buffer.WriteByte(b)
				}
			}

			continue
		}
		buffer.WriteByte(b)
	}
	return pos
}
