package textio

import (
	"bufio"
	"github.com/omakoto/go-common/src/common"
	"os"
)

var (
	BufferedStdout = bufio.NewWriter(os.Stdout)
)

func init() {
	common.AtExit(func() {
		BufferedStdout.Flush()
	})
}
