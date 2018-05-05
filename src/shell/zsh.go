package shell

import (
	"fmt"
	"os"
	"strconv"
)

const (
	zleBuffer = "BUFFER"
	zleCursor = "CURSOR"
)

type ZshProxy struct {
}

func GetZshProxy() Proxy {
	return &ZshProxy{}
}

func (b *ZshProxy) GetCommandLine() (string, int) {
	s := os.Getenv(zleBuffer)
	l, err := strconv.Atoi(os.Getenv(zleCursor))
	if err != nil || l < 0 {
		l = len(s)
	}
	return s, l
}

func (b *ZshProxy) PrintUpdateCommandLineEvalStr(commandLine string, cursorPos int) {
	fmt.Print(zleBuffer)
	fmt.Print("=")
	fmt.Println(Escape(commandLine))

	fmt.Print(zleCursor)
	fmt.Print("=")
	fmt.Println(strconv.Itoa(cursorPos))
}

func (s *ZshProxy) Split(text string) []Token {
	return SplitToTokens(text) // TODO Not sure if it's 100% compatible to bash.
}

func (s *ZshProxy) Escape(text string) string {
	return Escape(text) // TODO Not sure if it's 100% compatible to bash.
}

func (s *ZshProxy) Unescape(text string) string {
	return Unescape(text) // TODO Not sure if it's 100% compatible to bash.
}
