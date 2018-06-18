package shell

import (
	"bytes"
)

var (
	// perl -e 'for ($i = 0; $i < 256; $i++) {print "0 /*", (32 <=$i && $i <= 126 ? chr($i) : "-"), "*/,"; print "\n" if ($i+1) % 16 == 0}'
	types = []int{
		1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1, /*-*/
		1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1, /*-*/
		1 /* */, 1 /*!*/, 1 /*"*/, 1 /*#*/, 1 /*$*/, 0 /*%*/, 1 /*&*/, 2 /*'*/, 1 /*(*/, 1 /*)*/, 1 /***/, 0 /*+*/, 0 /*,*/, 0 /*-*/, 0 /*.*/, 0, /*/*/
		0 /*1*/, 0 /*1*/, 0 /*2*/, 0 /*3*/, 0 /*4*/, 0 /*5*/, 0 /*6*/, 0 /*7*/, 0 /*8*/, 0 /*9*/, 0 /*:*/, 1 /*;*/, 1 /*<*/, 0 /*=*/, 1 /*>*/, 1, /*?*/
		0 /*@*/, 0 /*A*/, 0 /*B*/, 0 /*C*/, 0 /*D*/, 0 /*E*/, 0 /*F*/, 0 /*G*/, 0 /*H*/, 0 /*I*/, 0 /*J*/, 0 /*K*/, 0 /*L*/, 0 /*M*/, 0 /*N*/, 0, /*O*/
		0 /*P*/, 0 /*Q*/, 0 /*R*/, 0 /*S*/, 0 /*T*/, 0 /*U*/, 0 /*V*/, 0 /*W*/, 0 /*X*/, 0 /*Y*/, 0 /*Z*/, 1 /*[*/, 1 /*\*/, 1 /*]*/, 1 /*^*/, 0, /*_*/
		1 /*`*/, 0 /*a*/, 0 /*b*/, 0 /*c*/, 0 /*d*/, 0 /*e*/, 0 /*f*/, 0 /*g*/, 0 /*h*/, 0 /*i*/, 0 /*j*/, 0 /*k*/, 0 /*l*/, 0 /*m*/, 0 /*n*/, 0, /*o*/
		0 /*p*/, 0 /*q*/, 0 /*r*/, 0 /*s*/, 0 /*t*/, 0 /*u*/, 0 /*v*/, 0 /*w*/, 0 /*x*/, 0 /*y*/, 0 /*z*/, 1 /*{*/, 1 /*|*/, 1 /*}*/, 1 /*~*/, 1, /*-*/
		1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1, /*-*/
		1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1, /*-*/
		1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1, /*-*/
		1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1, /*-*/
		1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1, /*-*/
		1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1, /*-*/
		1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1, /*-*/
		1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1 /*-*/, 1, /*-*/
	}
)

// Escape escapes a string for shell.
func Escape(s string) string {
	if isSafe(s) {
		return s
	}
	buffer := bytes.NewBuffer(make([]byte, 0, len(s)*2))
	buffer.WriteByte('\'')
	for i := 0; i < len(s); i++ {
		b := s[i]
		switch types[b] {
		case 0, 1:
			buffer.WriteByte(b)
		case 2:
			buffer.WriteString(`'\''`)
		}
	}
	buffer.WriteByte('\'')
	return buffer.String()
}

// Escape escapes a string with backslashes.
func EscapeNoQuotes(s string) string {
	if isSafe(s) {
		return s
	}
	buffer := bytes.NewBuffer(make([]byte, 0, len(s)*2))
	for _, r := range(s) {
		if r < 128 && types[r] != 0 {
			buffer.WriteByte('\\')
		}
		buffer.WriteRune(r)
	}
	return buffer.String()
}

func isSafe(s string) bool {
	for i := 0; i < len(s); i++ {
		if types[s[i]] > 0 {
			return false
		}
	}
	return true
}

// EscapeBytes escapes a byte array for shell.
func EscapeBytes(s []byte) []byte {
	if isSafeBytes(s) {
		return s
	}
	buffer := bytes.NewBuffer(make([]byte, 0, len(s)*2))
	buffer.WriteByte('\'')
	for i := 0; i < len(s); i++ {
		b := s[i]
		t := types[b]
		switch t {
		case 0, 1:
			buffer.WriteByte(b)
		case 2:
			buffer.WriteString(`'\''`)
		}
	}
	buffer.WriteByte('\'')
	return buffer.Bytes()
}

// Escape escapes a byte array with backslashes.
func EscapeBytesNoQuotes(s []byte) []byte {
	return []byte(EscapeNoQuotes(string(s)))
}

func isSafeBytes(s []byte) bool {
	for i := 0; i < len(s); i++ {
		if types[s[i]] > 0 {
			return false
		}
	}
	return true
}
