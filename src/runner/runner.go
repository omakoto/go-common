package runner

import (
	"fmt"
	"github.com/omakoto/go-common/src/common"
	"os"
	"runtime"
	"strings"
	"syscall"
)

const envSkipGen = "GO_RUNNER_SKIP_GEN"

const scriptContent = `#!/bin/sh

export ` + envSkipGen + `=1

exec go run "$0.go" "${@}"
`

// GenWrapper is supposed to be called by a go program executed by `go run`, and creates
// a wrapper shell script for the program.
// See misc/runner-test.go
func GenWrapper() {
	if os.Getenv(envSkipGen) == "1" {
		return
	}

	if _, file, _, ok := runtime.Caller(1); !ok {
		panic("runtime.Caller failed")
	} else {
		if script, ok := strings.CutSuffix(file, ".go"); !ok {
			panic("Unexpected filename extension in " + file)
		} else {
			common.Verbosef("Creating %s ...", script)

			f, err := os.OpenFile(script, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
			common.Checkf(err, "OpenFile failed: %s", err)
			_, err = f.WriteString(scriptContent)
			common.Checkf(err, "WriteString: %s", err)
			f.Close()

			common.Verbosef("Running %s ...", script)
			err = syscall.Exec(script, os.Args[1:], os.Environ())
			panic(fmt.Sprintf("Exec failed: %v", err))
		}
	}
}
