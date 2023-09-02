package main

import (
	"github.com/omakoto/go-common/src/runner"
	"os"
)

func main() {
	filename := os.Getenv("GOFILE")
	if filename == "" {
		panic("$GOFILE not set.")
	}

	runner.Generate(filename)
}
