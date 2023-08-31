package main

import (
	"fmt"
	"github.com/omakoto/go-common/src/cmdchain"
)

func main() {
	//fmt.Printf("on!\n")
	cmd := cmdchain.New().Command("bash", "-c", "for n in {0..9}; do echo $n; sleep 0.2 ; done")

	//cmd.MustRunAndWait()
	// fmt.Printf("%s", cmd.MustRunAndGetString())
	cmd.MustRunAndStreamStrings(func(s string) {
		fmt.Printf("%s", s)
	})
}
