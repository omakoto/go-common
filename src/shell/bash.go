package shell

import (
	"fmt"
	"os"
	"strconv"
)

const (
	readlineLine  = "READLINE_LINE"
	readlinePoint = "READLINE_POINT"
)

type BashProxy struct {
}

func GetBashProxy() Proxy {
	return &BashProxy{}
}

// GetCommandLine return the current command line and the cursor position from the READLINE_* environmental variables.
func (s *BashProxy) GetCommandLine() (string, int) {
	line := os.Getenv(readlineLine)
	l, err := strconv.Atoi(os.Getenv(readlinePoint))
	if err != nil || l < 0 {
		l = len(line)
	}
	return line, l
}

// PrintUpdateCommandLineEvalStr prints a string that can be evaled by bash to update the READLINE_* environmental variables
// to update the current command line.
func (s *BashProxy) PrintUpdateCommandLineEvalStr(commandLine string, cursorPos int) {
	fmt.Print(readlineLine)
	fmt.Print("=")
	fmt.Println(Escape(commandLine))

	fmt.Print(readlinePoint)
	fmt.Print("=")
	fmt.Println(strconv.Itoa(cursorPos))
}

func (s *BashProxy) Split(text string) []Token {
	return SplitToTokens(text)
}

func (s *BashProxy) Escape(text string) string {
	return Escape(text)
}

func (s *BashProxy) Unescape(text string) string {
	return Unescape(text)
}
