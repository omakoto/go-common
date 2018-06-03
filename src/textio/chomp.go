package textio

// Chomp returns an original slice with the last LF cut, if any.
func Chomp(s []byte) []byte {
	r, _ := Chomped(s)
	return r
}

// Chomped returns an original slice with the last LF cut and a slice containing LF, if any.
func Chomped(s []byte) ([]byte, []byte) {
	if len(s) == 0 {
		return s, nil
	}
	last := len(s) - 1
	if s[last] == '\n' {
		return s[0:last], s[last : last+1]
	}
	return s, nil
}
