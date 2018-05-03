package textio

// Chomp returns an original slice with the last LF cut, if any.
func Chomp(s []byte) []byte {
	if len(s) == 0 {
		return s
	}
	last := len(s)-1
	if s[last] == '\n' {
		return s[0:last]
	}
	return s
}
