package shell

// Shell tokenizer for posix + bash extension ($'..' and $"..")

import (
	"bytes"
)

type Token struct {
	Word  string
	Index int
}

type splitter struct {
	// Input.
	text []rune
	next int

	// State
	wasSpecial bool
	start      int
	buffer     bytes.Buffer
	hasRunes   bool

	// Result
	resultTokens []Token
}

func newSplitter(text string) splitter {
	return splitter{
		text:         []rune(text),
		buffer:       bytes.Buffer{},
		start:        -1,
		resultTokens: make([]Token, 0),
	}
}

func (s *splitter) peek() (rune, bool, int) {
	if s.next < len(s.text) {
		return s.text[s.next], true, s.next
	}
	return '\x00', false, -1
}

func (s *splitter) read() (rune, bool, int) {
	r, ok, index := s.peek()

	if ok {
		s.next++
	}
	return r, ok, index
}

func (s *splitter) pushRuneNoSpecial(r rune, index int) {
	s.buffer.WriteRune(r)
	s.hasRunes = true
	if s.start < 0 {
		s.start = index
	}
}

func (s *splitter) pushRune(r rune, index int) {
	s.pushRuneNoSpecial(r, index)
	s.wasSpecial = isSpecialChar(r)
}

func (s *splitter) onWordBoundary() {
	if s.hasRunes {
		s.resultTokens = append(s.resultTokens, Token{s.buffer.String(), s.start})
		s.hasRunes = false
		s.start = -1
		s.wasSpecial = false
		s.buffer = bytes.Buffer{}
	}
}

func (s *splitter) eatSingleQuote() {
	for {
		r, ok, index := s.read()
		if !ok {
			return
		}
		if r == '\'' {
			s.pushRune(r, index)
			return
		}
		s.pushRune(r, index)
	}
}

func (s *splitter) eatDoubleQuote(end rune) {
	for {
		r, ok, index := s.read()
		if !ok {
			return
		}
		if r == end {
			s.pushRune(r, index)
			return
		}
		if s.maybeEatDollar(r, index) {
			continue
		}
		if r == '\\' {
			s.pushRuneNoSpecial(r, index)
			r, ok, _ = s.read()
			if !ok {
				return
			}
		}
		s.pushRune(r, index)
	}
}

func (s *splitter) maybeEatDollar(r rune, index int) bool {
	if r == '$' {
		s.pushRune(r, index)

		next, ok, index := s.peek()
		if !ok || isWhitespace(next) {
			return true
		}
		if next == '(' {
			s.read()
			s.pushRuneNoSpecial(next, index)
			s.tokenize(')')
			return true
		}
		if next == '{' {
			s.read()
			s.pushRuneNoSpecial(next, index)
			s.tokenize('}')
			return true
		}
		if next == '\'' {
			s.read()
			s.pushRuneNoSpecial(next, index)
			s.eatSingleQuote()
			return true
		}
		if next == '"' {
			s.read()
			s.pushRuneNoSpecial(next, index)
			s.eatDoubleQuote('"')
			return true
		}
		return true
	}
	return false
}

func isWhitespace(r rune) bool {
	switch r {
	case ' ', '\t', '\r', '\n', '\v':
		return true
	}
	return false
}

func isSpecialChar(r rune) bool {
	switch r {
	case ';', '!', '<', '>', '(', ')', '|', '&':
		return true
	}
	return false
}

func isCommandSeparatorChar(r rune) bool {
	switch r {
	case ';', '(', ')', '|', '&':
		return true
	}
	return false
}

func (s *splitter) tokenize(end int) {
	for {
		r, ok, index := s.read()
		if !ok {
			break
		}
		if end >= 0 && end == int(r) {
			s.pushRuneNoSpecial(r, index)
			break
		}
		if end < 0 && isWhitespace(r) {
			s.onWordBoundary()
			continue
		}
		if s.wasSpecial != isSpecialChar(r) {
			s.onWordBoundary()
		}
		if r == '\\' {
			s.pushRune(r, index)
			r, ok, index = s.read()
			if !ok {
				break
			}
			s.pushRune(r, index)
			continue
		}
		if r == '\'' {
			s.pushRune(r, index)
			s.eatSingleQuote()
			continue
		}
		if r == '"' {
			s.pushRune(r, index)
			s.eatDoubleQuote('"')
			continue
		}
		if r == '`' {
			s.pushRune(r, index)
			s.eatDoubleQuote('`')
			continue
		}
		if s.maybeEatDollar(r, index) {
			continue
		}
		// If # follows a whitespace,the rest is a comment.
		if !s.hasRunes && r == '#' {
			for {
				s.pushRuneNoSpecial(r, index)
				r, ok, index = s.read()
				if !ok {
					s.onWordBoundary()
					return
				}
			}
		}
		s.pushRune(r, index)
	}
	if end < 0 {
		s.onWordBoundary()
	}
}

// Split splits a whole command line into tokens.
// Example: "cat fi\ le.txt|grep -V ^# >'out$$.txt' # Find non-comment lines."
// -> output: cat, fi\ le.txt, |, grep, -V, ^#, >, 'out$$.txt', # Find non-comment lines.
func Split(text string) []string {
	s := newSplitter(text)
	s.tokenize(-1)

	ret := make([]string, 0, len(s.resultTokens))
	for _, s := range s.resultTokens {
		ret = append(ret, s.Word)
	}
	return ret
}

func SplitToTokens(text string) []Token {
	s := newSplitter(text)
	s.tokenize(-1)
	return s.resultTokens
}

func IsCommandSeparator(text string) bool {
	for _, r := range text {
		if !isCommandSeparatorChar(r) {
			return false
		}
	}
	return true
}
