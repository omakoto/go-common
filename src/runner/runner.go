package runner

import (
	_ "embed"
	"github.com/omakoto/go-common/src/common"
	"github.com/omakoto/go-common/src/must"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

const envSkipGen = "GO_RUNNER_SKIP_GEN"

//go:embed runner.txt
var scriptContent string

type Options struct {
	WrapperPath string
}

func (o *Options) wrapperPath(mainSource string) string {
	if o.WrapperPath != "" {
		return o.WrapperPath
	}
	w, ok := strings.CutSuffix(mainSource, ".go")
	if !ok {
		panic("Unexpected filename extension in " + mainSource)
	}
	return w
}

func extractOptions(options []Options) Options {
	if len(options) == 0 {
		return Options{}
	}
	return options[0]
}

// GenWrapper is supposed to be called by a go program executed by `go run`, and creates
// a wrapper shell script for the program.
// See misc/runner-test.go
func GenWrapper(options ...Options) {
	if os.Getenv(envSkipGen) == "1" {
		return
	}

	if _, file, _, ok := runtime.Caller(1); !ok {
		panic("runtime.Caller failed")
	} else {
		generate(file, extractOptions(options))
	}
}

func getRelPath(sourcePath, wrapperPath string) (replath, wrapperResolvedPath string) {
	sourceDir := filepath.Dir(sourcePath)
	wrapperResolvedPath = filepath.Join(sourceDir, wrapperPath)
	wrapperDir := filepath.Dir(wrapperResolvedPath)
	return must.Must2(filepath.Rel(wrapperDir, sourcePath)), wrapperResolvedPath
}

func generate(mainSource string, options Options) {
	//fmt.Fprintf(os.Stderr, "s=%s\n", mainSource) // mainSource should be a fullpath here.
	wrapperFile := options.wrapperPath(mainSource)

	script := scriptContent
	script = strings.Replace(script, "{ENV}", envSkipGen, -1)

	rel, wrapperResolvedPath := getRelPath(mainSource, wrapperFile)

	common.Verbosef("Rel=%s", rel)
	common.Verbosef("Wrapper=%s", wrapperResolvedPath)
	script = strings.Replace(script, "{SOURCE}", rel, -1)

	f := must.Must2(os.OpenFile(wrapperResolvedPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755))

	_ = must.Must2(f.WriteString(script))
	must.Must(f.Close())

	must.Must(syscall.Exec(wrapperResolvedPath, os.Args, os.Environ()))
}
