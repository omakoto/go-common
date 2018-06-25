// Shell unescape for posix + bash extension ($'..' and $"..")
package shell

import "bytes"

func hasQuote(text string) bool {
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '\'', '"', '\\':
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

func eatHex(text string, pos *int, maxLen int) (uint64, bool) {
	startPos := *pos
	var ret uint64
	for *pos < len(text) && maxLen > 0 {
		b := text[*pos]
		var v uint8
		if '0' <= b && b <= '9' {
			v = (b - '0')
		} else if 'a' <= b && b <= 'f' {
			v = (b - 'a' + 10)
		} else if 'A' <= b && b <= 'F' {
			v = (b - 'A' + 10)
		} else {
			break
		}
		ret *= 16
		ret += uint64(v)
		*pos++
		maxLen--
	}
	return ret, *pos > startPos
}

func eatOctal(text string, pos *int, maxLen int) (uint64, bool) {
	startPos := *pos
	var ret uint64
	for *pos < len(text) && maxLen > 0 {
		b := text[*pos]
		var v uint8
		if '0' <= b && b <= '7' {
			v = (b - '0')
		} else {
			break
		}
		ret *= 8
		ret += uint64(v)
		*pos++
		maxLen--
	}
	return ret, *pos > startPos
}

func Unescape(text string) string {
	if !hasQuote(text) {
		return text
	}
	buffer := bytes.NewBuffer(make([]byte, 0, len(text)))
	pos := 0
	var b byte
	for nextByte(text, &pos, &b) {
		if b == '\\' {
			if nextByte(text, &pos, &b) {
				buffer.WriteByte(b)
			}
			continue
		}
		if b == '\'' {
			for nextByte(text, &pos, &b) {
				if b == '\'' {
					break
				}
				buffer.WriteByte(b)
			}
			continue
		}
		if (b == '$' && (pos < len(text)) && text[pos] == '"') || (b == '"') {
			if b == '$' {
				pos++
			}
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
				case 'c':
					if nextByte(text, &pos, &b); b != '\'' {
						buffer.WriteByte(b & 0x1f)
					} else {
						buffer.Write([]byte("\\c"))
					}
				case 'x':
					if v, ok := eatHex(text, &pos, 2); ok {
						buffer.WriteByte(uint8(v))
					} else {
						buffer.Write([]byte("\\c"))
					}
				case 'u':
					if v, ok := eatHex(text, &pos, 4); ok {
						buffer.WriteRune(rune(v))
					} else {
						buffer.Write([]byte("\\u"))
					}
				//case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				//	if v, ok := eatOctal(text, &pos, 3); ok {
				//		buffer.WriteByte(uint8(v))
				//	} else {
				//		buffer.WriteByte('\\')
				//		buffer.WriteByte(b)
				//	}
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

func UnescapeBytes(b []byte) []byte {
	return []byte(Unescape(string(b)))
}
