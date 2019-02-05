package utils

import (
	"github.com/omakoto/go-common/src/common"
	"io/ioutil"
	"os"
	"os/exec"
)

// ReadPdfAsText reads a given PDF file as text.
func ReadPdfAsText(file string, keepLayout bool) ([]byte, error) {
	args := make([]string, 0)
	if keepLayout {
		args = append(args, "-layout")
	}
	args = append(args, file, "-")
	cmd := exec.Command("pdftotext", args...)

	cmd.Stdin = nil
	cmd.Stderr = os.Stderr

	in, err := cmd.StdoutPipe()
	common.Checke(err)
	defer in.Close()

	err = cmd.Start()
	common.Checkf(err, "Failed executing pdftotext: %v", err)

	return ioutil.ReadAll(in)
}
