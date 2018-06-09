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

// StringChomp returns an original slice with the last LF cut, if any.
func StringChomp(s string) string {
	r, _ := StringChomped(s)
	return r
}

// StringChomped returns an original slice with the last LF cut and a slice containing LF, if any.
func StringChomped(s string) (string, string) {
	if len(s) == 0 {
		return s, ""
	}
	last := len(s) - 1
	if s[last] == '\n' {
		return s[0:last], s[last : last+1]
	}
	return s, ""
}
