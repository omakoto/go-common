package shell

import (
	"bytes"
	"strings"

	"github.com/omakoto/go-common/src/utils"
)

var (
	reNeedsEscaping = utils.NewLazyRegexp(`[^a-zA-Z0-9\-\.\_\/\+\^\,\=\:]`)
)

// Escape a string for shell.
func Escape(s string) string {
	if !reNeedsEscaping.Pattern().MatchString(s) {
		return s
	}
	var buffer bytes.Buffer
	buffer.WriteString("'")
	buffer.WriteString(strings.Replace(s, `'`, `'\''`, -1))
	buffer.WriteString("'")
	return buffer.String()
}
