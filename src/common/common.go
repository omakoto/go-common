package common

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

var (
	Quiet = false

	cachedBinName = ""

	DebugEnabled = false
)

func init() {
	if GetBinEnv("DEBUG") == "1" || os.Getenv("DEBUG") == "1" {
		DebugEnabled = true
	}
}

// MustGetExecutable returns the filename of the self binary.
func MustGetExecutable() string {
	me, err := os.Executable()
	Check(err, "os.Executable() failed")
	return me
}

// MustGetBinName returns the filename of the currently running executable.
func MustGetBinName() string {
	if cachedBinName == "" {
		cachedBinName = filepath.Base(MustGetExecutable())
	}
	return cachedBinName
}

func GetBinEnv(suffix string) string {
	return os.Getenv(fmt.Sprintf("%s_%s", strings.Replace(strings.ToUpper(MustGetBinName()), "-", "_", -1), suffix))
}

func MustGetenv(name string) string {
	ret := os.Getenv(name)
	if ret == "" {
		Fatal(name + " not set")
	}
	return ret
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

func CheckPanic(err error, message string) {
	if err == nil {
		return
	}
	panic(fmt.Sprintf("%s: %s", message, err))
}

func CheckPanicf(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}
	CheckPanic(err, fmt.Sprintf(format, args...))
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

func Debug(message string) {
	if !DebugEnabled {
		return
	}
	fmt.Fprint(os.Stderr, message, maybeLf(message))
}

func Debugf(format string, args ...interface{}) {
	if !DebugEnabled {
		return
	}
	Debug(fmt.Sprintf(format, args...))
}

func Dump(prefix string, object interface{}) {
	if !DebugEnabled {
		return
	}
	Debugf("%s%s", prefix, spew.Sdump(object))
}