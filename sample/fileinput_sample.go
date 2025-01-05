package main

import (
	"fmt"

	"github.com/omakoto/go-common/src/fileinput"
)

func main() {
	// fmt.Printf("ok!\n")

	for text, fi := range fileinput.FileInput() {
		fmt.Printf("%s:%d: %s\n", fi.Filename, fi.Line, text)
	}
}
