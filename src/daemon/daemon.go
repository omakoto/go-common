package daemon

import (
	"fmt"
	"github.com/omakoto/go-common/src/fileutils"
	"github.com/omakoto/go-common/src/must"
	"github.com/omakoto/go-common/src/textio"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/omakoto/go-common/src/common"
)

const DaemonMarker = "__DAEMON_MARKER__"

type Options struct {
	Cwd         string
	PidFilename string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (o *Options) getPidFile() string {
	pidFilename := o.PidFilename
	if pidFilename == "" {
		pidFilename = filepath.Base(os.Args[0]) + "_pid.txt"
	}
	return must.Must2(os.UserHomeDir()) + string(filepath.Separator) + "." + pidFilename
}

func Start() bool {
	return StartWithOptions(Options{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
}

func StartWithOptions(options Options) bool {
	if options.Cwd == "" {
		options.Cwd = must.Must2(os.UserHomeDir())
	}
	if os.Getenv(DaemonMarker) == "" {
		doParent(options)
		return true
	} else {
		doChild(options)
		return false
	}
}

func IsRunning(options Options) bool {
	return isRunning(getPid(options))
}

func doParent(options Options) {
	if isRunning(getPid(options)) {
		common.Panicf("Daemon is already running: options=%v", options)
	}
	bin := must.Must2(filepath.Abs(os.Args[0]))

	pidFile := options.getPidFile()

	os.Setenv(DaemonMarker, "x")
	cmd := exec.Command(bin, os.Args[1:]...)
	cmd.Stdin = options.Stdin
	cmd.Stdout = options.Stdout
	cmd.Stderr = options.Stderr

	common.Debugf("Spawning daemon... %v\n", cmd)
	err := cmd.Start()
	common.Check(err, "failed to spawn a daemon process")

	pid := cmd.Process.Pid
	must.Must(textio.WriteStringToFile(pidFile, strconv.Itoa(pid), 0600))
}

func doChild(options Options) {
	os.Unsetenv(DaemonMarker)

	os.Chdir(options.Cwd)

	signal.Ignore(syscall.SIGHUP)

	fmt.Printf("Daemon started with pid %d\n", os.Getpid())
}

func Stop() bool {
	return StopWithOptions(Options{})
}

func getPid(options Options) int {
	pidFile := options.getPidFile()

	if !fileutils.FileExists(pidFile) {
		common.Debugf("PID file %s doesn't exist", pidFile)
		return -1
	}
	pid, err := textio.ReadStringFromFile(pidFile)
	if err != nil {
		common.Warnf("PID file %s not readable: %s", pidFile, err)
		return -1
	}
	pid = textio.StringChomp(pid)
	ipid, err := strconv.Atoi(pid)
	if err != nil {
		common.Warnf("Invalid PID format in file %s: '%s'", pidFile, pid)
		return -1
	}

	common.Debugf("Previous PID: %d", ipid)
	return ipid
}

func isRunning(pid int) bool {
	if pid < 0 {
		return false
	}
	proc := fmt.Sprintf("/proc/%d/stat", pid)
	stat, err := textio.ReadStringFromFile(proc)
	if err != nil {
		common.Debugf("Failed to read %s: %s", proc, err)
		return false
	}
	stat = strings.Trim(stat, " \n")
	fields := strings.Split(stat, " ")
	if len(fields) < 2 {
		common.Debugf("Invalid format from %s: \"%s\"", proc, stat)
		return false
	}
	myName := filepath.Base(os.Args[0])
	return fields[1] == "("+myName+")"
}

func StopWithOptions(options Options) bool {
	pid := getPid(options)
	if pid < 0 {
		return true
	}
	if !isRunning(pid) {
		return true
	}
	err := syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		common.Warnf("Failed to send signal to pid %d: %s", pid, err)
		return false
	}

	timeout := time.Now().Add(time.Second * 5)
	for {
		time.Sleep(50 * time.Millisecond)
		if !isRunning(pid) {
			return true
		}
		if time.Now().After(timeout) {
			common.Warnf("Pid %d didn't terminate.", pid)
			return false
		}
	}
}
