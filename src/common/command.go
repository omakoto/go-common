package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
)

var (
	Quiet = false

	cachedBinName = ""
)

// MustGetBinName returns the filename of the currently running executable.
func MustGetBinName() string {
	if cachedBinName == "" {
		me, err := os.Executable()
		Check(err, "os.Executable() failed")
		cachedBinName = filepath.Base(me)
	}
	return cachedBinName
}

func maybePrintStackTrack() {
	if GetBinEnv("PRINT_STACK") == "1" {
		debug.PrintStack()
	}
}

func maybeLf(msg string) string {
	l := len(msg)
	if l > 0 && msg[l-1] == '\n' {
		return ""
	}
	return "\n"
}

func Fatal(message string) {
	if !Quiet {
		Warn(message)
		maybePrintStackTrack()
	}
	ExitFailure()
}

func Fatalf(format string, args ...interface{}) {
	Fatal(fmt.Sprintf(format, args))
}

func Warn(message string) {
	if !Quiet {
		fmt.Fprint(os.Stderr, MustGetBinName(), ": ", message, maybeLf(message))
	}
}

func Warnf(format string, args ...interface{}) {
	Warn(fmt.Sprintf(format, args))
}

func Check(err error, message string) {
	if err == nil {
		return
	}
	Fatal(fmt.Sprintf("%s: %s", message, err))
}

func Checkf(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	Check(err, fmt.Sprintf(format, args...))
}

func OrFatalf(condition bool, format string, args ...interface{}) {
	if !condition {
		Fatalf(format, args...)
	}
}

func OrWarnf(condition bool, format string, args ...interface{}) {
	if !condition {
		Warnf(format, args...)
	}
}
