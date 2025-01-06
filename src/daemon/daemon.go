package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/omakoto/go-common/src/common"
)

const DaemonMarker = "__DAEMON_MARKER__"

type DaemonOptions struct {
	Cwd string
}

func Start() bool {
	return StartWithOptions(DaemonOptions{})
}

func StartWithOptions(options DaemonOptions) bool {
	var err error
	if options.Cwd == "" {
		options.Cwd, err = os.UserHomeDir()
		common.Checke(err)
	}
	if os.Getenv(DaemonMarker) == "" {
		doParent(options)
		return true
	} else {
		doChild(options)
		return false
	}
}

func doParent(_ DaemonOptions) {
	os.Setenv(DaemonMarker, "x")

	bin, err := filepath.Abs(os.Args[0])
	common.Checke(err)

	cmd := exec.Command(bin, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	common.Debugf("Spawning daemon... %v\n", cmd)
	err = cmd.Start()
	common.Check(err, "failed to spawn a daemon process")
}

func doChild(options DaemonOptions) {
	os.Unsetenv(DaemonMarker)

	os.Chdir(options.Cwd)

	signal.Ignore(syscall.SIGHUP)

	fmt.Printf("Daemon started with pid %d\n", os.Getpid())
}
