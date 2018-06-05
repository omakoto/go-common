package common

import (
	"os"
	"sync"
)

type ExitFunc func()

var (
	atExitsSem = &sync.Mutex{}
	atExits    []ExitFunc
)

type exitStatus struct {
	code int
}

// ExitSuccess should be used within RunAndExit to cleanly finishes the process with a success code.
func ExitSuccess() {
	Exit(true)
}

// ExitFailure should be used within RunAndExit to cleanly finishes the process with a failure code.
func ExitFailure() {
	Exit(false)
}

// Exit should be used within RunAndExit to cleanly finishes the process.
func Exit(success bool) {
	status := 1
	if success {
		status = 0
	}
	ExitWithStatus(status)
}

// ExitWithStatus should be used within RunAndExit to cleanly finishes the process with a given status code.
func ExitWithStatus(status int) {
	panic(exitStatus{status})
}

// AtExit registers an at-exit hook function.
func AtExit(f ExitFunc) {
	atExitsSem.Lock()
	atExits = append(atExits, f)
	atExitsSem.Unlock()
}

// RunAtExits runs (and removes) all registered AtExit functions. RunAndExit will call it automatically,
// so no need to call it when you use RunAndExit.
func RunAtExits() {
	for {
		atExitsSem.Lock()
		if len(atExits) == 0 {
			atExitsSem.Unlock()
			return
		}
		last := len(atExits) - 1
		lastFunc := atExits[last]
		atExits[last] = nil
		atExits = atExits[0:last]

		atExitsSem.Unlock()

		lastFunc()
	}
}

func runWithRescue(f func() int) (result int) {
	defer RunAtExits()
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(exitStatus); ok {
				result = e.code
			} else {
				panic(r)
			}
		}
	}()
	result = f()
	return
}

// RunAndExit executes a given function. Within the function, util.Exit* functions can be used to finish the process cleanly.
func RunAndExit(f func() int) {
	os.Exit(runWithRescue(f))
}
