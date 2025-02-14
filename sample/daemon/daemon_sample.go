package main

import (
	"flag"
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/daemon"
	"os"
	"time"
)

var (
	doDaemon = flag.Bool("d", false, "start as daemon")
	doQuit   = flag.Bool("q", false, "stop daemon")
)

func main() {
	flag.Parse()
	if *doQuit {
		if daemon.Stop() {
			os.Exit(0)
		} else {
			common.Warn("Failed to stop daemon")
			os.Exit(0)
		}
	}

	if *doDaemon {
		if daemon.Start() {
			fmt.Printf("Started daemon.\n")
			os.Exit(0) // Parent process
		}
		actualMain()
	} else {
		// Run directly
		fmt.Printf("Start...\n")
		actualMain()
	}
}

func actualMain() {
	fmt.Printf("Main started... pid=%d\n", os.Getpid())
	for {
		time.Sleep(2 * time.Second)
		fmt.Printf("Still running...\n")
	}
}
