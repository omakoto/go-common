package fileinput

import (
	"bufio"
	"fmt"
	"iter"
	"os"

	"github.com/omakoto/go-common/src/common"
	"github.com/otiai10/copy"
)

// Options is options for [FileInput].
type Options struct {
	// If true, STDOUT will be redirected to the input file. For inline replacements.
	InlineReplace bool

	// Input files. Defaults to [os.Args].
	Files []string

	// Backup suffix. Defaults to ".bak".
	BackupSuffix string
}

// FileInfo contains information about the current line.
type FileInfo struct {
	// Filename of the current input file.
	Filename string

	// Line number within the file. (0 based)
	Line int
}

type FileInputSeq func(yield func(text, file string, line int) bool)

// FileInput is similar to Python's fileinput.
func FileInput(options_ ...Options) iter.Seq2[string, FileInfo] {
	options := Options{}
	if len(options_) > 0 {
		options = options_[0]
	}

	replace := options.InlineReplace
	suffix := options.BackupSuffix
	if suffix == "" {
		suffix = ".bak"
	}
	files := options.Files
	if len(files) == 0 {
		files = os.Args[1:]
		if len(files) == 0 {
			files = append(files, "/dev/stdin")
		}
	}

	opener := func(file string) (*os.File, *os.File, error) {
		if !replace {
			in, err := os.Open(file)
			if err != nil {
				return nil, nil, fmt.Errorf("unable to open file '%s': %w", file, err)
			}
			return in, nil, nil
		} else {
			backup := file + suffix
			err := copy.Copy(file, backup, copy.Options{PermissionControl: copy.PerservePermission})
			if err != nil {
				return nil, nil, fmt.Errorf("unable to create backup file '%s' for '%s': %w", backup, file, err)
			}
			in, err := os.Open(backup)
			if err != nil {
				return nil, nil, fmt.Errorf("unable to open file '%s': %w", file, err)
			}
			out, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC, 0)
			if err != nil {
				return nil, nil, fmt.Errorf("unable to open file '%s': %w", file, err)
			}
			os.Stdout = out

			return in, out, nil
		}
	}

	doSingle := func(file string, yield func(text string, info FileInfo) bool) {
		in, out, err := opener(file)
		common.Checke(err)
		defer in.Close()
		if out != nil {
			defer out.Close()
		}

		sc := bufio.NewScanner(in)

		fi := FileInfo{file, 0}

		for sc.Scan() {
			if !yield(sc.Text(), fi) {
				return
			}
			fi.Line++
		}
	}

	return func(yield func(text string, info FileInfo) bool) {
		for _, file := range files {
			doSingle(file, yield)
		}
	}
}
