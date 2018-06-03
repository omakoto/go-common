package shell

import (
	"github.com/omakoto/go-common/src/common"
	"os"
	"path/filepath"
)

type Proxy interface {
	GetCommandLine() (commandLine string, cursorPos int)
	PrintUpdateCommandLineEvalStr(commandLine string, cursorPos int)
	Split(text string) []Token
	Escape(text string) string
	Unescape(text string) string
}

type defaultShellProxy struct {
}

func (s *defaultShellProxy) GetCommandLine() (string, int) {
	return "", 0
}

func (s *defaultShellProxy) PrintUpdateCommandLineEvalStr(commandLine string, cursorPos int) {
}

func (s *defaultShellProxy) Split(text string) []Token {
	return SplitToTokens(text) // TODO Use posix-compat version?
}

func (s *defaultShellProxy) Escape(text string) string {
	return Escape(text) // TODO Use posix-compat version?
}

func (s *defaultShellProxy) Unescape(text string) string {
	return Unescape(text) // TODO Use posix-compat version?
}

func GetSupportedProxy() Proxy {
	shell := filepath.Base(os.Getenv("SHELL"))

	switch shell {
	case "bash":
		return GetBashProxy()
	case "zsh":
		return GetZshProxy()
	}
	return nil
}

func MustGetSupportedProxy() Proxy {
	sh := GetSupportedProxy()
	common.OrFatalf(sh != nil, "Unsupported shell.\n")
	return sh
}

func GetProxy() Proxy {
	ret := GetSupportedProxy()
	if ret != nil {
		return ret
	}
	return &defaultShellProxy{}
}
