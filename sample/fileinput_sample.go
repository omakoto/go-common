package main

import (
	"flag"
	"fmt"

	"github.com/omakoto/go-common/src/fileinput"
)

var (
	inline = flag.Bool("replace", false, "Inline replace")
)

func main() {
	flag.Parse()
	// fmt.Printf("ok!\n")

	for text, fi := range fileinput.FileInput(fileinput.Options{InlineReplace: *inline, Files: flag.Args()}) {
		fmt.Printf("%s:%d: %s\n", fi.Filename, fi.Line, text)
	}
}
